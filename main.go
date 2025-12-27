package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Film struct {
	Title string
	Genre string
}

var films []Film
var tmpl *template.Template

const (
	totalData = 500000
	maxInput  = 50000
	repeat    = 1000
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Generate 500000 data film...")
	generateFilms(totalData)

	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateFilms(n int) {
	genres := []string{"Romance", "Horror", "Action", "Comedy"}
	titles := []string{
		"Eternal Love", "Broken Heart", "Midnight Story", "Last Promise",
		"Dark Memories", "Haunted Night", "Final Battle", "Hidden Truth",
		"Forever Yours", "Silent Scream", "Endless Journey", "Lost Soul",
	}

	for i := 0; i < n; i++ {
		films = append(films, Film{
			Title: fmt.Sprintf("%s %d", titles[rand.Intn(len(titles))], i),
			Genre: genres[rand.Intn(len(genres))],
		})
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	genre := r.FormValue("genre")
	method := r.FormValue("method")
	n, _ := strconv.Atoi(r.FormValue("n"))

	if n > maxInput {
		n = maxInput
	}
	if n <= 0 {
		n = 1
	}

	var result []Film
	start := time.Now()

	for i := 0; i < repeat; i++ {
		if method == "iterative" {
			result = searchIterative(films[:n], genre)
		} else {
			result = []Film{}
			searchRecursive(films[:n], genre, 0, &result)
		}
	}

	elapsed := float64(time.Since(start).Nanoseconds()) / 1e6
	elapsed = elapsed / repeat

	data := struct {
		Genre  string
		Method string
		N      int
		Time   float64
		Result []Film
	}{
		Genre:  genre,
		Method: method,
		N:      n,
		Time:   elapsed,
		Result: result,
	}

	tmpl.Execute(w, data)
}

func searchIterative(data []Film, genre string) []Film {
	var res []Film
	for _, film := range data {
		if film.Genre == genre {
			res = append(res, film)
		}
	}
	return res
}

func searchRecursive(data []Film, genre string, idx int, res *[]Film) {
	if idx >= len(data) {
		return
	}
	if data[idx].Genre == genre {
		*res = append(*res, data[idx])
	}
	searchRecursive(data, genre, idx+1, res)
}
