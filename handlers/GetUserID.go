package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

const (
	getUserIDQuery = `SELECT user_id FROM Session WHERE session_id = ?`
)

// getUserIDByCookie retrieves the user ID associated with a session cookie.
//
// Parameters:
//   - r: An http.Request object containing the client's request information,
//     including cookies.
//   - db: A pointer to sql.DB representing the database connection.
//
// Returns:
//   - A string containing the user ID if successful.
//   - An error if any step in the process fails (e.g., cookie retrieval,
//     database query preparation, or query execution).
func getUserIDByCookie(r *http.Request, db *sql.DB) (string, error) {
	yummyCookie, err := r.Cookie("sessionID")
	if err != nil {
		return "", fmt.Errorf("error getting cookie: %v", err)
	}

	stmt, err := db.Prepare(getUserIDQuery)
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	var userID string

	err = stmt.QueryRow(yummyCookie.Value).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("error getting user ID: %v", err)
	}

	return userID, err
}

func GetPostOwnerID(PostID string, db *sql.DB) (string, error) {
	stmt, err := db.Prepare(`SELECT UserID FROM Post WHERE PostID =?;`)
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	var userID string
	err = stmt.QueryRow(PostID).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("post not found")
	} else if err != nil {
		return "", fmt.Errorf("error getting user ID: %v", err)
	}

	return userID, nil
}

func GetCommentOwnerID(CommentID string, db *sql.DB) (string, error) {
	stmt, err := db.Prepare(`SELECT UserID FROM Comment WHERE CommentID =?;`)
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	var userID string
	err = stmt.QueryRow(CommentID).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("comment not found")
	} else if err != nil {
		return "", fmt.Errorf("error getting user ID: %v", err)
	}

	return userID, nil
}
