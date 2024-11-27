package DB

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	/*
		// TODO: Add more tables and relationships as needed for the main and optional projects.
		TODO: Ensure to use proper data types, constraints, and indexes where appropriate.
		// TODO: Use FKs in all the tables and enforce the Refrential integrity(PRAGMA).
		// TODO: insert cookies table in the ERD and it relations with the rest of the tables.
		! for @musabt: review the code for any improvements.
	*/

	enforcementOfFKs = `PRAGMA FOREIGN_KEYS = 1;`

	CreateUserTableQuery = `CREATE TABLE IF NOT EXISTS User(
		UserID INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		firstname TEXT NOT NULL,
		lastname TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE CHECK(email LIKE '%@%.%'),
        password TEXT NOT NULL,
		gender TEXT NOT NULL CHECK(gender IN ('M', 'F')),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        privilege INTEGER NOT NULL CHECK(privilege >= 1 AND privilege <= 3)
	);`
	CreatePostTableQuery = `CREATE TABLE IF NOT EXISTS Post(
        PostID INTEGER PRIMARY KEY AUTOINCREMENT,
        UserID INTEGER NOT NULL,
		PostDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		title TEXT NOT NULL,
        content TEXT NOT NULL,
	    ImagePath TEXT,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreateCategoryTableQuery = `CREATE TABLE IF NOT EXISTS Category(
		CategoryID INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
        description TEXT NOT NULL
	);`
	CreatePostCategoryTableQuery = `CREATE TABLE IF NOT EXISTS PostCategory(
		PostID INTEGER NOT NULL,
        CategoryID INTEGER NOT NULL,
        PRIMARY KEY (PostID, CategoryID),
        FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
        FOREIGN KEY (CategoryID) REFERENCES Category(CategoryID) ON DELETE CASCADE
	);`
	CreateCommentTableQuery = `CREATE TABLE IF NOT EXISTS Comment(
		CommentID INTEGER PRIMARY KEY AUTOINCREMENT,
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
        content TEXT NOT NULL,
		CmtDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
	);`
	CreatePostLikeTableQuery = `CREATE TABLE IF NOT EXISTS PostLike(
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
		PRIMARY KEY (PostID, UserID),
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
    );`
	CreatePostDislikeTableQuery = `CREATE TABLE IF NOT EXISTS PostDislike(
        PostID INTEGER NOT NULL,
        UserID INTEGER NOT NULL,
		PRIMARY KEY (PostID, UserID),
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
    );`
	///// TODO: review the tables and check how fesiable the code will be...
	CreateCommentLikeTableQuery = `CREATE TABLE IF NOT EXISTS CommentLike(
		CommentID INTEGER NOT NULL,
		UserID INTEGER NOT NULL,
		PostID INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (CommentID, UserID),
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE
	);`
	CreateCommentDislikeTableQuery = `CREATE TABLE IF NOT EXISTS CommentDislike(
		CommentID INTEGER NOT NULL,
		UserID INTEGER NOT NULL,
		PostID INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (CommentID, UserID),
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE CASCADE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE
	);`
	// ! we change this way of storing the images
	// CreatePostImageTableQuery = `CREATE TABLE IF NOT EXISTS PostImage(
	// 	ImageID INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	PostID INTEGER NOT NULL,
	// 	image_data BLOB NOT NULL,
	// 	image_filename TEXT NOT NULL,
	// 	image_mimetype TEXT NOT NULL,
	// 	image_width INTEGER NOT NULL CHECK(image_width > 0 AND image_width <= 8192),
	// 	image_height INTEGER NOT NULL CHECK(image_height > 0 AND image_height <= 8192),
	// 	image_size INTEGER NOT NULL CHECK(image_size > 0 AND image_size <= 20971520),
	// 	FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE CASCADE
	// );`
	// ! need to be checked. @mmahmooda
	// * added UserToNotify (to know who's the user to get the notification)
	CreateNotificationTableQuery = `CREATE TABLE IF NOT EXISTS Notification (
		NotificationID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserID INTEGER NOT NULL,  -- User receiving the notification
		UserToNotify INTEGER NOT NULL,  -- User who is getting the notification (null if system notification)
		PostID INTEGER,           -- Post related to the notification (nullable if comment only)
		CommentID INTEGER,        -- Comment related to the notification (nullable if only a like)
		NotificationType TEXT NOT NULL CHECK(NotificationType IN ('PostLike', 'PostDislike', 'Comment', 'CommentLike', 'CommentDislike')),
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		IsRead BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (UserToNotify) REFERENCES User(UserID) ON DELETE CASCADE,
		FOREIGN KEY (PostID) REFERENCES Post(PostID) ON DELETE SET NULL,
		FOREIGN KEY (CommentID) REFERENCES Comment(CommentID) ON DELETE SET NULL
	);`
	sessionTableQuery = `CREATE TABLE IF NOT EXISTS Session (
		session_id TEXT PRIMARY KEY,
		user_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expiry_date TIMESTAMP,
		ip_address TEXT,
		FOREIGN KEY (user_id) REFERENCES User(UserID) ON DELETE CASCADE
	);`

	// // ------------------------------------------------------------ // //

	// ! WE WERE IDIOTS THINKING THAT WE GOTTA DO IT THIS WAY 👇💀💀
	// TODO: to store the cookies in a map instead of storing them in the database.
)

