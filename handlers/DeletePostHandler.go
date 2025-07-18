package handlers

import (
	"database/sql"
	"fmt"
	"forum/DB"
	"log"
	"net/http"
)

// DelPostHandler handles HTTP POST requests to delete a post from the forum.
// It expects a form value "postId" representing the ID of the post to be deleted.
// The function opens a connection to the SQLite database "meow.db", deletes the post with the given ID,
// and redirects the client to the "/categories" page.
//
// If the request method is not POST, it returns a "Method not allowed" error.
// If there is an error opening the database connection or deleting the post,
// it returns an "Internal Server Error" response.
//
// The function uses the HX-Redirect header to perform the client-side redirect.
func DelPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin or moderator
	cookie, err := r.Cookie("sessionID")
	if err != nil || !isValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	privilege, err := getPrivilege(cookie.Value)
	if err != nil || privilege < 2 { // Must be moderator (2) or admin (3)
		http.Error(w, "Unauthorized - Moderator or Admin access required", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	r.ParseForm()
	postID := r.FormValue("postId")

	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	err = DB.DelPost(db, postID)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "message": "Post deleted successfully"}`)
}
