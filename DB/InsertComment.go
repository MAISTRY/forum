package DB

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	insertCommentQuery = `
        INSERT INTO "Comment" (PostID, UserID, content)
        VALUES (?,?,?)
    `
	SelectUsernameQuery = `
		SELECT username 
		FROM User 
		WHERE UserID = ?
	`
)

func InsertComment(postID int, userID int, comment string) (int64, error) {

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		return -1, fmt.Errorf("error Opening database: %v", err)
	}
	defer db.Close()

	result, err := db.Exec(insertCommentQuery, postID, userID, comment)
	if err != nil {
		return -1, fmt.Errorf("error insert in the database: %v", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return commentID, nil
}
