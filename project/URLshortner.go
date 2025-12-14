package project

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string    `json:"id"`
	Originalurl  string    `json:"original_url"`
	Shorturl     string    `json:"short_url"`
	Creationdate time.Time `json:"creation_date"`
}

var UrlDB = make(map[string]URL)

func GenerateshortUrl(Originalurl string) string {
	hasher := md5.New()
	hasher.Write([]byte(Originalurl))
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data:", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Encodetostring:", hash)
	fmt.Println("final string:", hash[:8])
	return hash[:8]
}
func Createurl(originalurl string) string {
	shorturl := GenerateshortUrl(originalurl)
	id := shorturl
	UrlDB[id] = URL{
		Id:           id,
		Originalurl:  originalurl,
		Shorturl:     shorturl,
		Creationdate: time.Now(),
	}
	return shorturl
}
func Geturl(id string) (URL, error) {
	url, ok := UrlDB[id]
	if !ok {
		return URL{}, errors.New("errors not found")
	}
	return url, nil
}
func RootpageUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the URL shortner! Use /shortner to create short urls. ")

}
func ShortURLhandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shorturl_ := Createurl(data.URL)
	response := struct {
		Shorturl string `json:"shorturl"`
	}{Shorturl: shorturl_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func redirecturlhandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := Geturl(id)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)

	}
	http.Redirect(w, r, url.Originalurl, http.StatusFound)
}
func URLshortner() {
	http.HandleFunc("/", RootpageUrl)
	http.HandleFunc("/shortner", ShortURLhandler)
	http.HandleFunc("/redirect/", redirecturlhandler)

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)

}
