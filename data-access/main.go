package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// connStr := "user=postgres dbname=data-access sslmode=verify-full"
	connStr := "host=localhost port=5432 user=postgres password=dataaccess dbname=recordings sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to Connect the database %v", err)
	}

	log.Println("Successfully connected to the data base")

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	defer rows.Close()

	log.Println("Query executed successfully. Results:")

	for rows.Next() {
		var id uint
		var title string
		var artist string
		var price float64

		if err := rows.Scan(&id, &title, &artist, &price); err != nil {
			log.Fatalf("Failed to scan the row: %v", err)
		}

		log.Printf("Album: %s, Artist:%s, Price:%f", title, artist, price)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("%v", err)
	}
}
