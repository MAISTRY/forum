package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type UserComment struct {
	CommentID   int    `json:"comment_id"`
	PostID      int    `json:"post_id"`
	PostTitle   string `json:"post_title"`
	PostAuthor  string `json:"post_author"`
	Comment     string `json:"comment"`
	CommentDate string `json:"comment_date"`
	Likes       int    `json:"likes"`
	Dislikes    int    `json:"dislikes"`
}

type ActivityData struct {
	UserComments []UserComment `json:"user_comments"`
}

func ActivityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromSession(r)
	if userID == 0 {
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

	userComments, err := getUserComments(db, userID)
	if err != nil {
		log.Printf("Error getting user comments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	activityData := ActivityData{
		UserComments: userComments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activityData)
}

func getUserComments(db *sql.DB, userID int) ([]UserComment, error) {
	query := `
		SELECT 
			c.CommentID,
			c.PostID,
			p.title as post_title,
			pu.username as post_author,
			c.content as comment,
			c.CmtDate as comment_date,
			COALESCE(cl.likes, 0) as likes,
			COALESCE(cd.dislikes, 0) as dislikes
		FROM Comment c
		JOIN Post p ON c.PostID = p.PostID
		JOIN User pu ON p.UserID = pu.UserID
		LEFT JOIN (
			SELECT CommentID, COUNT(*) as likes 
			FROM CommentLike 
			GROUP BY CommentID
		) cl ON c.CommentID = cl.CommentID
		LEFT JOIN (
			SELECT CommentID, COUNT(*) as dislikes 
			FROM CommentDislike 
			GROUP BY CommentID
		) cd ON c.CommentID = cd.CommentID
		WHERE c.UserID = ?
		ORDER BY c.CmtDate DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying user comments: %v", err)
	}
	defer rows.Close()

	var comments []UserComment
	for rows.Next() {
		var comment UserComment
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.PostTitle,
			&comment.PostAuthor,
			&comment.Comment,
			&comment.CommentDate,
			&comment.Likes,
			&comment.Dislikes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %v", err)
	}

	return comments, nil
}
