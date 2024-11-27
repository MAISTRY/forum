package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var (
	commentsQuery = `
        SELECT
			cm.CommentID,
            cm.content,
			cm.CmtDate,
            u.username,
            COALESCE(cl.CommentLikes, 0) AS likes,
            COALESCE(cd.CommentDislikes, 0) AS dislikes
        FROM
            Comment cm
        JOIN
            User u ON cm.UserID = u.UserID
        LEFT JOIN (
            SELECT CommentID, COUNT(*) AS CommentLikes FROM CommentLike GROUP BY CommentID
        ) AS cl ON cm.CommentID = cl.CommentID
        LEFT JOIN (
            SELECT CommentID, COUNT(*) AS CommentDislikes FROM CommentDislike GROUP BY CommentID
			) AS cd ON cm.CommentID = cd.CommentID
			WHERE
            cm.PostID = ?
			`
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	postIDParam := r.URL.Query().Get("postid")
	if postIDParam == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDParam)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	var count int
	checkQuery := `SELECT COUNT(*) FROM Post WHERE PostID = ?`
	checkError := db.QueryRow(checkQuery, postID).Scan(&count)
	if checkError != nil {
		http.Error(w, "Error checking post validity", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "Post not found or deleted", http.StatusNotFound)
		return
	}

	commentRows, err := db.Query(commentsQuery, postID)
	if err != nil {
		http.Error(w, "Error querying comments", http.StatusInternalServerError)
		return
	}
	defer commentRows.Close()

	comments := []Comment{}
	for commentRows.Next() {
		var cmt Comment
		if err := commentRows.Scan(&cmt.CmtID, &cmt.Content, &cmt.CmtDate, &cmt.Username, &cmt.Likes, &cmt.Dislikes); err != nil {
			http.Error(w, "Error scanning comments", http.StatusInternalServerError)
			return
		}
		comments = append(comments, cmt)
	}

	// if len(comments) == 0 {
	// 	http.Error(w, "No comments found", http.StatusNotFound)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
