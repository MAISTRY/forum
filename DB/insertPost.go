package DB

import (
	"database/sql"
	"fmt"
)

// InsertPost inserts a new post into the database, including the associated categories.
// It takes the following parameters:
// - db: a pointer to an *sql.DB instance representing the database connection
// - title: the title of the post
// - content: the content of the post
// - imagePath: the path to the image associated with the post
// - categories: a slice of strings representing the categories associated with the post
// - usrID: the ID of the user who created the post
// It returns an error if any part of the operation fails.
func InsertPost(db *sql.DB, title, content, imagePath string, categories []string, usrID int) error {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmtPost, err := tx.Prepare("INSERT INTO Post (UserID, title, content, ImagePath) VALUES (?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmtPost.Close()

	result, err := stmtPost.Exec(usrID, title, content, imagePath)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting post: %v", err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	stmtPostCategory, err := tx.Prepare("INSERT INTO PostCategory (PostID, CategoryID) VALUES (?,?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmtPostCategory.Close()

	for _, category := range categories {
		var categoryID int
		err = tx.QueryRow("SELECT CategoryID FROM Category WHERE title = ?", category).Scan(&categoryID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error getting category ID: %v", err)
		}

		_, err = stmtPostCategory.Exec(postID, categoryID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting into PostCategory: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
