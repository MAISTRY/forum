package handlers

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	GetUserByEmailQuery = `SELECT password, UserID FROM User WHERE username = ? OR email = ?`
	InsertSessionQuery  = `INSERT INTO Session (session_id, user_id, created_at, expiry_date, ip_address) VALUES (?,?,?,?,?)`
	checkSessionQuery   = `SELECT session_id FROM Session WHERE user_id = ?`
)

// LoginHandler handles user login requests.
// It authenticates the user using their username or email and password.
// If the credentials are valid, it creates a new session token, stores it in the database,
// and sets a cookie with the session token for subsequent requests.
// If the user already has an active session, it deletes the previous session and creates a new one.
//
// Parameters:
//   - w: An http.ResponseWriter to write the response.
//   - r: An *http.Request containing the incoming HTTP request.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	password := r.FormValue("password")
	username := r.FormValue("username")

	var hashedPassword string
	var usrID string

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error Opening the DB %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("Error Pinging the DB %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	stmt, err := db.Prepare(GetUserByEmailQuery)
	if err != nil {
		log.Printf("Error preparing the statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(username, username).Scan(&hashedPassword, &usrID)
	if err != nil {
		log.Printf("Error Querying the DB %v\n", err)
		if err == sql.ErrNoRows {
			http.Error(w, "username or email doesn't exist", http.StatusOK)
		} else {
			http.Error(w, "Internal Server Error", http.StatusOK)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Error Comparing Hashes %v\n", err)
		http.Error(w, "Invalid credentials", http.StatusOK)
		return
	}

	insertSession := func() error {
		sessionToken, err := utils.GenerateSessionToken()
		if err != nil {
			log.Printf("Error generating session token: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return fmt.Errorf(`error generating session token: %v\n`, err)
		}
		expiryDate := time.Now().Add(72 * time.Hour)
		ipAddr := utils.GetIP(r)

		stmt, err = db.Prepare(InsertSessionQuery)
		if err != nil {
			log.Printf("Error preparing the insertion statement: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return fmt.Errorf(`error preparing the insertion statement: %v\n`, err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(sessionToken, usrID, time.Now(), expiryDate, ipAddr)
		if err != nil {
			log.Printf("Error Inserting to the DB: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return fmt.Errorf(`error inserting to the DB: %v\n`, err)
		}

		yummyCookie := &http.Cookie{
			Name:     "sessionID",
			Value:    sessionToken,
			Expires:  expiryDate,
			Path:     "/",
			Secure:   true, // Set to true if testing with HTTPS
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
		}

		http.SetCookie(w, yummyCookie)
		log.Printf("Attempting to set cookie: %s\n", yummyCookie.String())
		return nil
	}

	// ! START: CHECK if user have another session
	stmt, err = db.Prepare(checkSessionQuery)
	if err != nil {
		log.Printf("Error preparing the statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(usrID)
	err = row.Err()
	if err != nil {
		log.Printf("Error querying the DB: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	if row == nil {
		err = insertSession()
		if err != nil {
			log.Printf("Error inserting to the DB after checking previous session: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
	} else {
		// Delete old session and cookie
		stmt, err := db.Prepare("DELETE FROM Session WHERE user_id =?")
		if err != nil {
			log.Printf("Error preparing the deletion statement: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(usrID)
		if err != nil {
			log.Printf("Error deleting from the DB: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}

		expiredCookie := &http.Cookie{
			Name:     "sessionID",
			Value:    "",
			Expires:  time.Unix(0, 0),
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, expiredCookie)

		err = insertSession()
		if err != nil {
			log.Printf("Error inserting to the DB after deleting previous session: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
	}

	// Redirect the user after login
	w.Header().Set("HX-Redirect", "/home")
	w.WriteHeader(http.StatusOK)

	// ! END: CHECK if user have another session
	//https://community.auth0.com/t/how-to-log-user-out-after-cookie-expiration-in-go/151913
	w.Write([]byte("Login successful"))

	w.Header().Set("HX-Redirect", "/")
	fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)
}
