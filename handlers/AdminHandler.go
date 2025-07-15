package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// AdminStats represents the statistics for the admin dashboard
type AdminStats struct {
	AdminCount     int `json:"AdminCount"`
	ModeratorCount int `json:"ModeratorCount"`
	PostCount      int `json:"PostCount"`
	CommentCount   int `json:"CommentCount"`
}

// AdminUser represents a user in the admin management interface
type AdminUser struct {
	UserID    int    `json:"UserID"`
	Username  string `json:"Username"`
	Email     string `json:"Email"`
	Privilege int    `json:"Privilege"`
	CreatedAt string `json:"CreatedAt"`
}

// ModerationRequest represents a moderation request
type ModerationRequest struct {
	RequestID     int    `json:"RequestID"`
	UserID        int    `json:"UserID"`
	Username      string `json:"Username"`
	RequestDate   string `json:"RequestDate"`
	Status        string `json:"Status"`
	AdminResponse string `json:"AdminResponse,omitempty"`
	ResponseDate  string `json:"ResponseDate,omitempty"`
}

// UserPromotionRequest represents a request to promote/demote a user
type UserPromotionRequest struct {
	UserID    int `json:"userId"`
	Privilege int `json:"privilege"`
}

// ModerationResponseRequest represents a response to a moderation request
type ModerationResponseRequest struct {
	RequestID int    `json:"requestId"`
	Status    string `json:"status"`
}

