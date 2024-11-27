package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/DB"
	"log"
	"net/http"
)

// PostLikeHandler handles HTTP requests for liking or disliking a post.
// It manages the state of likes and dislikes for a given post and user,
// updating the database accordingly. The function returns a JSON response
// with the updated like and dislike counts for the post.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response.
//   - r *http.Request: The HTTP request containing the post ID and user information.
//
// The function does not return any values directly, but writes the response to the http.ResponseWriter.
// In case of success, it returns a JSON object with the following structure:
//
//	{
//	  "message": "Liked post",
//	  "LikeCount": <number of likes>,
//	  "DislikeCount": <number of dislikes>
//	}
func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("error opening database:  %v\n", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	r.ParseForm()
	postID := r.FormValue("postId")

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		// http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	isLiked, err := IsLiked(db, userID, postID)
	if err != nil {
		log.Printf("Error checking if user has liked post %v\n", err)
		http.Error(w, "Error checking if user has liked post", http.StatusInternalServerError)
		return
	}
	isDisliked, err := IsDisliked(db, userID, postID)
	if err != nil {
		log.Printf("Error checking if user has disliked post %v\n", err)
		http.Error(w, "Error checking if user has disliked post", http.StatusInternalServerError)
		return
	}

	insertPostLikeNotification := func(db *sql.DB, userID, postID string) {
		stmt, err := db.Prepare(`
        INSERT INTO Notification (UserID, UserToNotify, PostID, NotificationType)
        VALUES (?,?,?,?);
    	`)
		if err != nil {
			log.Printf("Error preparing the statement: %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		PostOwnerID, err := GetPostOwnerID(postID, db)
		if err != nil {
			log.Printf("Error getting post owner ID %v\n", err)
			http.Error(w, "Error getting post owner ID", http.StatusInternalServerError)
			return
		}
		_, err = stmt.Exec(userID, PostOwnerID, postID, "PostLike")
		if err != nil {
			log.Printf("Error inserting to the DB %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if isLiked {
		_, err = db.Exec(deleteLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the like %v\n", err)
			http.Error(w, "Error inserting like", http.StatusInternalServerError)
			return
		}
	} else if isDisliked {
		_, err = db.Exec(deleteDisLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error deleting the dislike %v\n", err)
			http.Error(w, "Error deleting dislike", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the like %v\n", err)
			http.Error(w, "Error inserting like", http.StatusInternalServerError)
			return
		}
		insertPostLikeNotification(db, userID, postID)
	} else {
		_, err = db.Exec(insertLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the like %v\n", err)
			http.Error(w, "Error inserting like", http.StatusInternalServerError)
			return
		}
		insertPostLikeNotification(db, userID, postID)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	likeCount, err := DB.GetLikeCount(db, postID)
	if err != nil {
		log.Printf("Error getting like count %v\n", err)
		http.Error(w, "Error getting like count", http.StatusInternalServerError)
		return
	}
	dislikeCount, err := DB.GetDislikeCount(db, postID)
	if err != nil {
		log.Printf("Error getting dislike count %v\n", err)
		http.Error(w, "Error getting dislike count", http.StatusInternalServerError)
		return
	}

	Action.Message = "Liked post"
	Action.LikeCount = likeCount
	Action.DislikeCount = dislikeCount

	log.Println(Action)
	json.NewEncoder(w).Encode(Action)
}

// CommentLikeHandler handles HTTP requests for liking or disliking a comment.
// It manages the state of likes and dislikes for a given comment and user,
// updating the database accordingly. The function returns a JSON response
// with the updated like and dislike counts for the comment.
//
// Parameters:
//   - w http.ResponseWriter: The response writer to send the HTTP response.
//   - r *http.Request: The HTTP request containing the comment ID and user information.
//
// The function does not return any values directly, but writes the response to the http.ResponseWriter.
// In case of success, it returns a JSON object with the following structure:
//
//	{
//	  "message": "liked comment",
//	  "LikeCount": <number of likes>,
//	  "DislikeCount": <number of dislikes>
//	}
func CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error opening the DB %v\n", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	r.ParseForm()
	commentID := r.FormValue("commentId")

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user ID %v\n", err)
		// http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	isCommentLiked, err := IsCommentLiked(db, userID, commentID)
	if err != nil {
		log.Printf("Error checking if user has liked comment %v\n", err)
		http.Error(w, "Error checking if user has liked comment", http.StatusInternalServerError)
		return
	}
	isCommentDisliked, err := IsCommentDisliked(db, userID, commentID)
	if err != nil {
		log.Printf("Error checking if user has disliked comment %v\n", err)
		http.Error(w, "Error checking if user has disliked comment", http.StatusInternalServerError)
		return
	}

	insertPostLikeNotification := func(db *sql.DB, userID, commentID string) {
		stmt, err := db.Prepare(`
        INSERT INTO Notification (UserID, UserToNotify, CommentID, NotificationType)
        VALUES (?,?,?,?);
    `)
		if err != nil {
			log.Printf("Error preparing the statement: %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		PostOwnerID, err := GetCommentOwnerID(commentID, db)
		if err != nil {
			log.Printf("Error getting post owner ID %v\n", err)
			http.Error(w, "Error getting post owner ID", http.StatusInternalServerError)
			return
		}
		_, err = stmt.Exec(userID, PostOwnerID, commentID, "CommentLike")
		if err != nil {
			log.Printf("Error inserting to the DB %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if isCommentLiked {
		_, err = db.Exec(deleteCommentLikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error deleting the comment dislike %v\n", err)
			http.Error(w, "Error disliking comment", http.StatusInternalServerError)
			return
		}
	} else if isCommentDisliked {
		_, err = db.Exec(deleteCommentDislikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment like %v\n", err)
			http.Error(w, "Error inserting comment like", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertCommentLikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment dislike %v\n", err)
			http.Error(w, "Error inserting comment dislike", http.StatusInternalServerError)
			return
		}
		insertPostLikeNotification(db, userID, commentID)
	} else {
		_, err = db.Exec(insertCommentLikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment dislike %v\n", err)
			http.Error(w, "Error inserting comment dislike", http.StatusInternalServerError)
			return
		}
		insertPostLikeNotification(db, userID, commentID)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	commentLikeCount, err := DB.GetCommentLikeCount(db, commentID)
	if err != nil {
		log.Printf("Error getting comment like count %v\n", err)
		http.Error(w, "Error getting comment like count", http.StatusInternalServerError)
		return
	}
	commentDislikeCount, err := DB.GetCommentDislikeCount(db, commentID)
	if err != nil {
		log.Printf("Error getting comment dislike count %v\n", err)
		http.Error(w, "Error getting comment dislike count", http.StatusInternalServerError)
		return
	}

	Action.Message = "Liked comment"
	Action.LikeCount = commentLikeCount
	Action.DislikeCount = commentDislikeCount

	log.Println(Action)
	json.NewEncoder(w).Encode(Action)
}
