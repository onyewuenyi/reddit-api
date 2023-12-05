package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// Global variable at the package level, it is accessible from all files within the same package
var DB *sqlx.DB

/*
Encapsulation refers to the practice of bundling related data and
functions into a single unit, often for the purposes of information hiding and abstraction

Encapsulation often means defining methods
on types and keeping related data and functions within the same package

App is a struct that encapsulates a *sqlx.DB pointer.
We encapsulate the database connection in a struct App
*/
type App struct {
	DB *sqlx.DB
}

type Users struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	Date_created string `db:"date_created"` // Corrected database column name
	Karma        string `db:"karma"`
}

type Post struct {
	Id       int       `db:"id"`
	Link     string    `db:"link"`
	Title    string    `db:"title"`
	User_id  int       `db:"user_id"`
	PostDate time.Time `db:"post_date" sql:"timestamptz"`
}

type Comment struct {
	Id           int    `db:"id"`
	Post_id      int    `db:"post_id"`
	Text         string `db:"text"`
	User_id      int    `db:"user_id"`
	Comment_date string `db:"comment_date"`
}

type Replies struct {
	Id         int    `db:"id"`
	Comment_id int    `db:"comment_id"`
	Text       string `db:"text"`
	User_id    int    `db:"user_id"`
	Reply_date string `db:"reply_date"`
}

/*
http method GET endpoint /healthcheck -> 200
*/
func (app *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) handlePost(w http.ResponseWriter, r *http.Request) {
	/*
		Retrive posts from db and store in a slice posts
		Set res content type
		Encode posts slice as JSON
		Wr to res writer
	*/
	if r.Method == http.MethodGet {
		log.Println("Handling GET posts")
		queryParameters := r.URL.Query()
		id := queryParameters.Get("id")
		if id != "" {
			// make network req to retreive post with id
			var post Post
			query := "SELECT * FROM post WHERE id = $1"
			err := app.DB.Get(&post, query, id)
			if err != nil {

			}

			// set header req type that I will send

		} else {
			// Retrive all post in store it in posts, which is a list of Post structs
			var posts []Post
			query := "SELECT * FROM posts"
			err := app.DB.Select(&posts, query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error retrieving posts", err)
				return
			}
			// in http spec it is recommended to set header before writing
			// Set res headers
			w.Header().Set("Content-Type", "application/json")
			// Set the res status to 200 since I am retrieving data
			w.WriteHeader(http.StatusOK)

			// Convert posts to json and wr res
			err = json.NewEncoder(w).Encode(posts)

			// Return http status 500 if an error occured from encoding
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else if r.Method == http.MethodPost {
		var post Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if post.Link == "" {
			http.Error(w, "Link cannot be empty", http.StatusBadRequest)
			return
		}
		if post.Title == "" {
			http.Error(w, "Title cannot be empty", http.StatusBadRequest)
			return
		}
		if post.User_id == 0 {
			http.Error(w, "User_id cannot be empty", http.StatusBadRequest)
			return
		}

		_, err = app.DB.NamedExec(`INSERT INTO posts (title, user_id, link) VALUES (:title, :user_id, :link)`, &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Ensure that resources like the database connection are cleaned up if an err occurs
	defer r.Body.Close()
}

func (app *App) commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

	} else if r.Method == http.MethodPost {
		// Post comment for a specific post
		// Edit a comment for a specific post

	}
}

func (app *App) initDB() {
	var err error
	app.DB, err = sqlx.Connect("postgres", "postgres://postgres:postgrespw@localhost:55001/postgres?sslmode=disable")
	if err != nil {
		// TODO: update log.Fatalf
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

func (app *App) closeDB() {
	err := app.DB.Close()
	if err != nil {
		log.Fatalf("Unable to close database: %v\n", err)
	}
}

func main() {
	// init
	log.SetOutput(os.Stdout)
	log.SetPrefix("Reddit-API")
	app := &App{}
	app.initDB()
	defer app.closeDB()

	http.HandleFunc("/healthcheck", app.healthCheck)
	http.HandleFunc("/v1/posts", app.handlePost)
	http.HandleFunc("/v1/posts/{id}", app.handlePost)
	// http.HandleFunc("/v1/comments", app.CreateComment)
	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
