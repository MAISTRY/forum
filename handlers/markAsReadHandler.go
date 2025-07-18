package handlers

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func MarkAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Debug: log all form values
	log.Printf("All form values: %v", r.Form)
	log.Printf("PostForm values: %v", r.PostForm)

	NotificationID := r.FormValue("notificationID")
	log.Printf("Received NotificationID: '%s'", NotificationID)

	if NotificationID == "" {
		log.Printf("NotificationID is empty")
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("Marking notification %s as read", NotificationID)

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get user ID from session to ensure they can only mark their own notifications as read
	userIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user ID from cookie: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("User %s attempting to mark notification %s as read", userIDStr, NotificationID)

	// Update notification only if it belongs to the current user
	result, err := db.Exec(`UPDATE Notification SET IsRead = 1 WHERE NotificationID = ? AND UserToNotify = ?`, NotificationID, userIDStr)
	if err != nil {
		log.Printf("Error updating notification: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Rows affected: %d", rowsAffected)

	if rowsAffected == 0 {
		log.Printf("No rows affected - notification not found or unauthorized")
		http.Error(w, "Notification not found or unauthorized", http.StatusNotFound)
		return
	}

	log.Printf("Successfully marked notification %s as read", NotificationID)
	w.WriteHeader(http.StatusOK)
}
