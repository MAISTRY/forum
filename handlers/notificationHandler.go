package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type NotificaionBody struct {
	NotificationID   int    `json:"notification_id"`
	UserID           int    `json:"user_id"`
	UserToNotify     int    `json:"user_to_notify"`
	PostID           any    `json:"post_id"`
	CommentID        any    `json:"comment_id"`
	NotificationType string `json:"notification_type"`
	CreatedAt        string `json:"created_at"`
	IsRead           bool   `json:"is_read"`
	Username         string `json:"username"`
	PostTitle        string `json:"post_title"`
	CommentContent   string `json:"comment_content"`
}

func NotificaionHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		http.Error(w, "Internal Server Error 1", http.StatusOK)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error 2", http.StatusOK)
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT
			n.NotificationID, n.UserID, n.UserToNotify, n.PostID, n.CommentID,
			n.NotificationType, n.CreatedAt, n.IsRead,
			u.username,
			COALESCE(p.title, '') as post_title,
			COALESCE(c.content, '') as comment_content
		FROM Notification n
		JOIN User u ON n.UserID = u.UserID
		LEFT JOIN Post p ON n.PostID = p.PostID
		LEFT JOIN Comment c ON n.CommentID = c.CommentID
		WHERE n.UserToNotify = ?
		ORDER BY n.CreatedAt DESC
		LIMIT 50;
	`, userID)
	if err != nil {
		fmt.Printf("Error querying notifications: %v\n", err)
		http.Error(w, "Internal Server Error 3", http.StatusOK)
		return
	}
	defer rows.Close()

	var notifications []NotificaionBody
	for rows.Next() {
		var notification NotificaionBody
		err := rows.Scan(&notification.NotificationID, &notification.UserID, &notification.UserToNotify,
			&notification.PostID, &notification.CommentID, &notification.NotificationType,
			&notification.CreatedAt, &notification.IsRead, &notification.Username,
			&notification.PostTitle, &notification.CommentContent)
		if err != nil {
			fmt.Printf("Error scanning notification: %v\n", err)
			http.Error(w, "Internal Server Error 4", http.StatusOK)
			return
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Internal Server Error 5", http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Println(notifications)
	json.NewEncoder(w).Encode(notifications)
}

// NotificationCountHandler returns the count of unread notifications
func NotificationCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		http.Error(w, "Internal Server Error 1", http.StatusOK)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error 2", http.StatusOK)
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM Notification
		WHERE UserToNotify = ? AND IsRead = 0;
	`, userID).Scan(&count)
	if err != nil {
		fmt.Printf("Error querying notification count: %v\n", err)
		http.Error(w, "Internal Server Error 3", http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
