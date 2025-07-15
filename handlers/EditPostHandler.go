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

type EditPostRequest struct {
	PostID  string `json:"post_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EditPostResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EditPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.PostID == "" || strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(EditPostResponse{
			Success: false,
			Message: "Title and content cannot be empty",
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

	postID, err := strconv.Atoi(req.PostID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if user owns the post
	var postOwnerID int
	err = db.QueryRow("SELECT UserID FROM Post WHERE PostID = ?", postID).Scan(&postOwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(EditPostResponse{
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
		json.NewEncoder(w).Encode(EditPostResponse{
			Success: false,
			Message: "You can only edit your own posts",
		})
		return
	}

	// Update the post
	_, err = db.Exec("UPDATE Post SET title = ?, content = ? WHERE PostID = ?", 
		strings.TrimSpace(req.Title), strings.TrimSpace(req.Content), postID)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EditPostResponse{
		Success: true,
		Message: "Post updated successfully",
	})
}

// GetPostForEdit returns post data for editing
func GetPostForEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("postId")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
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

	// Get post data
	var post struct {
		PostID  int    `json:"post_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		UserID  int    `json:"user_id"`
	}

	err = db.QueryRow("SELECT PostID, title, content, UserID FROM Post WHERE PostID = ?", postID).
		Scan(&post.PostID, &post.Title, &post.Content, &post.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting post: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user owns the post
	if post.UserID != userID {
		http.Error(w, "You can only edit your own posts", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
