package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// Global variable at the package level, it is accessible from all files within the same package

// internal DB. This will be coverted into a DB later
var DB []Post

/*
Encapsulation refers to the practice of bundling related data and
functions into a single unit, often for the purposes of information hiding and abstraction

Encapsulation often means defining methods
on types and keeping related data and functions within the same package

App is a struct that encapsulates a *sqlx.DB pointer.
We encapsulate the database connection in a struct App
*/
type App struct {
	DB []Post
}

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Username string `json:"username"`
}

func printDB(db []Post) {
	fmt.Println("DB Contents:")
	for _, post := range db {
		fmt.Printf("ID: %d, Title: %s, Link: %s, Username: %s\n", post.ID, post.Title, post.Link, post.Username)
	}
}

func (app *App) handlePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// store the decoded JSON data representing a Post object
		var newPost Post

		// The json package to decode the JSON data from the request body (r.Body) into the newPost variable.
		// The Decode method decodes the JSON data and returns an error if there was any issue during decoding.
		err := json.NewDecoder(r.Body).Decode(&newPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		app.DB = append(app.DB, newPost)
		printDB(app.DB)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}
}

func main() {
	app := &App{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Resource for this API is post. Post contain comments
	http.HandleFunc("/v1/posts/", app.handlePosts)

	// IP address and port
	// Docker NOTE: To make the server accessible from
	// outside the container, you need to bind it to all network interfaces (0.0.0.0) instead.
	addr := "0.0.0.0:8080"
	fmt.Printf("Server is running on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}

}