// AdminStatsHandler returns statistics for the admin dashboard
func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stats := AdminStats{}

	// Get admin count
	err = db.QueryRow("SELECT COUNT(*) FROM User WHERE privilege = 3").Scan(&stats.AdminCount)
	if err != nil {
		log.Printf("Error getting admin count: %v", err)
		stats.AdminCount = 0
	}

	// Get moderator count
	err = db.QueryRow("SELECT COUNT(*) FROM User WHERE privilege = 2").Scan(&stats.ModeratorCount)
	if err != nil {
		log.Printf("Error getting moderator count: %v", err)
		stats.ModeratorCount = 0
	}

	// Get post count
	err = db.QueryRow("SELECT COUNT(*) FROM Post").Scan(&stats.PostCount)
	if err != nil {
		log.Printf("Error getting post count: %v", err)
		stats.PostCount = 0
	}

	// Get comment count
	err = db.QueryRow("SELECT COUNT(*) FROM Comment").Scan(&stats.CommentCount)
	if err != nil {
		log.Printf("Error getting comment count: %v", err)
		stats.CommentCount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// AdminUsersHandler returns users for admin management
func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	searchTerm := r.URL.Query().Get("search")
	var query string
	var args []interface{}

	if searchTerm != "" {
		query = "SELECT UserID, username, email, privilege, created_at FROM User WHERE username LIKE ? OR email LIKE ? ORDER BY privilege DESC, username"
		searchPattern := "%" + searchTerm + "%"
		args = []interface{}{searchPattern, searchPattern}
	} else {
		query = "SELECT UserID, username, email, privilege, created_at FROM User ORDER BY privilege DESC, username"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var user AdminUser
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Privilege, &user.CreatedAt)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// AdminPromoteUserHandler promotes a user to a higher privilege level
func AdminPromoteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UserPromotionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate privilege level
	if req.Privilege < 1 || req.Privilege > 3 {
		http.Error(w, "Invalid privilege level", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update user privilege
	_, err = db.Exec("UPDATE User SET privilege = ? WHERE UserID = ?", req.Privilege, req.UserID)
	if err != nil {
		log.Printf("Error updating user privilege: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User promoted successfully",
	})
}

// AdminDemoteUserHandler demotes a user to a lower privilege level
func AdminDemoteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UserPromotionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate privilege level
	if req.Privilege < 1 || req.Privilege > 3 {
		http.Error(w, "Invalid privilege level", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get current admin's user ID to prevent self-demotion
	adminUserID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting admin user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	adminID, err := strconv.Atoi(adminUserID)
	if err != nil {
		log.Printf("Error converting admin ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prevent self-demotion
	if adminID == req.UserID {
		http.Error(w, "You cannot demote yourself", http.StatusBadRequest)
		return
	}

	// Check if this would leave no admins (only when demoting an admin)
	// First get the current privilege of the user being demoted
	var currentPrivilege int
	err = db.QueryRow("SELECT privilege FROM User WHERE UserID = ?", req.UserID).Scan(&currentPrivilege)
	if err != nil {
		log.Printf("Error getting current user privilege: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Only check admin count if we're demoting an admin
	if currentPrivilege == 3 && req.Privilege < 3 {
		var adminCount int
		err = db.QueryRow("SELECT COUNT(*) FROM User WHERE privilege = 3").Scan(&adminCount)
		if err != nil {
			log.Printf("Error checking admin count: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if adminCount <= 1 {
			http.Error(w, "Cannot demote the last admin", http.StatusBadRequest)
			return
		}
	}

	// Update user privilege
	_, err = db.Exec("UPDATE User SET privilege = ? WHERE UserID = ?", req.Privilege, req.UserID)
	if err != nil {
		log.Printf("Error updating user privilege: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User demoted successfully",
	})
}

// AdminModerationRequestsHandler returns moderation requests
func AdminModerationRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
		SELECT mr.RequestID, mr.UserID, u.username, mr.RequestDate, mr.Status, 
		       COALESCE(mr.AdminResponse, '') as AdminResponse, 
		       COALESCE(mr.ResponseDate, '') as ResponseDate
		FROM ModerationRequest mr
		JOIN User u ON mr.UserID = u.UserID
		ORDER BY mr.RequestDate DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying moderation requests: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []ModerationRequest
	for rows.Next() {
		var req ModerationRequest
		err := rows.Scan(&req.RequestID, &req.UserID, &req.Username, &req.RequestDate,
			&req.Status, &req.AdminResponse, &req.ResponseDate)
		if err != nil {
			log.Printf("Error scanning moderation request: %v", err)
			continue
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// AdminRespondRequestHandler responds to a moderation request
func AdminRespondRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ModerationResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate status
	if req.Status != "approved" && req.Status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get admin ID
	cookie, _ := r.Cookie("sessionID")
	adminPrivilege, _ := getPrivilege(cookie.Value)
	if adminPrivilege != 3 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get admin user ID
	adminUserID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting admin user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	adminID, err := strconv.Atoi(adminUserID)
	if err != nil {
		log.Printf("Error converting admin ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update moderation request
	_, err = tx.Exec(`
		UPDATE ModerationRequest
		SET Status = ?, AdminID = ?, ResponseDate = ?
		WHERE RequestID = ?`,
		req.Status, adminID, time.Now().Format("2006-01-02 15:04:05"), req.RequestID)
	if err != nil {
		log.Printf("Error updating moderation request: %v", err)
		http.Error(w, "Failed to update request", http.StatusInternalServerError)
		return
	}

	// If approved, promote user to moderator
	if req.Status == "approved" {
		// Get user ID from request
		var userID int
		err = tx.QueryRow("SELECT UserID FROM ModerationRequest WHERE RequestID = ?", req.RequestID).Scan(&userID)
		if err != nil {
			log.Printf("Error getting user ID from request: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Promote user to moderator (privilege level 2)
		_, err = tx.Exec("UPDATE User SET privilege = 2 WHERE UserID = ?", userID)
		if err != nil {
			log.Printf("Error promoting user: %v", err)
			http.Error(w, "Failed to promote user", http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Request processed successfully",
	})
}

// CreateModerationRequestHandler creates a new moderation request
func CreateModerationRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	cookie, err := r.Cookie("sessionID")
	if err != nil || !isValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is a normal user (privilege 1)
	privilege, err := getPrivilege(cookie.Value)
	if err != nil || privilege != 1 {
		http.Error(w, "Only normal users can request moderation", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get user ID
	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if user already has a pending request
	var existingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM ModerationRequest WHERE UserID = ? AND Status = 'pending'", userIDInt).Scan(&existingCount)
	if err != nil {
		log.Printf("Error checking existing requests: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingCount > 0 {
		http.Error(w, "You already have a pending moderation request", http.StatusBadRequest)
		return
	}

	// Create new moderation request
	_, err = db.Exec("INSERT INTO ModerationRequest (UserID) VALUES (?)", userIDInt)
	if err != nil {
		log.Printf("Error creating moderation request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Moderation request created successfully",
	})
}

// isAdmin checks if the current user is an admin
func isAdmin(r *http.Request) bool {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return false
	}

	privilege, err := getPrivilege(cookie.Value)
	if err != nil {
		return false
	}

	return privilege == 3
}
