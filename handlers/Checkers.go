package handlers

import (
	"database/sql"
	"fmt"
)

const (
	insertLikeQuery = `
        INSERT INTO PostLike (UserID, PostID)
        VALUES (?,?);
    `
	insertDislikeQuery = `
	    INSERT INTO PostDislike (UserID, PostID)
		VALUES (?,?);
	`
	deleteLikeQuery = `
		DELETE FROM PostLike
		WHERE UserID = ? AND PostID = ?;
	`
	deleteDisLikeQuery = `
		DELETE FROM PostDislike
		WHERE UserID = ? AND PostID = ?;
	`
	getLikeQuery = `
	    SELECT FROM COUNT(*) FROM PostLike
		WHERE UserID =? AND PostID =?;
	`
	getDislikeQuery = `
	    SELECT FROM COUNT(*) FROM PostDislike
		WHERE UserID =? AND PostID =?;
	`
	//! for comments:
	deleteCommentDislikeQuery = `
	    DELETE FROM CommentDislike
		WHERE UserID =? AND CommentID =?;
	`
	deleteCommentLikeQuery = `
	    DELETE FROM CommentLike
		WHERE UserID =? AND CommentID =?;
	`
	insertCommentDislikeQuery = `
	    INSERT INTO CommentDislike (UserID, CommentID)
		VALUES (?,?);
	`
	insertCommentLikeQuery = `
	    INSERT INTO CommentLike (UserID, CommentID)
		VALUES (?,?);
	`
	getLikesCountQuery = `
	    SELECT COUNT(*) FROM PostLike WHERE PostID =?;
	`
	getDislikesCountQuery = `
	    SELECT COUNT(*) FROM PostDislike WHERE PostID =?;
	`
)

// IsLiked checks if a user has liked a specific post.
//
// Parameters:
//   - db: A pointer to the SQL database connection.
//   - userID: The ID of the user checking for the like.
//   - postID: The ID of the post being checked.
//
// Returns:
//   - bool: True if the user has liked the post, false otherwise.
//   - error: An error if the database query fails, nil otherwise.
func IsLiked(db *sql.DB, userID, postID string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM PostLike WHERE UserID = ? AND PostID = ?", userID, postID).Scan(&count)
    if err != nil {
        fmt.Printf("Error checking if user has liked post %v\n", err)
        return false, err
    }
    return count > 0, nil
}

// IsDisliked checks if the given user has disliked the specified post.
//
// Parameters:
//   - db: A pointer to the SQL database connection.
//   - userID: The ID of the user checking for the dislike.
//   - postID: The ID of the post being checked.
//
// Returns:
//   - bool: True if the user has disliked the post, false otherwise.
//   - error: An error if the database query fails, nil otherwise.
func IsDisliked(db *sql.DB, userID, postID string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM PostDislike WHERE UserID = ? AND PostID = ?", userID, postID).Scan(&count)
    if err != nil {
        fmt.Printf("Error checking if user has disliked post %v\n", err)
        return false, err
    }
    return count > 0, nil
}

// IsCommentDisliked checks if a user has disliked a specific comment.
//
// Parameters:
//   - db: A pointer to the SQL database connection.
//   - userID: The ID of the user checking for the dislike.
//   - commentID: The ID of the comment being checked.
//
// Returns:
//   - bool: True if the user has disliked the comment, false otherwise.
//   - error: An error if the database query fails, nil otherwise.
func IsCommentDisliked(db *sql.DB, userID, commentID string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM CommentDislike WHERE UserID =? AND CommentID =?", userID, commentID).Scan(&count)
    if err != nil {
        fmt.Printf("Error checking if user has disliked comment %v\n", err)
        return false, err
    }
    return count > 0, nil
}

// IsCommentLiked checks if a user has liked a specific comment.
//
// Parameters:
//   - db: A pointer to the SQL database connection.
//   - userID: The ID of the user checking for the like.
//   - commentID: The ID of the comment being checked.
//
// Returns:
//   - bool: True if the user has liked the comment, false otherwise.
//   - error: An error if the database query fails, nil otherwise.
func IsCommentLiked(db *sql.DB, userID, commentID string) (bool, error) {
    var count int
    err := db.QueryRow("SELECT COUNT(*) FROM CommentLike WHERE UserID =? AND CommentID =?", userID, commentID).Scan(&count)
    if err != nil {
        fmt.Printf("Error checking if user has liked comment %v\n", err)
        return false, err
    }
    return count > 0, nil
}
