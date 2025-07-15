package handlers

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	GetUserByUnameQuery     = `SELECT UserID FROM User WHERE username = ?;`
	GetUserByUserEmailQuery = `SELECT UserID FROM User WHERE email = ?;`
	InsertNewUserQuery      = `INSERT INTO User (username, firstname, lastname, email, password, gender) VALUES (?,?,?,?,?,?);`
	insertNewSessionQuery   = `INSERT INTO Session (session_id, user_id, created_at, expiry_date, ip_address) VALUES (?,?,?,?,?)`
)

// RegisterHandler handles the user registration process.
// It validates the registration form data, checks for existing users,
// creates a new user account, and sets up a session for the newly registered user.
//
// Parameters:
//   - w: http.ResponseWriter - The response writer to send the HTTP response.
//   - r: *http.Request - The HTTP request containing the registration form data.
//
// This function does not return any values, but it writes to the http.ResponseWriter:
//   - On success: Sets a session cookie and redirects to the home page.
//   - On failure: Sends an appropriate error message back to the client.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusOK)
		log.Printf("error connecting to the database: %s\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusOK)
		log.Printf("error pinging the database: %s\n", err)
		return
	}

	r.ParseForm()
	username := r.FormValue("newUsername")
	firstNmae := r.FormValue("fistName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("Email")
	psswd := r.FormValue("newPassword")
	cnfrmPswd := r.FormValue("ConfirmNewPassword")
	gender := r.FormValue("gender")

	// ! STSRT: to check if the user already exists in the database.
	// prep the stmt to get the user by their username.
	stmt, err := db.Prepare(GetUserByUnameQuery)
	if err != nil {
		log.Printf("error preparing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()
	var existingUser int
	err = stmt.QueryRow(username).Scan(&existingUser)
	if err != nil {
		log.Printf("error querying the DB: %v\n", err)
		if err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
	} else {
		log.Printf("Username is already taken\n")
		http.Error(w, "Username is already taken", http.StatusOK)
		r.Form = nil
		return
	}

	// prep the stmt to get the user by the userEmail.
	stmt, err = db.Prepare(GetUserByUserEmailQuery)
	if err != nil {
		log.Printf("error preparing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()

	var existingEmail string
	err = stmt.QueryRow(email).Scan(&existingEmail)
	if err != nil {
		log.Printf("error querying the DB: %v\n", err)
		if err != sql.ErrNoRows {
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
	} else {
		log.Printf("Email is already taken\n")
		http.Error(w, "Email is already taken", http.StatusOK)
		r.Form = nil
		return
	}

	isValidEmail := func(email string) bool {
		regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(regex)
		return re.MatchString(email)
	}
	if !isValidEmail(email) {
		log.Printf("Invalid email format\n")
		http.Error(w, "Invalid email format!", http.StatusOK)
		r.Form = nil
		return
	}
	// ! END

	// ! START: check the password validity.
	if len(psswd) < 8 {
		log.Printf("Password must be at least 8 characters long\n")
		http.Error(w, "Password must be at least 8 characters long", http.StatusOK)
		delete(r.Form, "newPassword")
		delete(r.Form, "ConfirmNewPassword")
		return
	}

	if psswd != cnfrmPswd {
		log.Printf("Passwords do not match\n")
		http.Error(w, "Passwords do not match", http.StatusOK)
		delete(r.Form, "ConfirmNewPassword")
		return
	}
	//! END

	// ! START: hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(psswd), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error generating hash: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	// ! END

	// ! START: insert the user into the database.
	stmt, err = db.Prepare(InsertNewUserQuery)
	if err != nil {
		log.Printf("error preparing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()

	if gender == "" {
		log.Printf("Invalid gender\n")
		http.Error(w, "Select a gender!", http.StatusOK)
		r.Form = nil
		return
	}
	if gender != "M" && gender != "F" {
		log.Printf("Invalid gender\n")
		http.Error(w, "Invalid gender!", http.StatusOK)
		r.Form = nil
		return
	}

	_, err = stmt.Exec(username, firstNmae, lastName, email, hashedPassword, gender)
	if err != nil {
		log.Printf("error executing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	// ! END

	log.Printf("User %s registered successfully\n", username)

	// *** create the 🍪 and redirect the user to the homepage. *** \\

	stmt, err = db.Prepare(GetUserByUnameQuery)
	if err != nil {
		log.Printf("error preparing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()

	var userID int

	err = stmt.QueryRow(username).Scan(&userID)
	if err != nil {
		log.Printf("error getting user ID: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	sessionToken, err := utils.GenerateSessionToken()
	if err != nil {
		log.Printf("error generating session token: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	expiryDate := time.Now().Add(72 * time.Hour)
	yummyCookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionToken,
		Path:     "/",
		Expires:  expiryDate,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	stmt, err = db.Prepare(insertNewSessionQuery)
	if err != nil {
		log.Printf("error preparing statement: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer stmt.Close()

	ipAddr := utils.GetIP(r)
	_, err = stmt.Exec(sessionToken, userID, time.Now(), expiryDate, ipAddr)
	if err != nil {
		log.Printf("error inserting session into the database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	http.SetCookie(w, yummyCookie)
	w.Write([]byte("Registration successful"))
	w.Header().Set("HX-Redirect", "/")
	fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)

}
