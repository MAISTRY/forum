package handlers

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func MarkAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	NotificationID := r.FormValue("notificationID")
	if NotificationID == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM Notification WHERE NotificationID = ?; `, NotificationID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
