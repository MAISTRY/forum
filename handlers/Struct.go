package handlers

type Err struct {
	ErrorMessage string `json:"errorMessage"`
	Statuscode   int    `json:"statuscode"`
}

type Comment struct {
	CmtID    int    `json:"CmtID"`
	Content  string `json:"CmtContent"`
	CmtDate  string `json:"CmtDate"`
	Username string `json:"CmtUsername"`
	Likes    int    `json:"CmtLikes"`
	Dislikes int    `json:"CmtDislikes"`
}

type Post struct {
	PostID     int      `json:"PostID"`
	UserID     int      `json:"UserID"`
	PostDate   string   `json:"PostDate"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	ImagePath  *string  `json:"imagePath"`
	Username   string   `json:"username"`
	Likes      int      `json:"Likes"`
	Dislikes   int      `json:"Dislikes"`
	CmtCount   int      `json:"CmtCount"`
	Categories []string `json:"Categories"`
}

type categories struct {
	CategoryName string `json:"CategoryName"`
	Posts        []Post `json:"Posts"`
}

var Action action

type action struct {
	Message      string `json:"Message"`
	LikeCount    int    `json:"LikeCount"`
	DislikeCount int    `json:"DislikeCount"`
}

type CommentedPost struct {
	UserID     int    `json:"UserID"`
	UserName   string `json:"UserName"`
	CommentID  int    `json:"CommentID"`
	PostID     int    `json:"PostID"`
	Comment    string `json:"Comment"`
	CreateDate string `json:"CreateDate"`
	Likes      int    `json:"Likes"`
	Dislikes   int    `json:"Dislikes"`
}

type CommentRequest struct {
	PostID  string `json:"postId"`
	Comment string `json:"comment"`
}
