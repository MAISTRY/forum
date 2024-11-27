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
		http.Error(w, `The Title isn't present!!!`, http.StatusOK)
		return
	}
	if content == "" {
		http.Error(w, `The Content isn't present!!!`, http.StatusOK)
		return
	}
	if len(categoriesFromForm) == 0 {
		http.Error(w, `The Categories aren't present!!!`, http.StatusOK)
		return
	}

	var imagePath string
	file, fileHead, err := r.FormFile("image")
	if err != nil {
		fmt.Printf(`post without image!`)
	}
	if file != nil {
		filename := filepath.Base(fileHead.Filename)
		storePath := filepath.Join(uploadDir, filename)
		imagePath = filepath.Join(dataDir, filename)

		i := 1
		for _, err := os.Stat(storePath); err == nil; i++ {
			storePath = filepath.Join(uploadDir, fmt.Sprintf("%d_%s", i, filename))
			imagePath = filepath.Join(imagePath, fmt.Sprintf("%d_%s", i, filename))
		}

		filePlace, err := os.Create(storePath)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error saving image"}`, http.StatusInternalServerError)
			return
		}
		defer filePlace.Close()

		_, err = file.Seek(0, 0)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error reading file"}`, http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(filePlace, file)
		if err != nil {
			http.Error(w, `{"success": false, "message": "Error writing file"}`, http.StatusInternalServerError)
			return
		}

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
