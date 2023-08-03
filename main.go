package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import PostgreSQL driver

	"github.com/jmoiron/sqlx"
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

type Post struct {
	PostID   int    `db:"post_id"`
	Username string `db:"username"`
	Title    string `db:"title"`
	Link     string `db:"link"`
	Upvotes  int    `db:"upvotes"`
	Content  string `db:"content"`
}

// Get all posts GET/api/posts/
func (app *App) getAllPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

	}

}

// convert to a method of App struct
// POST /v1/posts
// GET /v1/posts/:id
// POST /v1/posts/:id
// POST /v1/posts/:id/resume
// DELETE /v1/posts/:id
// GET /v1/posts
// GET /v1/posts/search

/*
http method GET endpoint "/healthcheck" -> 200
*/
func (app *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Val input data before inserting it into the database to
	// prevent errors and security vulnerabilities

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Body val
	if post.Username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}
	if post.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}
	if post.Link == "" {
		http.Error(w, "Link cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = app.DB.NamedExec(`INSERT INTO posts (username, title, link, upvotes, content) VALUES (:username, :title, :link, :upvotes, :content)`, &post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) initDB() {
	var err error
	app.DB, err = sqlx.Connect("postgres", "postgres://postgres:postgrespw@localhost:55000/postgres?sslmode=disable")
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
	app := &App{}
	app.initDB()
	defer app.closeDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})

	http.HandleFunc("/healthcheck", app.healthCheck)
	http.HandleFunc("/v1/posts", app.createPost)
	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
