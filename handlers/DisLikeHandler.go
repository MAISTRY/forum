package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/DB"
	"log"
	"net/http"
)

// PostDisLikeHandler handles the HTTP POST request for disliking a post.
// It manages the dislike action on a post, including checking if the post
// is already liked or disliked by the user, and updates the database accordingly.
// After processing, it returns the updated like and dislike counts for the post
// and creates a notification for the post owner.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request details and form data.
//
// The function doesn't return any value, but writes the response to the http.ResponseWriter.
// It sends a JSON response with the updated like and dislike counts, or an error if any occurs.
func PostDisLikeHandler(w http.ResponseWriter, r *http.Request) {
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
	postID := r.FormValue("postId")

	userID, err := getUserIDByCookie(r, db)
	if err != nil {
		log.Printf("Error getting user ID %v\n", err)
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

	insertPostDislikeNotification := func(db *sql.DB, userID, postID string) {
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
		_, err = stmt.Exec(userID, PostOwnerID, postID, "PostDislike")
		if err != nil {
			log.Printf("Error inserting to the DB %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if isDisliked {
		_, err = db.Exec(deleteDisLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error deleting the dislike %v\n", err)
			http.Error(w, "Error disliking post", http.StatusInternalServerError)
			return
		}
	} else if isLiked {
		_, err = db.Exec(deleteLikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the like %v\n", err)
			http.Error(w, "Error inserting like", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertDislikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the dislike %v\n", err)
			http.Error(w, "Error inserting dislike", http.StatusInternalServerError)
			return
		}
		insertPostDislikeNotification(db, userID, postID)
	} else {
		_, err = db.Exec(insertDislikeQuery, userID, postID)
		if err != nil {
			log.Printf("Error inserting the dislike %v\n", err)
			http.Error(w, "Error inserting dislike", http.StatusInternalServerError)
			return
		}
		insertPostDislikeNotification(db, userID, postID)
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

	Action.Message = "DisLiked post"
	Action.LikeCount = likeCount
	Action.DislikeCount = dislikeCount

	log.Println(Action)
	json.NewEncoder(w).Encode(Action)
}

// CommentDislikeHandler handles the HTTP POST request for disliking a comment.
// It manages the dislike action on a comment, including checking if the comment
// is already liked or disliked by the user, and updates the database accordingly.
// After processing, it returns the updated like and dislike counts for the comment
// and creates a notification for the comment owner.
//
// Parameters:
//   - w: http.ResponseWriter to write the HTTP response.
//   - r: *http.Request containing the HTTP request details and form data.
//
// The function doesn't return any value, but writes the response to the http.ResponseWriter.
// It sends a JSON response with the updated like and dislike counts, or an error if any occurs.
func CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
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

	insertCommentDislikeNotification := func(db *sql.DB, userID, commentID string) {
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
		_, err = stmt.Exec(userID, PostOwnerID, commentID, "CommentDislike")
		if err != nil {
			log.Printf("Error inserting to the DB %v\n", err)
			http.Error(w, "Error preping the notification: Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if isCommentDisliked {
		_, err = db.Exec(deleteCommentDislikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error deleting the comment dislike %v\n", err)
			http.Error(w, "Error disliking comment", http.StatusInternalServerError)
			return
		}
	} else if isCommentLiked {
		_, err = db.Exec(deleteCommentLikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment like %v\n", err)
			http.Error(w, "Error inserting comment like", http.StatusInternalServerError)
			return
		}
		_, err = db.Exec(insertCommentDislikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment dislike %v\n", err)
			http.Error(w, "Error inserting comment dislike", http.StatusInternalServerError)
			return
		}
		insertCommentDislikeNotification(db, userID, commentID)
	} else {
		_, err = db.Exec(insertCommentDislikeQuery, userID, commentID)
		if err != nil {
			log.Printf("Error inserting the comment dislike %v\n", err)
			http.Error(w, "Error inserting comment dislike", http.StatusInternalServerError)
			return
		}
		insertCommentDislikeNotification(db, userID, commentID)
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

	Action.Message = "DisLiked comment"
	Action.LikeCount = commentLikeCount
	Action.DislikeCount = commentDislikeCount

	log.Println(Action)
	json.NewEncoder(w).Encode(Action)
}
