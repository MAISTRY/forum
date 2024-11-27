package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type LikedPost struct {
	PostID     int    `json:"PostID"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	CreateDate string `json:"CreateDate"`
	Author     string `json:"Author"`
}

type DislikedPost struct {
	PostID     int    `json:"PostID"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	CreateDate string `json:"CreateDate"`
	Author     string `json:"Author"`
}

type CreatedPost struct {
	PostID     int    `json:"PostID"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	CreateDate string `json:"CreateDate"`
	Likes      int    `json:"Likes"`
	Dislikes   int    `json:"Dislikes"`
	Comments   int    `json:"Comments"`
}

type Profile struct {
	UserID        int             `json:"UserID"`
	Username      string          `json:"Username"`
	CreatedPosts  []CreatedPost   `json:"CreatedPosts"`
	UserComments  []CommentedPost `json:"UserComments"`
	LikedPosts    []LikedPost     `json:"LikedPosts"`
	DislikedPosts []DislikedPost  `json:"DislikedPosts"`
}

// const (
// 	createdPostQuery = `
// 		SELECT user_id FROM Session WHERE session_id = ?
// 	`
// 	CommentsQuery = `
// 		SELECT user_id FROM Session WHERE session_id = ?
// 	`
// 	likedPostQuery = `
// 		SELECT user_id FROM Session WHERE session_id = ?
// 	`
// 	DislikedPostQuery = `
// 		SELECT user_id FROM Session WHERE session_id = ?
// 	`
// )

// type (
//
//	createdPost struct {
//	}
//	Comments struct {
//	}
//	likedPost struct {
//	}
//	DislikedPost struct {
//	}
//	Profile struct {
//	}
//
// )

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
		UserID:        userID,
		CreatedPosts:  getCreatedPosts(userID),
		UserComments:  getUserComments(userID),
		LikedPosts:    getLikedPosts(userID),
		DislikedPosts: getDislikedPosts(userID),
	}

	// Add debug logs here
	log.Printf("Created Posts: %+v", profile.CreatedPosts)
	log.Printf("User Comments: %+v", profile.UserComments)
	log.Printf("Liked Posts: %+v", profile.LikedPosts)

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func getCreatedPosts(userID int) []CreatedPost {
	query := `
    SELECT p.PostID, p.title, p.content, p.PostDate,
    COUNT(DISTINCT l.PostID) as likes,
    COUNT(DISTINCT d.PostID) as dislikes,
    COUNT(DISTINCT c.CommentID) as comments
    FROM Post p
    LEFT JOIN PostLike l ON p.PostID = l.PostID
    LEFT JOIN PostDislike d ON p.PostID = d.PostID
    LEFT JOIN Comment c ON p.PostID = c.PostID
    WHERE p.UserID = ?
    GROUP BY p.PostID
`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		return nil
	}
	defer rows.Close()

	var posts []CreatedPost
	for rows.Next() {
		var post CreatedPost
		err := rows.Scan(&post.PostID, &post.Title, &post.Content,
			&post.CreateDate, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts
}

func getUserComments(userID int) []CommentedPost {
	fmt.Printf("Fetching liked posts for user %d\n", userID) //debug
	query := `
        SELECT c.CommentID, c.PostID, c.content, c.created_at
        FROM Comment c
        JOIN Post p ON c.PostID = p.PostID
        WHERE c.UserID = ?
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []CommentedPost
	for rows.Next() {
		var comment CommentedPost
		rows.Scan(&comment.CommentID, &comment.PostID, &comment.Comment, &comment.CreateDate)
		comments = append(comments, comment)
	}
	return comments
}

func getLikedPosts(userID int) []LikedPost {
	fmt.Printf("Fetching disliked posts for user %d\n", userID) //debug
	query := `
        SELECT p.PostID, p.title, p.content, p.PostDate, u.username
        FROM Post p
        JOIN PostLike l ON p.PostID = l.PostID
        JOIN User u ON p.UserID = u.UserID
        WHERE l.UserID = ?
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []LikedPost
	for rows.Next() {
		var post LikedPost
		rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CreateDate, &post.Author)
		posts = append(posts, post)
	}
	return posts
}

func getDislikedPosts(userID int) []DislikedPost {
	query := `
        SELECT p.PostID, p.title, p.content, p.PostDate, u.username
        FROM Post p
        JOIN PostDislike d ON p.PostID = d.PostID
        JOIN User u ON p.UserID = u.UserID
        WHERE d.UserID = ?
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []DislikedPost
	for rows.Next() {
		var post DislikedPost
		rows.Scan(&post.PostID, &post.Title, &post.Content, &post.CreateDate, &post.Author)
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
