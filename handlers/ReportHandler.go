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

// PostReport represents a post report
type PostReport struct {
	ReportID      int     `json:"ReportID"`
	PostID        int     `json:"PostID"`
	PostTitle     string  `json:"PostTitle"`
	PostContent   string  `json:"PostContent"`
	PostAuthor    string  `json:"PostAuthor"`
	ModeratorID   int     `json:"ModeratorID"`
	ModeratorName string  `json:"ModeratorName"`
	ReportDate    string  `json:"ReportDate"`
	Reason        string  `json:"Reason"`
	Status        string  `json:"Status"`
	AdminResponse string  `json:"AdminResponse"`
	AdminID       *int    `json:"AdminID"`
	ResponseDate  *string `json:"ResponseDate"`
}

// ReportPostHandler handles moderator reports of posts
func ReportPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is moderator or admin
	if !isModerator(r) {
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

	// Get moderator ID
	moderatorIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting moderator ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	moderatorID, err := strconv.Atoi(moderatorIDStr)
	if err != nil {
		log.Printf("Error converting moderator ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse form data
	r.ParseForm()
	postIDStr := r.FormValue("postId")
	reason := r.FormValue("reason")

	if postIDStr == "" || reason == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if post exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Post WHERE PostID = ?)", postID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking post existence: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if already reported by this moderator
	var reportExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM PostReport WHERE PostID = ? AND ModeratorID = ? AND Status = 'pending')", postID, moderatorID).Scan(&reportExists)
	if err != nil {
		log.Printf("Error checking existing report: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if reportExists {
		http.Error(w, "Post already reported by you", http.StatusConflict)
		return
	}

	// Insert report
	_, err = db.Exec("INSERT INTO PostReport (PostID, ModeratorID, Reason) VALUES (?, ?, ?)", postID, moderatorID, reason)
	if err != nil {
		log.Printf("Error inserting report: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post reported successfully",
	})
}

// AdminReportsHandler returns all post reports for admin dashboard
func AdminReportsHandler(w http.ResponseWriter, r *http.Request) {
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
		SELECT
			pr.ReportID, pr.PostID,
			COALESCE(p.title, '[Post Deleted]') as title,
			COALESCE(p.content, '[Post content no longer available]') as content,
			COALESCE(u1.username, '[Unknown User]') as post_author,
			pr.ModeratorID, u2.username as moderator_name, pr.ReportDate, pr.Reason,
			pr.Status, pr.AdminResponse, pr.AdminID, pr.ResponseDate
		FROM PostReport pr
		LEFT JOIN Post p ON pr.PostID = p.PostID
		LEFT JOIN User u1 ON p.UserID = u1.UserID
		JOIN User u2 ON pr.ModeratorID = u2.UserID
		ORDER BY pr.ReportDate DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying reports: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reports []PostReport
	for rows.Next() {
		var report PostReport
		var responseDate sql.NullString
		var adminID sql.NullInt64
		var adminResponse sql.NullString

		err := rows.Scan(
			&report.ReportID, &report.PostID, &report.PostTitle, &report.PostContent, &report.PostAuthor,
			&report.ModeratorID, &report.ModeratorName, &report.ReportDate, &report.Reason,
			&report.Status, &adminResponse, &adminID, &responseDate,
		)
		if err != nil {
			log.Printf("Error scanning report: %v", err)
			continue
		}

		// Handle nullable fields
		if adminResponse.Valid {
			report.AdminResponse = adminResponse.String
		} else {
			report.AdminResponse = ""
		}

		if adminID.Valid {
			adminIDInt := int(adminID.Int64)
			report.AdminID = &adminIDInt
		}
		if responseDate.Valid {
			report.ResponseDate = &responseDate.String
		}

		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

// AdminRespondReportHandler handles admin responses to post reports
func AdminRespondReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// Get admin ID
	adminIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting admin ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	adminID, err := strconv.Atoi(adminIDStr)
	if err != nil {
		log.Printf("Error converting admin ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var requestData struct {
		ReportID int    `json:"reportId"`
		Status   string `json:"status"`
		Response string `json:"response"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if requestData.Status != "approved" && requestData.Status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	// If approved, delete the post
	if requestData.Status == "approved" {
		// Get the post ID from the report
		var postID int
		err = db.QueryRow("SELECT PostID FROM PostReport WHERE ReportID = ?", requestData.ReportID).Scan(&postID)
		if err != nil {
			log.Printf("Error getting post ID from report: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Delete the post (this will cascade to comments, likes, etc.)
		_, err = db.Exec("DELETE FROM Post WHERE PostID = ?", postID)
		if err != nil {
			log.Printf("Error deleting post: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Update the report status
	_, err = db.Exec(
		"UPDATE PostReport SET Status = ?, AdminResponse = ?, AdminID = ?, ResponseDate = ? WHERE ReportID = ?",
		requestData.Status, requestData.Response, adminID, time.Now(), requestData.ReportID,
	)
	if err != nil {
		log.Printf("Error updating report: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Report response saved successfully",
	})
}

// isModerator checks if the user is a moderator or admin
func isModerator(r *http.Request) bool {
	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		return false
	}
	defer db.Close()

	userIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		return false
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return false
	}

	var privilege int
	err = db.QueryRow("SELECT privilege FROM User WHERE UserID = ?", userID).Scan(&privilege)
	if err != nil {
		return false
	}

	return privilege >= 2 // Moderator (2) or Admin (3)
}

// UserReportsHandler returns reports made by the current user
func UserReportsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	userIDStr, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Error converting user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	query := `
		SELECT
			pr.ReportID, pr.PostID,
			COALESCE(p.title, '[Post Deleted]') as title,
			COALESCE(p.content, '[Post content no longer available]') as content,
			COALESCE(u1.username, '[Unknown User]') as post_author,
			pr.ReportDate, pr.Reason, pr.Status, pr.AdminResponse, pr.ResponseDate
		FROM PostReport pr
		LEFT JOIN Post p ON pr.PostID = p.PostID
		LEFT JOIN User u1 ON p.UserID = u1.UserID
		WHERE pr.ModeratorID = ?
		ORDER BY pr.ReportDate DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying user reports: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type UserReport struct {
		ReportID      int     `json:"ReportID"`
		PostID        int     `json:"PostID"`
		PostTitle     string  `json:"PostTitle"`
		PostContent   string  `json:"PostContent"`
		PostAuthor    string  `json:"PostAuthor"`
		ReportDate    string  `json:"ReportDate"`
		Reason        string  `json:"Reason"`
		Status        string  `json:"Status"`
		AdminResponse string  `json:"AdminResponse"`
		ResponseDate  *string `json:"ResponseDate"`
	}

	var reports []UserReport
	for rows.Next() {
		var report UserReport
		var responseDate sql.NullString
		var adminResponse sql.NullString

		err := rows.Scan(
			&report.ReportID, &report.PostID, &report.PostTitle, &report.PostContent, &report.PostAuthor,
			&report.ReportDate, &report.Reason, &report.Status, &adminResponse, &responseDate,
		)
		if err != nil {
			log.Printf("Error scanning user report: %v", err)
			continue
		}

		// Handle nullable fields
		if adminResponse.Valid {
			report.AdminResponse = adminResponse.String
		} else {
			report.AdminResponse = ""
		}

		if responseDate.Valid {
			report.ResponseDate = &responseDate.String
		}

		reports = append(reports, report)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
