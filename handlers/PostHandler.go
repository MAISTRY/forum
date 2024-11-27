package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var (
	MultipostQuery = `
        SELECT 
            p.PostID, 
            p.UserID,
            p.PostDate,
            p.title,
            p.content,
			p.ImagePath,
            u.username,
            COALESCE(pl.likes, 0) AS likes,
            COALESCE(pdl.dislike, 0) AS dislikes,
            COALESCE(cmt.comments, 0) AS comments
        FROM 
            Post p
        JOIN 
            User u ON p.UserID = u.UserID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS likes FROM PostLike GROUP BY PostID
        ) AS pl ON p.PostID = pl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS dislike FROM PostDislike GROUP BY PostID
        ) AS pdl ON p.PostID = pdl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS comments FROM Comment GROUP BY PostID
        ) AS cmt ON p.PostID = cmt.PostID
        ORDER BY 
			p.PostDate DESC
    `
	categoryQuery = `
        SELECT c.title
        FROM Category c
        JOIN PostCategory pc ON c.CategoryID = pc.CategoryID
        WHERE pc.PostID = ?
    `
)

// PostHandler handles HTTP requests for retrieving post information.
// It can fetch either a single post or multiple posts based on the presence of a 'postid' query parameter.
//
// Parameters:
//   - w: http.ResponseWriter - Used to write the HTTP response.
//   - r: *http.Request - The incoming HTTP request containing query parameters and other request details.
//
// The function does not return any value directly, but writes the response to the http.ResponseWriter:
//   - For a single post (when 'postid' is provided): Responds with a JSON object containing details of the requested post.
//   - For multiple posts (when 'postid' is not provided): Responds with a JSON array of post objects.
//
// In case of errors, appropriate HTTP error statuses and messages are written to the response.
func PostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	postRows, err := db.Query(MultipostQuery)
	if err != nil {
		http.Error(w, "Error querying posts", http.StatusInternalServerError)
		return
	}
	defer postRows.Close()

	var posts []Post
	for postRows.Next() {
		var post Post
		if err := postRows.Scan(
			&post.PostID, &post.UserID, &post.PostDate, &post.Title, &post.Content, &post.ImagePath, &post.Username,
			&post.Likes, &post.Dislikes, &post.CmtCount,
		); err != nil {
			http.Error(w, "Error scanning post details", http.StatusInternalServerError)
			return
		}

		categoryRows, err := db.Query(categoryQuery, post.PostID)
		if err != nil {
			http.Error(w, "Error querying categories", http.StatusInternalServerError)
			return
		}
		defer categoryRows.Close()

		categories := []string{}
		for categoryRows.Next() {
			var category string
			if err := categoryRows.Scan(&category); err != nil {
				http.Error(w, "Error scanning category", http.StatusInternalServerError)
				return
			}
			categories = append(categories, category)
		}
		post.Categories = categories
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}
