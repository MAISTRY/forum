package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CheckAuthHandler is a HTTP handler function that checks the authentication status and privilege level of a user.
// It retrieves the session ID from the request cookies and validates it against the database.
// If the session is valid, it retrieves the user's privilege level and returns a JSON response indicating the authentication status and privilege level.
// If the session is invalid, it returns a JSON response indicating the authentication status as false and privilege level as 0.
//
// Parameters:
//   - w: An http.ResponseWriter object to write the response.
//   - r: An http.Request object containing the request data.
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	privilege := 0
	cookie, err := r.Cookie("sessionID")
	if err != nil || !isValidSession(cookie.Value) {
		jsonResp := fmt.Sprintf(`{
            "authenticated": false,
            "privilege": %d}`, privilege)

		// fmt.Println(jsonResp)
		w.Write([]byte(jsonResp))
		return
	}

	privilege, err = getPrivilege(cookie.Value)
	if err != nil {
		log.Printf("Error getting privilege: %v\n", err)
		return
	}

	jsonResp := fmt.Sprintf(`{
        "authenticated": true,
        "privilege": %d}`, privilege)

	// fmt.Println(jsonResp)
	w.Write([]byte(jsonResp))
}

// getPrivilege retrieves the privilege level of a user associated with a given session ID.
//
// The function opens a connection to the SQLite database, prepares two SQL statements,
// and executes them to retrieve the user ID and privilege level associated with the session ID.
// If any error occurs during the database operations, the function logs the error and returns -1 and an error.
//
// Parameters:
//   - session_id: A string representing the unique identifier of the session.
//
// Returns:
//   - An integer representing the privilege level of the user associated with the session ID.
//     If an error occurs, the function returns -1.
//   - An error object indicating any error that occurred during the database operations.
func getPrivilege(session_id string) (int, error) {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening DB: %v\n", err)
		return -1, fmt.Errorf("error opening DB: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`select user_id from Session where session_id = ?`)
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		return -1, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	var user_id int
	err = stmt.QueryRow(session_id).Scan(&user_id)
	if err != nil {
		log.Printf("Error getting user ID: %v\n", err)
		return -1, fmt.Errorf("error getting user ID: %v", err)
	}

	var privilege int
	stmt, err = db.Prepare(`select privilege from User where UserID = ?`)
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		return -1, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(user_id).Scan(&privilege)
	if err != nil {
		log.Printf("Error getting privilege: %v\n", err)
		return -1, fmt.Errorf("error getting privilege: %v", err)
	}

	return privilege, nil
}

// isValidSession checks if a given session ID is valid and not expired.
//
// Parameters:
//   - sessionID: A string representing the unique identifier of the session to validate.
//
// Returns:
//
//	A boolean value indicating whether the session is valid (true) or invalid (false).
//	The session is considered invalid if it doesn't exist in the database,
//	if there's an error querying the database, or if the session has expired.
func isValidSession(sessionID string) bool {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening DB: %v\n", err)
		return false
	}
	defer db.Close()

	var expiryDate time.Time
	err = db.QueryRow(`SELECT expiry_date FROM Session WHERE session_id = ?`, sessionID).Scan(&expiryDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No session found for the provided session ID.")
		} else {
			log.Printf("Error querying session: %v\n", err)
		}
		return false
	}

	if time.Now().After(expiryDate) {
		log.Println("Session has expired.")
		return false
	}

	// fmt.Println("Session is valid.")
	return true
}
