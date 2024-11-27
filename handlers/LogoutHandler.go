package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

// LogoutHandler handles the user logout process.
// It invalidates the user's session, clears the session cookie,
// and redirects the user to the home page.
//
// Parameters:
//   - w: http.ResponseWriter - The response writer to send the HTTP response.
//   - r: *http.Request - The HTTP request received from the client.
//
// The function does not return any value, but it writes to the response writer
// and sets headers to manage the logout process and redirection.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	// get the 🍪
	yummyCookie, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Printf("error getting cookie: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		w.Header().Set("HX-Redirect", "/")
		fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)
		return
	}

	sessionToken := yummyCookie.Value
	err = deleteSession(sessionToken)
	if err != nil {
		fmt.Printf("error getting cookie: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		w.Header().Set("HX-Redirect", "/")
		fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "sessionID",
		Value:  "",
		Path:   "/",
		Secure: true,
		MaxAge: -1,
	})

	// w.Write([]byte("Logout successful"))

	w.Header().Set("HX-Redirect", "/")
	fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)

}

// deleteSession removes a session from the database based on the provided session ID.
//
// Parameters:
//   - session_id: string - The unique identifier of the session to be deleted.
//
// Returns:
//   - error: An error if any step in the process fails, nil otherwise.
//     Possible errors include database connection issues, SQL query preparation problems,
//     or failures in executing the delete operation.
func deleteSession(session_id string) error {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		return fmt.Errorf("error opening the DB %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging the DB %v", err)
	}

	stmt, err := db.Prepare("DELETE FROM Session WHERE session_id =?")
	if err != nil {
		return fmt.Errorf("error preparing the statement %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(session_id)
	if err != nil {
		return fmt.Errorf("error deleting from the DB %v", err)
	}

	return nil
}
