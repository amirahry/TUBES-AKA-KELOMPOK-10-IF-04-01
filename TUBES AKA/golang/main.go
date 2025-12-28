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

var (
	films []Film
	tmpl  *template.Template
)

const (
	totalData = 500000
	maxInput  = 50000
	repeat    = 1000
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Generate data film...")
	generateFilms(totalData)

	// HUBUNGKAN KE HTML (TIDAK GANTI NAMA FILE)
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Gagal load template:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/search", searchHandler)

	fmt.Println("Server running di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ================= DATA =================
func generateFilms(n int) {
	genres := []string{"Romance", "Horror", "Action", "Comedy"}
	titles := []string{
		"Eternal Love", "Broken Heart", "Midnight Story", "Last Promise",
		"Dark Memories", "Haunted Night", "Final Battle", "Hidden Truth",
	}

	films = make([]Film, 0, n)

	for i := 0; i < n; i++ {
		films = append(films, Film{
			Title: fmt.Sprintf("%s %d", titles[rand.Intn(len(titles))], i),
			Genre: genres[rand.Intn(len(genres))],
		})
	}
}

// ================= HANDLER =================
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	genre := r.FormValue("genre")
	method := r.FormValue("method")

	n, err := strconv.Atoi(r.FormValue("n"))
	if err != nil || n <= 0 {
		n = 1
	}
	if n > maxInput {
		n = maxInput
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

	elapsed := float64(time.Since(start).Nanoseconds()) / 1e6 / repeat

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

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ================= ALGORITMA =================
func searchIterative(data []Film, genre string) []Film {
	res := make([]Film, 0)
	for _, f := range data {
		if f.Genre == genre {
			res = append(res, f)
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
