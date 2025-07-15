package DB

import (
	"database/sql"
	"log"
	// 	"fmt"
)

const CategoryInsertion = `INSERT INTO Category (title, description) VALUES (?, ?)`

// func InsertCategory(db *sql.DB, title, description, usrID int) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error starting transaction: %v", err)
// 	}

// 	stmtPost, err := tx.Prepare("INSERT INTO Category (title, description) VALUES (?,?)")
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error preparing statement: %v", err)
// 	}
// 	defer stmtPost.Close()

// 	result, err := stmtPost.Exec(title, description)
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error inserting post: %v", err)
// 	}

// 	postID, err := result.LastInsertId()
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error getting last insert ID: %v", err)
// 	}

// 	stmtPostCategory, err := tx.Prepare("INSERT INTO PostCategory (PostID, CategoryID) VALUES (?,?)")
// 	if err != nil {
// 		tx.Rollback()
// 		return fmt.Errorf("error preparing statement: %v", err)
// 	}
// 	defer stmtPostCategory.Close()

// 	for _, category := range categories {
// 		var categoryID int
// 		err = tx.QueryRow("SELECT CategoryID FROM Category WHERE title = ?", category).Scan(&categoryID)
// 		if err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("error getting category ID: %v", err)
// 		}

// 		_, err = stmtPostCategory.Exec(postID, categoryID)
// 		if err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("error inserting into PostCategory: %v", err)
// 		}
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return fmt.Errorf("error committing transaction: %v", err)
// 	}

// 	return nil
// }

func InitailTableFiller(db *sql.DB) {
	if _, err := db.Exec(CategoryInsertion); err != nil {
		log.Fatalf("error inserting data into the category table: %v", err)
	}
}
