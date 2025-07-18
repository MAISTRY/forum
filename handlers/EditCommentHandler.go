package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type EditCommentRequest struct {
	CommentID string `json:"comment_id"`
	Content   string `json:"content"`
}

type EditCommentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EditCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.CommentID == "" || strings.TrimSpace(req.Content) == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(EditCommentResponse{
			Success: false,
			Message: "Comment content cannot be empty",
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

	commentID, err := strconv.Atoi(req.CommentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the comment
	var commentOwnerID int
	err = db.QueryRow("SELECT UserID FROM Comment WHERE CommentID = ?", commentID).Scan(&commentOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(EditCommentResponse{
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
		json.NewEncoder(w).Encode(EditCommentResponse{
			Success: false,
			Message: "You can only edit your own comments",
		})
		return
	}

	// Update the comment
	_, err = db.Exec("UPDATE Comment SET content = ? WHERE CommentID = ?", 
		strings.TrimSpace(req.Content), commentID)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EditCommentResponse{
		Success: true,
		Message: "Comment updated successfully",
	})
}

// GetCommentForEdit returns comment data for editing
func GetCommentForEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commentID := r.URL.Query().Get("commentId")
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
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

	// Get comment data
	var comment struct {
		CommentID int    `json:"comment_id"`
		Content   string `json:"content"`
		UserID    int    `json:"user_id"`
	}

	err = db.QueryRow("SELECT CommentID, content, UserID FROM Comment WHERE CommentID = ?", commentID).
		Scan(&comment.CommentID, &comment.Content, &comment.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting comment: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user owns the comment
	if comment.UserID != userID {
		http.Error(w, "You can only edit your own comments", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
