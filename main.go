package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// convert to a method of App struct
// New createUser method on App struct
func (app *App) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Get IP address of the client that connected to your server.
	// The clinet may be src client, proxy, or a load balancer.
	ip := r.RemoteAddr

	// If the request came through a proxy, RemoteAddr might contain the IP:Port, so let's split it:
	ip = strings.Split(ip, ":")[0]

	// Another way to get the IP if it came through a proxy is to check the "X-Forwarded-For" header:
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded // This could be a comma-separated list of IPs, the client's IP is the first one
	}

	fmt.Printf("Client IP Address: %s\n", ip)

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = app.DB.NamedExec(`INSERT INTO users (name) VALUES (:name)`, &u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

/*
The InitDB method is defined on the App struct, so it can be called on instances of App.
This encapsulates the database connection and initialization within the App struct.

You can also move the createUser handler function into this struct and turn
it into a method, and you can add any other application-specific data to the App struct.
This can help keep your code organized and allow you to easily share data within your
application without using global variables
*/
func (app *App) initDB(dataSourceName string) {
	var err error
	app.DB, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer DB.Close()
}

func main() {
	app := &App{}
	app.initDB("user=username password=password dbname=dbname sslmode=disable")
	http.HandleFunc("/user", app.createUser)
	http.ListenAndServe(":8080", nil)
}
