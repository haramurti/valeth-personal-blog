package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// 1. Model: Ini bentuk data lo
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    string    `json:"date"`
}

var dataFile = "data/data.json"

// 2. Helper: Buat baca file JSON (Si Pustakawan)
func getPosts() ([]Post, error) {
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		// Kalau file belum ada, balikin array kosong
		if os.IsNotExist(err) {
			return []Post{}, nil
		}
		return nil, err
	}
	var posts []Post
	json.Unmarshal(file, &posts)
	return posts, nil
}

// 3. Handler: Halaman Depan (Buat Tamu)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, _ := getPosts()
	
	// Nanti ini diganti pake render HTML dari folder 'views'
	// Sekarang kita tampilin teks mentah dulu buat tes
	fmt.Fprintf(w, "<h1>Blog Aletha</h1>")
	for _, p := range posts {
		fmt.Fprintf(w, "<h3>%s</h3><p>%s</p><hr>", p.Title, p.Content)
	}
}

// 4. Handler: Admin Login (Hardcode Sederhana)
func adminHandler(w http.ResponseWriter, r *http.Request) {
	// Basic Auth Logic (Satpam Galak)
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != "rahasia123" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Eits, mau ngapain? Cuma yang punya hati yang boleh masuk.", 401)
		return
	}

	fmt.Fprintf(w, "<h1>Welcome Admin!</h1><p>Di sini tempat nulis curhatan.</p>")
}

func main() {
	// Routing
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/admin", adminHandler)

	// Serve Static Files (CSS/Gambar) dari folder public
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	fmt.Println("Server jalan di http://localhost:8080")
	fmt.Println("Gas ngoding bang!")
	http.ListenAndServe(":8080", nil)
}