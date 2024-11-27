package handlers

import (
	"database/sql"
	"fmt"
	"forum/DB"
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
		http.Error(w, "Method not allowed", http.StatusOK)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	r.ParseForm()
	postID := r.FormValue("postId")

	err = DB.DelPost(db, postID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("HX-Redirect", "/")
	fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/categories"></head></html>`)
}
