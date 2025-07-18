package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// CategoryRequest represents a request to add a new category
type CategoryRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CategoryDeleteRequest represents a request to delete a category
type CategoryDeleteRequest struct {
	CategoryID int `json:"categoryId"`
}

// AdminCategory represents a category in the admin interface
type AdminCategory struct {
	CategoryID  int    `json:"CategoryID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"UserID"`
}

// AdminCategoriesHandler returns all categories for admin management
func AdminCategoriesHandler(w http.ResponseWriter, r *http.Request) {
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
		SELECT CategoryID, title, description, UserID
		FROM Category
		ORDER BY title
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []AdminCategory
	for rows.Next() {
		var category AdminCategory

		err := rows.Scan(&category.CategoryID, &category.Title, &category.Description,
			&category.UserID)
		if err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}

		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// AdminAddCategoryHandler adds a new category
func AdminAddCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Title == "" {
		http.Error(w, "Category title is required", http.StatusBadRequest)
		return
	}

	if req.Description == "" {
		http.Error(w, "Category description is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get admin user ID
	adminUserID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting admin user ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if category already exists
	var existingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM Category WHERE title = ?", req.Title).Scan(&existingCount)
	if err != nil {
		log.Printf("Error checking existing category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingCount > 0 {
		http.Error(w, "Category with this title already exists", http.StatusBadRequest)
		return
	}

	// Insert new category
	_, err = db.Exec("INSERT INTO Category (title, description, UserID) VALUES (?, ?, ?)",
		req.Title, req.Description, adminUserID)
	if err != nil {
		log.Printf("Error inserting category: %v", err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Category created successfully",
	})
}

// AdminDeleteCategoryHandler deletes a category
func AdminDeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is admin
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CategoryDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CategoryID <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete PostCategory relationships first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM PostCategory WHERE CategoryID = ?", req.CategoryID)
	if err != nil {
		log.Printf("Error deleting post category relationships: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Delete the category
	result, err := tx.Exec("DELETE FROM Category WHERE CategoryID = ?", req.CategoryID)
	if err != nil {
		log.Printf("Error deleting category: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if category was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
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
		"message": "Category deleted successfully",
	})
}

// PublicCategoriesHandler returns all categories for public use (no admin required)
func PublicCategoriesHandler(w http.ResponseWriter, r *http.Request) {
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

	query := `
		SELECT CategoryID, title, description
		FROM Category
		ORDER BY title
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying categories: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []struct {
		CategoryID  int    `json:"CategoryID"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	for rows.Next() {
		var category struct {
			CategoryID  int    `json:"CategoryID"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		err := rows.Scan(&category.CategoryID, &category.Title, &category.Description)
		if err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}

		categories = append(categories, category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
