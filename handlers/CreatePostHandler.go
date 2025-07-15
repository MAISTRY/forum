package handlers

import (
	"database/sql"
	"fmt"
	"forum/DB"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	uploadDir = "static/uploads/"
	dataDir   = "../uploads/"
)

// CreatePostHandler handles the creation of a new post in the forum.
// It processes the form data, including title, content, and an optional image,
// and inserts the post into the database.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response.
//   - r *http.Request: The HTTP request containing the form data for the new post.
//
// The function does not return any value, but writes a JSON response to the
// http.ResponseWriter indicating the success or failure of the post creation.
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, `{"success": false, "message": "Database connection error"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Error getting user ID"}`, http.StatusInternalServerError)
		return
	}

	title := (r.FormValue("title"))
	content := (r.FormValue("content"))
	categoriesFromForm := r.Form["categories"]

	if title == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success": false, "message": "Post title is required"}`, http.StatusBadRequest)
		return
	}
	if content == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success": false, "message": "Post content is required"}`, http.StatusBadRequest)
		return
	}
	if len(categoriesFromForm) == 0 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success": false, "message": "Please select at least one category"}`, http.StatusBadRequest)
		return
	}

	// Handle duplicate post titles by adding a number
	originalTitle := title
	i := 1
	for {
		var existingCount int
		duplicateQuery := `SELECT COUNT(*) FROM Post WHERE UserID = ? AND title = ?`
		err = db.QueryRow(duplicateQuery, userID, title).Scan(&existingCount)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Error checking for duplicate titles"}`, http.StatusInternalServerError)
			return
		}

		if existingCount == 0 {
			break // Title is unique, we can use it
		}

		// Title exists, add/increment number
		title = fmt.Sprintf("%s %d", originalTitle, i)
		i++
	}

	var imagePath string
	file, fileHead, err := r.FormFile("image")
	if err != nil {
		// No image uploaded, this is fine
		fmt.Printf("Post without image\n")
	}
	if file != nil {
		defer file.Close()

		// Validate file size (max 10MB)
		const maxFileSize = 10 << 20 // 10MB
		if fileHead.Size > maxFileSize {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Image file too large. Maximum size is 10MB."}`, http.StatusBadRequest)
			return
		}

		// Validate file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
		}

		// Read first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Error reading image file"}`, http.StatusInternalServerError)
			return
		}

		contentType := http.DetectContentType(buffer)
		if !allowedTypes[contentType] {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Invalid file type. Only JPEG and PNG images are allowed."}`, http.StatusBadRequest)
			return
		}

		// Reset file pointer to beginning
		_, err = file.Seek(0, 0)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Error reading file"}`, http.StatusInternalServerError)
			return
		}

		filename := filepath.Base(fileHead.Filename)
		// Remove any path separators from filename for security
		filename = filepath.Base(filename)

		// If filename is empty, generate one
		if filename == "" || filename == "." {
			filename = fmt.Sprintf("upload_%d.jpg", time.Now().Unix())
		}

		storePath := filepath.Join(uploadDir, filename)
		imagePath = filepath.Join(dataDir, filename)

		// Handle duplicate filenames by adding a number prefix
		i := 1
		for {
			if _, err := os.Stat(storePath); os.IsNotExist(err) {
				break // File doesn't exist, we can use this name
			}
			newFilename := fmt.Sprintf("%d_%s", i, filename)
			storePath = filepath.Join(uploadDir, newFilename)
			imagePath = filepath.Join(dataDir, newFilename)
			i++
		}

		// Create the file
		filePlace, err := os.Create(storePath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Error saving image"}`, http.StatusInternalServerError)
			return
		}
		defer filePlace.Close()

		// Copy file content
		_, err = io.Copy(filePlace, file)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"success": false, "message": "Error writing file"}`, http.StatusInternalServerError)
			return
		}

		fmt.Printf("Image saved successfully: %s\n", storePath)
	}

	UsrID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, `{"success": false, "message": "Error converting user id"}`, http.StatusInternalServerError)
		return
	}

	err = DB.InsertPost(db, title, content, imagePath, categoriesFromForm, UsrID)
	if err != nil {
		fmt.Printf("Error inserting post: %v", err)
		http.Error(w, `{"success": false, "message": "Error inserting post"}`, http.StatusInternalServerError)
		return
	}

	// PostTable, err := db.Exec(insertPostQuery, userID, title, content, imagePath)
	// if err != nil {
	// 	http.Error(w, `{"success": false, "message": "Error querying posts"}`, http.StatusInternalServerError)
	// 	return
	// }
	// PostID, err := PostTable.LastInsertId()
	// if err != nil {
	// 	http.Error(w, `{"success": false, "message": "Error getting post id"}`, http.StatusInternalServerError)
	// 	return
	// }

	// for _, categoryID := range r.Form["category"] {
	// 	_, err := db.Exec(insertPostCategoryQuery, PostID, categoryID)
	// 	if err != nil {
	// 		http.Error(w, `{"success": false, "message": "Error inserting post category"}`, http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	w.Write([]byte(`Post created successfully`))

	w.Header().Set("HX-Redirect", "/")
	fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="0;url=/home"></head></html>`)

}
