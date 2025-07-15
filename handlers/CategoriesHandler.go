package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
)

// CategoriesHandler handles HTTP GET requests for retrieving categorized posts and their associated comments.
// It fetches posts from the database, organizes them by category, and includes comment information for each post.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response.
//   - r *http.Request: The HTTP request received from the client.
//
// The function doesn't return any value directly, but writes a JSON response to the http.ResponseWriter.
// The JSON response contains an array of category objects, each with a list of posts and their comments.
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Internal Server Error", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `
        SELECT 
            p.PostID, 
            p.UserID,
            p.PostDate,
            p.title,
            p.content,
            p.ImagePath,
            u.username,
            c.title AS category,
            COALESCE(pl.likes, 0) AS likes,
            COALESCE(pdl.dislike, 0) AS dislikes,
            COALESCE(cmt.comments, 0) AS comments
        FROM 
            Post p
        JOIN 
            User u ON p.UserID = u.UserID
        JOIN 
            PostCategory pc ON p.PostID = pc.PostID
        JOIN 
            Category c ON pc.CategoryID = c.CategoryID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS likes FROM PostLike GROUP BY PostID
        ) AS pl ON p.PostID = pl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS dislike FROM PostDislike GROUP BY PostID
        ) AS pdl ON p.PostID = pdl.PostID
        LEFT JOIN (
            SELECT PostID, COUNT(*) AS comments FROM Comment GROUP BY PostID
        ) AS cmt ON p.PostID = cmt.PostID
        ORDER BY 
            p.PostDate DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error querying posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	categoriesMap := make(map[string]map[int]Post)
	postCategoriesMap := make(map[int][]string)

	for rows.Next() {
		var post Post
		var category string

		if err := rows.Scan(
			&post.PostID, &post.UserID, &post.PostDate, &post.Title, &post.Content,
			&post.ImagePath, &post.Username, &category, &post.Likes, &post.Dislikes,
			&post.CmtCount,
		); err != nil {
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}

		if _, exists := postCategoriesMap[post.PostID]; !exists {
			postCategoriesMap[post.PostID] = []string{}
		}
		categoryExists := false
		for _, existingCategory := range postCategoriesMap[post.PostID] {
			if existingCategory == category {
				categoryExists = true
				break
			}
		}
		if !categoryExists {
			postCategoriesMap[post.PostID] = append(postCategoriesMap[post.PostID], category)
		}

		if categoriesMap[category] == nil {
			categoriesMap[category] = make(map[int]Post)
		}
		if _, exists := categoriesMap[category][post.PostID]; !exists {
			categoriesMap[category][post.PostID] = post
		}
	}

	var categoryGroup []categories
	for categoryName, posts := range categoriesMap {
		var category categories
		category.CategoryName = categoryName
		var categoryPosts []Post

		for postID, post := range posts {
			post.Categories = postCategoriesMap[postID]
			categoryPosts = append(categoryPosts, post)
		}

		sort.Slice(categoryPosts, func(i, j int) bool {
			return categoryPosts[i].PostDate > categoryPosts[j].PostDate
		})

		category.Posts = categoryPosts
		categoryGroup = append(categoryGroup, category)
	}

	sort.Slice(categoryGroup, func(i, j int) bool {
		return categoryGroup[i].CategoryName < categoryGroup[j].CategoryName
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryGroup)
}
