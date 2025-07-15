package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Profile struct {
	UserID       int    `json:"UserID"`
	CreatedPosts []Post `json:"CreatedPosts"`
	// UserComments  []CommentedPost `json:"UserComments"`
	LikedPosts    []Post `json:"LikedPosts"`
	DislikedPosts []Post `json:"DislikedPosts"`
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "meow.db")
	if err != nil {
		panic(err)
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Profile Handler Called")
	fmt.Printf("Request Method: %s\n", r.Method)

	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed:", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromSession(r)
	fmt.Printf("User ID from session: %d\n", userID)

	profile := Profile{
		UserID:       userID,
		CreatedPosts: getCreatedPosts(userID),
		// UserComments:  getUserComments(userID),
		LikedPosts:    getLikedPosts(userID),
		DislikedPosts: getDislikedPosts(userID),
	}

	// Add debug logs here
	// log.Printf("Created Posts: %+v", profile.CreatedPosts)
	// log.Printf("User Comments: %+v", profile.UserComments)
	// log.Printf("Liked Posts: %+v", profile.LikedPosts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func getCreatedPosts(userID int) []Post {
	query := `
		SELECT 
            p.PostID, 
            p.UserID,
            p.PostDate,
            p.title,
            p.content,
			p.ImagePath,
            u.username,
            COALESCE(pl.likes, 0) AS likes,
            COALESCE(pdl.dislike, 0) AS dislikes,
            COALESCE(cmt.comments, 0) AS comments
        FROM 
            Post p
        JOIN 
            User u ON p.UserID = u.UserID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS likes FROM PostLike GROUP BY PostID
        ) AS pl ON p.PostID = pl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS dislike FROM PostDislike GROUP BY PostID
        ) AS pdl ON p.PostID = pdl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS comments FROM Comment GROUP BY PostID
        ) AS cmt ON p.PostID = cmt.PostID
        WHERE 
			p.UserID = ?
        ORDER BY 
			p.PostDate DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostID, &post.UserID, &post.PostDate, &post.Title, &post.Content, &post.ImagePath, &post.Username,
			&post.Likes, &post.Dislikes, &post.CmtCount)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts
}

// todo: will make it after
// func getUserComments(userID int) []CommentedPost {
// 	fmt.Printf("Fetching liked posts for user %d\n", userID) //debug
// 	query := `
//         SELECT c.CommentID, c.PostID, c.content, c.created_at
//         FROM Comment c
//         JOIN Post p ON c.PostID = p.PostID
//         WHERE c.UserID = ?
//     `
// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil
// 	}
// 	defer rows.Close()

// 	var comments []CommentedPost
// 	for rows.Next() {
// 		var comment CommentedPost
// 		rows.Scan(&comment.CommentID, &comment.PostID, &comment.Comment, &comment.CreateDate)
// 		comments = append(comments, comment)
// 	}
// 	return comments
// }

func getLikedPosts(userID int) []Post {
	fmt.Printf("Fetching disliked posts for user %d\n", userID)
	query := `
        SELECT 
			p.PostID, 
			p.UserID,
			p.PostDate,
			p.title,
			p.content,
			p.ImagePath,
			u.username,
			COALESCE(pl.likes, 0) AS likes,
			COALESCE(pdl.dislike, 0) AS dislikes,
			COALESCE(cmt.comments, 0) AS comments
		FROM 
			Post p
		JOIN 
			User u ON p.UserID = u.UserID
		JOIN 
			PostLike l ON p.PostID = l.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS likes FROM PostLike GROUP BY PostID
		) AS pl ON p.PostID = pl.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS dislike FROM PostDislike GROUP BY PostID
		) AS pdl ON p.PostID = pdl.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS comments FROM Comment GROUP BY PostID
		) AS cmt ON p.PostID = cmt.PostID
		WHERE 
			l.UserID = ?
		ORDER BY 
			p.PostDate DESC

    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostID, &post.UserID, &post.PostDate, &post.Title, &post.Content, &post.ImagePath, &post.Username,
			&post.Likes, &post.Dislikes, &post.CmtCount)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts
}

func getDislikedPosts(userID int) []Post {
	query := `
        SELECT 
			p.PostID, 
			p.UserID,
			p.PostDate,
			p.title,
			p.content,
			p.ImagePath,
			u.username,
			COALESCE(pl.likes, 0) AS likes,
			COALESCE(pdl.dislike, 0) AS dislikes,
			COALESCE(cmt.comments, 0) AS comments
		FROM 
			Post p
		JOIN 
			User u ON p.UserID = u.UserID
		JOIN 
			PostDislike d ON p.PostID = d.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS likes FROM PostLike GROUP BY PostID
		) AS pl ON p.PostID = pl.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS dislike FROM PostDislike GROUP BY PostID
		) AS pdl ON p.PostID = pdl.PostID
		LEFT JOIN (
			SELECT PostID, COUNT(*) AS comments FROM Comment GROUP BY PostID
		) AS cmt ON p.PostID = cmt.PostID
		WHERE 
			d.UserID = ?
		ORDER BY 
			p.PostDate DESC

    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.PostID, &post.UserID, &post.PostDate, &post.Title, &post.Content, &post.ImagePath, &post.Username,
			&post.Likes, &post.Dislikes, &post.CmtCount)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts
}

func getUserIDFromSession(r *http.Request) int {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Printf("Cookie error: %v\n", err)
		return 0
	}

	fmt.Printf("Session token found: %s\n", cookie.Value)

	query := `SELECT user_id FROM Session WHERE session_id = ?`
	var userID int
	err = db.QueryRow(query, cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return 0
	}

	fmt.Printf("Found user ID: %d\n", userID)
	return userID
}
