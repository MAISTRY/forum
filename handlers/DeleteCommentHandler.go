package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// DeleteCommentHandler handles HTTP POST requests to delete a comment from the forum.
// It expects a form value "commentId" representing the ID of the comment to be deleted.
// Only admins (privilege level 3) can delete comments.
// The function opens a connection to the SQLite database "meow.db", deletes the comment with the given ID,
// and returns a success response.
//
// If the request method is not POST, it returns a "Method not allowed" error.
// If the user is not an admin, it returns an "Unauthorized" error.
// If there is an error opening the database connection or deleting the comment,
// it returns an "Internal Server Error" response.
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin (only admins can delete comments)
	cookie, err := r.Cookie("sessionID")
	if err != nil || !isValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	privilege, err := getPrivilege(cookie.Value)
	if err != nil || privilege != 3 {
		http.Error(w, "Unauthorized - Admin access required", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	r.ParseForm()
	commentID := r.FormValue("commentId")

	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	// Start transaction to ensure atomicity
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete comment likes first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM CommentLike WHERE CommentID = ?", commentID)
	if err != nil {
		log.Printf("Error deleting comment likes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Delete comment dislikes (foreign key constraint)
	_, err = tx.Exec("DELETE FROM CommentDislike WHERE CommentID = ?", commentID)
	if err != nil {
		log.Printf("Error deleting comment dislikes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Delete the comment itself
	result, err := tx.Exec("DELETE FROM Comment WHERE CommentID = ?", commentID)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if comment was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true, "message": "Comment deleted successfully"}`)
}
