package project2

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// URL model
type URL struct {
	Id           string    `json:"id"`
	Originalurl  string    `json:"original_url"`
	Shorturl     string    `json:"short_url"`
	Creationdate time.Time `json:"creation_date"`
}

// Handler struct that holds DB connection
type URLShortenerHandler struct {
	db *sql.DB
}

// Constructor: initialize DB and return handler
func NewURLShortenerHandler(dsn string) (*URLShortenerHandler, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Connected to MySQL!")
	return &URLShortenerHandler{db: db}, nil
}

// Generate short URL
func generateShortURL(originalurl string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalurl))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

// Create URL entry in DB
func (h *URLShortenerHandler) createURL(originalurl string) (string, error) {
	shorturl := generateShortURL(originalurl)
	id := shorturl

	_, err := h.db.Exec("INSERT INTO urls (id, original_url, short_url, creation_date) VALUES (?, ?, ?, ?)",
		id, originalurl, shorturl, time.Now())
	if err != nil {
		return "", err
	}
	return shorturl, nil
}

// Fetch URL from DB
func (h *URLShortenerHandler) getURL(id string) (URL, error) {
	var url URL
	row := h.db.QueryRow("SELECT id, original_url, short_url, creation_date FROM urls WHERE id = ?", id)
	err := row.Scan(&url.Id, &url.Originalurl, &url.Shorturl, &url.Creationdate)
	if err != nil {
		if err == sql.ErrNoRows {
			return URL{}, errors.New("URL not found")
		}
		return URL{}, err
	}
	return url, nil
}

// Handlers
func rootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the URL shortener! Use /shortner to create short URLs.")
}

func (h *URLShortenerHandler) shortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shorturl, err := h.createURL(data.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Shorturl string `json:"shorturl"`
	}{Shorturl: shorturl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *URLShortenerHandler) redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := h.getURL(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.Originalurl, http.StatusFound)
}

func URLshortner1() {
	dsn := "root:root@tcp(127.0.0.1:3306)/urlshortner"
	handler, err := NewURLShortenerHandler(dsn)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", rootPageURL)
	http.HandleFunc("/shortner", handler.shortURLHandler)
	http.HandleFunc("/redirect/", handler.redirectURLHandler)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
