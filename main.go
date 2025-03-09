package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/immanuel-254/blog/auth"
	"github.com/immanuel-254/blog/cmd"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	// connect to db
	db, err := sql.Open("sqlite3", os.Getenv("DB"))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	auth.DB = db

	// migrate to database
	goose.SetDialect("sqlite3")

	// Apply all "up" migrations
	err = goose.Up(auth.DB, "auth/migrations")
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")

	if len(os.Args) < 1 {
		panic("There has to be exactly one argument")
	} else {
		if os.Args[1] == "createadmin" {
			cmd.CreateAdminUser()
		} else if os.Args[1] == "runserver" {
			cmd.Api()
		} else {
			panic("Invalid Argument")
		}
	}

}
