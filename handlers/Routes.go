package handlers

import (
	"forum/auth"
	mdlware "forum/middleware"
	"net/http"
)

func Routes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("/", HomePage)

	router.HandleFunc("/auth/status", CheckAuthHandler)

	router.HandleFunc("/Data-userLogin", LoginHandler)
	router.HandleFunc("/Data-userLogout", LogoutHandler)
	router.HandleFunc("/Data-userRegister", RegisterHandler)

	// ! START Google and Github auth
	router.HandleFunc("/auth/google/login", auth.HandleGoogleLogin)
	router.HandleFunc("/auth/github/login", auth.HandleGitHubLogin)
	router.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", "google")
		r.URL.RawQuery = q.Encode()
		auth.HandleOAuthCallback(w, r)
	})

	router.HandleFunc("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		q.Add("provider", "github")
		r.URL.RawQuery = q.Encode()
		auth.HandleOAuthCallback(w, r)
	})
	// ! END Google and Github auth

	router.HandleFunc("/Data-Post", PostHandler)
	router.HandleFunc("/Data-PostLike", PostLikeHandler)
	router.HandleFunc("/Data-PostDisLike", PostDisLikeHandler)

	router.HandleFunc("/Data-Comment", CommentHandler)
	router.HandleFunc("/Data-CommentLike", CommentLikeHandler)
	router.HandleFunc("/Data-CommentDisLike", CommentDislikeHandler)

	router.HandleFunc("/Data-CreatPost", CreatePostHandler)
	router.HandleFunc("/Data-CreatComment", CreatCommentHandler)

	router.HandleFunc("/Data-Profile", ProfileHandler)
	router.HandleFunc("/Data-Activity", ActivityHandler)
	router.HandleFunc("/Data-Categories", CategoriesHandler)
	router.HandleFunc("/Data-PublicCategories", PublicCategoriesHandler)

	// Admin routes
	router.HandleFunc("/Data-AdminStats", AdminStatsHandler)
	router.HandleFunc("/Data-AdminUsers", AdminUsersHandler)
	router.HandleFunc("/Data-AdminPromoteUser", AdminPromoteUserHandler)
	router.HandleFunc("/Data-AdminDemoteUser", AdminDemoteUserHandler)
	router.HandleFunc("/Data-AdminModerationRequests", AdminModerationRequestsHandler)
	router.HandleFunc("/Data-AdminRespondRequest", AdminRespondRequestHandler)

	// Report routes
	router.HandleFunc("/Data-ReportPost", ReportPostHandler)
	router.HandleFunc("/Data-AdminReports", AdminReportsHandler)
	router.HandleFunc("/Data-AdminRespondReport", AdminRespondReportHandler)
	router.HandleFunc("/Data-UserReports", UserReportsHandler)
	router.HandleFunc("/Data-CreateModerationRequest", CreateModerationRequestHandler)
	router.HandleFunc("/Data-AdminCategories", AdminCategoriesHandler)
	router.HandleFunc("/Data-AdminAddCategory", AdminAddCategoryHandler)
	router.HandleFunc("/Data-AdminDeleteCategory", AdminDeleteCategoryHandler)

	// Edit routes
	router.HandleFunc("/Data-EditPost", EditPostHandler)
	router.HandleFunc("/Data-GetPostForEdit", GetPostForEditHandler)
	router.HandleFunc("/Data-EditComment", EditCommentHandler)
	router.HandleFunc("/Data-GetCommentForEdit", GetCommentForEditHandler)

	// Delete routes (admin/moderator)
	router.HandleFunc("/Data-DeletePost", DelPostHandler)
	router.HandleFunc("/Data-DeleteComment", DeleteCommentHandler)

	// User delete routes (own content only)
	router.HandleFunc("/Data-UserDeletePost", UserDeletePostHandler)
	router.HandleFunc("/Data-UserDeleteComment", UserDeleteCommentHandler)

	// Notification routes
	router.HandleFunc("/Data-Notifications", NotificaionHandler)
	router.HandleFunc("/Data-NotificationCount", NotificationCountHandler)
	router.HandleFunc("/Data-MarkAsRead", MarkAsReadHandler)

	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./static/uploads"))))
	router.Handle("/scripts/", http.StripPrefix("/scripts/", http.FileServer(http.Dir("./static/scripts"))))
	router.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./static/styles"))))
	router.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	//// router.Handle("/Data-creatPost", middleware.AuthenticateUser(http.HandlerFunc(CreatePostHandler)))
	return mdlware.RateLimiter(router)
	// return router
}
