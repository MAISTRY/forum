package DB

import "database/sql"

// GetLikeCount retrieves the number of likes for the specified post ID from the database.
// It takes a database connection and a post ID as input, and returns the count of likes
// and any error that occurred during the query.
func GetLikeCount(db *sql.DB, postID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM PostLike WHERE PostID =?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetDislikeCount retrieves the number of dislikes for the specified post ID from the database.
// It takes a database connection and a post ID as input, and returns the count of dislikes
// and any error that occurred during the query.
func GetDislikeCount(db *sql.DB, postID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM PostDislike WHERE PostID =?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentLikeCount retrieves the number of likes for the specified comment ID from the database.
// It takes a database connection and a comment ID as input, and returns the count of likes
// and any error that occurred during the query.
func GetCommentLikeCount(db *sql.DB, commentID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM CommentLike WHERE CommentID =?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentDislikeCount retrieves the number of dislikes for the specified comment ID from the database.
// It takes a database connection and a comment ID as input, and returns the count of dislikes
// and any error that occurred during the query.
func GetCommentDislikeCount(db *sql.DB, commentID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM CommentDislike WHERE CommentID =?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentsUsername retrieves the username for the specified user ID from the database.
// It takes a database connection and a user ID as input, and returns the username
// and any error that occurred during the query.
func GetCommentsUsername(db *sql.DB, userID string) (string, error) {
	var username string
	err := db.QueryRow(SelectUsernameQuery, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
