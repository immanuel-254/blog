package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/immanuel-254/blog/cmd"
	"github.com/immanuel-254/blog/database"
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

	database.DB = db

	// migrate to database
	goose.SetDialect("sqlite3")

	// Apply all "up" migrations
	err = goose.Up(database.DB, "auth/migrations")
	if err != nil {
		log.Fatalf("Failed to auth apply migrations: %v", err)
	}

	// Apply all "up" migrations
	err = goose.Up(database.DB, "blog/migrations", goose.WithAllowMissing())
	if err != nil {
		log.Fatalf("Failed to blog apply migrations: %v", err)
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
