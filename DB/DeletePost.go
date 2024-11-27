package DB

import (
	"database/sql"
	"fmt"
)

const (
	deletePostQuery = `
        DELETE FROM Post
        WHERE PostID =?;
    `

	deletePostCategoryQuery = `
        DELETE FROM PostCategory
        WHERE PostID =?;
    `

	deletePostLikeQuery = `
        DELETE FROM PostLike
        WHERE PostID =?;
    `

	deletePostDislikeQuery = `
        DELETE FROM PostDislike
        WHERE PostID =?;
    `

	deleteCommentQuery = `
        DELETE FROM Comment
        WHERE PostID =?;
    `
)

// DelPost deletes a post from the database, including all associated comments, likes, dislikes, and categories.
// It takes a database connection and the ID of the post to be deleted, and returns an error if any part of the deletion fails.
// The function uses a transaction to ensure that the entire deletion process is atomic - either all changes are committed or none are.
func DelPost(db *sql.DB, postID string) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error starting transaction: %v", err)
	}

	if _, err = tx.Exec(deleteCommentQuery, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting comments: %v", err)
	}

	if _, err = tx.Exec(deletePostLikeQuery, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting likes: %v", err)
	}

	if _, err = tx.Exec(deletePostDislikeQuery, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting dislikes: %v", err)
	}

	if _, err = tx.Exec(deletePostCategoryQuery, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting categories: %v", err)
	}

	if _, err = tx.Exec(deletePostQuery, postID); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting post: %v", err)
	}

	return nil

}
