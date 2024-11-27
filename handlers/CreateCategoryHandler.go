package handlers

// import (
// 	"database/sql"
// 	"net/http"
// )

var (
	AllcategoryQuery = `
        SELECT title, description 
		FROM Category
    `
)

// func CategoryHandler(w http.ResponseWriter, r *http.Request) {

// 	db, err := sql.Open("sqlite3", "meow.db")
// 	if err != nil {
// 		http.Error(w, `{"success": false, "message": "Database connection error"}`, http.StatusInternalServerError)
// 		return
// 	}
// 	defer db.Close()

// 	userID, err := getUserIDByCookie(r, db)
// 	if err != nil {
// 		http.Error(w, `{"success": false, "message": "Error getting user ID"}`, http.StatusInternalServerError)
// 		return
// 	}

// }