// InitDB initializes the database connection, creates the necessary tables, and performs any initial table filling.
// It opens a new SQLite database connection, checks the connection, and then calls the CreateTables function to create the required tables.
// If any errors occur during the database initialization or table creation, it logs a fatal error.
func InitDB() {
	db, err := sql.Open("sqlite3", "./meow.db")
	if err != nil {
		log.Fatalf("error creating database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	// * DONE
	CreateTables(db)
	// InitailTableFiller(db)
}

// CreateTables creates the necessary tables in the database.
// It enables foreign key constraints, creates the user, post, category, post_category, comment, like, dislike, and notification tables.
// If any errors occur during the table creation, it logs a fatal error.
// Finally, it logs a success message indicating that the tables were created successfully.
func CreateTables(db *sql.DB) {
	if _, err := db.Exec(enforcementOfFKs); err != nil {
		log.Fatalf("error enabling foreign keys: %v", err)
	}
	if _, err := db.Exec(CreateUserTableQuery); err != nil {
		log.Fatalf("error creating the user table: %v", err)
	}
	if _, err := db.Exec(CreatePostTableQuery); err != nil {
		log.Fatalf("error creating the post table: %v", err)
	}
	if _, err := db.Exec(CreateCategoryTableQuery); err != nil {
		log.Fatalf("error creating the category table: %v", err)
	}
	if _, err := db.Exec(CreatePostCategoryTableQuery); err != nil {
		log.Fatalf("error creating the post_category table: %v", err)
	}
	if _, err := db.Exec(CreateCommentTableQuery); err != nil {
		log.Fatalf("error creating the comment table: %v", err)
	}
	if _, err := db.Exec(CreatePostLikeTableQuery); err != nil {
		log.Fatalf("error creating the like table: %v", err)
	}
	if _, err := db.Exec(CreatePostDislikeTableQuery); err != nil {
		log.Fatalf("error creating the dislike table: %v", err)
	}
	if _, err := db.Exec(CreateCommentLikeTableQuery); err != nil {
		log.Fatalf("error creating the comment_like table: %v", err)
	}
	if _, err := db.Exec(CreateCommentDislikeTableQuery); err != nil {
		log.Fatalf("error creating the comment_dislike table: %v", err)
	}
	// if _, err := db.Exec(CreatePostImageTableQuery); err != nil {
	// 	log.Fatalf("error creating the post_image table: %v", err)
	// }
	if _, err := db.Exec(CreateNotificationTableQuery); err != nil {
		log.Fatalf("error creating the Notification table: %v", err)
	}
	if _, err := db.Exec(sessionTableQuery); err != nil {
		log.Fatalf("error creating the session table: %v", err)
	}

	log.Println("Tables created successfully...")
}
