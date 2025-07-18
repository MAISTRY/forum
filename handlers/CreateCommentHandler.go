package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/DB"
	"log"
	"net/http"
	"strconv"
)

func CreatCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postID := req.PostID
	comment := req.Comment
	if comment == "" {
		http.Error(w, "Comment cannot be empty", http.StatusOK)
		return
	}

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("The Error getting user ID %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting user ID to int %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	intPostID, err := strconv.Atoi(postID)
	if err != nil {
		log.Printf("Error converting Post ID to int %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	cmntID, err := DB.InsertComment(intPostID, intUserID, comment)
	if err != nil {
		log.Printf("Error inserting comment %v\n", err)
		http.Error(w, "Failed to post comment. Please try again.", http.StatusOK)
		return
	}

	username, err := DB.GetCommentsUsername(db, userID)
	if err != nil {
		log.Printf("Error getting username %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	// Create notification for post owner
	insertCommentNotification(db, userID, postID)

	commnetObject := CommentedPost{
		UserID:     intUserID,
		UserName:   username,
		CommentID:  int(cmntID),
		PostID:     intPostID,
		Comment:    comment,
		CreateDate: "now",
		Likes:      0,
		Dislikes:   0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commnetObject)
}

// insertCommentNotification creates a notification for the post owner when someone comments
func insertCommentNotification(db *sql.DB, userID, postID string) {
	stmt, err := db.Prepare(`
        INSERT INTO Notification (UserID, UserToNotify, PostID, NotificationType)
        VALUES (?,?,?,?);
    `)
	if err != nil {
		log.Printf("Error preparing comment notification statement: %v\n", err)
		return
	}
	defer stmt.Close()

	postOwnerID, err := GetPostOwnerID(postID, db)
	if err != nil {
		log.Printf("Error getting post owner ID for comment notification: %v\n", err)
		return
	}

	// Don't create notification if user is commenting on their own post
	if userID == postOwnerID {
		return
	}

	_, err = stmt.Exec(userID, postOwnerID, postID, "Comment")
	if err != nil {
		log.Printf("Error inserting comment notification: %v\n", err)
		return
	}
}
