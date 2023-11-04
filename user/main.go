package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const database_url string = "../database.db"
const create_user_table string = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT,
	);`

func main() {
	db, err := sql.Open("sqlite3", database_url)

	if err != nil {
		log.Println("Error opening database: ", err)
		os.Exit(1)
	}

	if _, err := db.Exec(create_user_table); err != nil {
		log.Println("Error creating users table: ", err)
		os.Exit(1)
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world, it's user1 service\n")
	})

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting USER server on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting USER server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Printf("Got signal: %s, exiting.", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
