package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

//	type NotificaionBody struct {
//		UserID           int    `json:"user_id"`
//		UserToNotifying  int    `json:"user_to_notifying"`
//		Notification     string `json:"notification"`
//		NotificationType string `json:"notification_type"`
//		NotificationDate string `json:"notification_date"`
//		Read             bool   `json:"read"`
//		PostID           int    `json:"post_id"`
//		CommentID        int    `json:"comment_id"`
//	}

type NotificaionBody struct {
	NotificationID   int    `json:"notification_id"`
	UserID           int    `json:"user_id"`
	UserToNotify     int    `json:"user_to_notify"`
	PostID           any    `json:"post_id"`
	CommentID        any    `json:"comment_id"`
	NotificationType string `json:"notification_type"`
	CreatedAt        string `json:"created_at"`
	IsRead           bool   `json:"is_read"`
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
		SELECT * FROM Notification
		WHERE UserToNotify = ? AND IsRead = 0;
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
		err := rows.Scan(&notification.NotificationID, &notification.UserID, &notification.UserToNotify, &notification.PostID, &notification.CommentID, &notification.NotificationType, &notification.CreatedAt, &notification.IsRead)
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
