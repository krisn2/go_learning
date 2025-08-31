package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=dataaccess dbname=recordings sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Successfully connected to the database")

	mux := http.NewServeMux()

	// A single handler for /albums that routes based on the HTTP method
	mux.HandleFunc("/albums", albumsHandler)
	mux.HandleFunc("/albums/by-id", getAlbumsById)
	mux.HandleFunc("/albums/by-artist", getAlbumsByArtist)

	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// albumsHandler handles requests to the /albums endpoint and dispatches them based on the HTTP method.
func albumsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addAlbumHandler(w, r)
	case http.MethodGet:
		getAllalbums(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllalbums retrieves all albums from the database and returns them as a JSON array.
func getAllalbums(w http.ResponseWriter, r *http.Request) {
	var albs []Album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		http.Error(w, "Failed to get all albums", http.StatusInternalServerError)
		log.Println("Error querying all albums:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Error scanning album row:", err)
			return
		}
		albs = append(albs, alb)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error during row iteration", http.StatusInternalServerError)
		log.Println("Error during rows iteration:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albs); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// getAlbumsByArtist retrieves albums by a specific artist from the database.
func getAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artist") // Use lowercase "artist" for the query parameter

	if artistName == "" {
		http.Error(w, "Missing 'artist' query parameter", http.StatusBadRequest)
		return
	}

	albums, err := albumsByArtist(artistName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// getAlbumsById retrieves a single album by its ID from the database.
func getAlbumsById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	alb, err := albumById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(alb); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

// addAlbumHandler handles the creation of a new album.
func addAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var alb Album
	if err := json.NewDecoder(r.Body).Decode(&alb); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	id, err := addAlbum(alb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "Album added with ID: %d"}`, id)
}

// albumsByArtist is a helper function to query albums by artist from the database.
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumById is a helper function to query a single album by its ID from the database.
func albumById(id int64) (Album, error) {
	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum is a helper function to insert a new album into the database.
func addAlbum(alb Album) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO album(title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
