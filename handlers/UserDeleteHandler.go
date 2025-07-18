package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserDeletePostHandler allows users to delete their own posts
func UserDeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	postID := r.FormValue("postId")
	if postID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "Post ID is required",
		})
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get user ID from session
	userIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the post
	var postOwnerID int
	err = db.QueryRow("SELECT UserID FROM Post WHERE PostID = ?", postIDInt).Scan(&postOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(DeleteResponse{
				Success: false,
				Message: "Post not found",
			})
			return
		}
		log.Printf("Error checking post ownership: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if postOwnerID != userID {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "You can only delete your own posts",
		})
		return
	}

	// Delete the post using transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete related data first (foreign key constraints)
	queries := []string{
		"DELETE FROM CommentLike WHERE CommentID IN (SELECT CommentID FROM Comment WHERE PostID = ?)",
		"DELETE FROM CommentDislike WHERE CommentID IN (SELECT CommentID FROM Comment WHERE PostID = ?)",
		"DELETE FROM Comment WHERE PostID = ?",
		"DELETE FROM PostLike WHERE PostID = ?",
		"DELETE FROM PostDislike WHERE PostID = ?",
		"DELETE FROM PostCategory WHERE PostID = ?",
		"DELETE FROM Notification WHERE PostID = ?",
		"DELETE FROM Post WHERE PostID = ?",
	}

	for _, query := range queries {
		_, err = tx.Exec(query, postIDInt)
		if err != nil {
			log.Printf("Error executing query %s: %v", query, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteResponse{
		Success: true,
		Message: "Post deleted successfully",
	})
}

// UserDeleteCommentHandler allows users to delete their own comments
func UserDeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	commentID := r.FormValue("commentId")
	if commentID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "Comment ID is required",
		})
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get user ID from session
	userIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the comment
	var commentOwnerID int
	err = db.QueryRow("SELECT UserID FROM Comment WHERE CommentID = ?", commentIDInt).Scan(&commentOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(DeleteResponse{
				Success: false,
				Message: "Comment not found",
			})
			return
		}
		log.Printf("Error checking comment ownership: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if commentOwnerID != userID {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeleteResponse{
			Success: false,
			Message: "You can only delete your own comments",
		})
		return
	}

	// Delete the comment using transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete related data first (foreign key constraints)
	queries := []string{
		"DELETE FROM CommentLike WHERE CommentID = ?",
		"DELETE FROM CommentDislike WHERE CommentID = ?",
		"DELETE FROM Notification WHERE CommentID = ?",
		"DELETE FROM Comment WHERE CommentID = ?",
	}

	for _, query := range queries {
		_, err = tx.Exec(query, commentIDInt)
		if err != nil {
			log.Printf("Error executing query %s: %v", query, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteResponse{
		Success: true,
		Message: "Comment deleted successfully",
	})
}
