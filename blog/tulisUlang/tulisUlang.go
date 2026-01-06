package tulisulang

//fully sendiri
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

func isloggedin(r *http.Request) bool {
	// Coba minta cookie 'session_token' dari browser
	cookie, err := r.Cookie("session_token")

	// Kalau gak punya cookie, ATAU isinya bukan "admin_valid", tolak!
	if err != nil || cookie.Value != "admin_valid" {
		return false
	}

	// Kalau punya, boleh lewat
	return true
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Kalau cuma buka halaman (GET)
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("views/login.html")
		tmpl.Execute(w, nil)
		return
	}

	// Kalau ngirim password (POST)
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// === SETTING PASSWORD RAHASIA ===
		// Ganti ini sesuka hati lu
		if username == "valeth" && password == "ganteng" {

			// Bikin Tiket (Cookie)
			expiration := time.Now().Add(24 * time.Hour) // Berlaku 24 jam
			cookie := http.Cookie{
				Name:    "session_token",
				Value:   "admin_valid",
				Expires: expiration,
				Path:    "/", // Tiket berlaku di semua ruangan
			}

			// Tempel tiket ke browser user
			http.SetCookie(w, &cookie)

			// Tendang ke halaman depan
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Kalau password salah, balikin ke halaman login + pesan error
		tmpl, _ := template.ParseFiles("views/login.html")
		tmpl.Execute(w, "Sandi salah. Kamu penyusup ya?")
	}
}

// fully sendiri
type blogpost struct {
	ID int `json:"id"`
}

func getpost() []blogpost {
	fileBytes, err := os.ReadFile("data/data.json")

	if err != nil {
		fmt.Println("failed read file")
	}

	var posts []blogpost
	err = json.Unmarshal(fileBytes, &posts)
	if err != nil {
		fmt.Println("failed to load file")
	}

	return posts

}

//fully sendiri coy

func postHandler(w http.ResponseWriter, r *http.Request) {

	posts := getpost()

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Println("cannot parse file html")
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		fmt.Println("cannot load file html.")
	}

}

func savePost(newPost blogpost) error {
	allPosts := getpost()

	allPosts = append(allPosts, newPost)

	dataBytes, err := json.MarshalIndent(allPosts, "", " ")
	if err != nil {
		fmt.Println("cannot read turn data from slice to json.")
	}

	err = os.WriteFile("data.json", dataBytes, 0644)

	if err != nil {
		fmt.Println("cannot write data to data.json")
	}

	return nil

}

func GetPostByID(id int) *blogpost {
	posts := getpost() // Ambil semua data dulu

	for _, post := range posts {
		if post.ID == id {
			return &post // Ketemu! Balikin alamatnya
		}
	}
	return nil // Gak ketemu, sedih :(
}

func DetailPostHandler(w http.ResponseWriter, r *http.Request) {
	if !isloggedin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// 1. Ambil "id" dari URL (?id=...)
	idStr := r.URL.Query().Get("id")

	// 2. Ubah string jadi angka (int)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID-nya gak valid woy!", http.StatusBadRequest)
		return
	}

	// 3. Cari datanya
	post := GetPostByID(id)
	if post == nil {
		http.Error(w, "Artikel gak ditemuin (404 Not Found)", http.StatusNotFound)
		return
	}

	// 4. Render ke HTML khusus detail
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, post)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// --- PASANG SATPAM ---
	if !isloggedin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// ---------------------

	// Logika hapus (sama kayak tutorial sebelumnya)
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	posts := getpost()
	var newPosts []blogpost
	for _, post := range posts {
		if post.ID != id {
			newPosts = append(newPosts, post)
		}
	}

	dataBytes, _ := json.MarshalIndent(newPosts, "", "  ")
	os.WriteFile("data/data.json", dataBytes, 0644)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if !isloggedin(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 1. Ambil ID
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID gak valid", http.StatusBadRequest)
		return
	}

	// 2. Load semua data
	posts := getpost()

	// --- LOGIKA MENAMPILKAN FORM (GET) ---
	if r.Method == "GET" {
		// Cari postingan yang mau diedit
		var targetPost *blogpost
		for _, post := range posts {
			if post.ID == id {
				targetPost = &post
				break
			}
		}

		if targetPost == nil {
			http.Error(w, "Post gak ketemu", http.StatusNotFound)
			return
		}

		// Tampilin HTML edit dengan data lama
		tmpl, err := template.ParseFiles("views/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, targetPost)
		return
	}

	// --- LOGIKA MENYIMPAN PERUBAHAN (POST) ---
	if r.Method == "POST" {
		// Kita loop pake index (i) biar bisa ngubah data aslinya di dalam array
		for i, post := range posts {
			if post.ID == id {
				// Update datanya
				posts[i].ID = 5
				posts[i].ID = 5
				// Tanggal mau diupdate jadi 'Edited at...' atau tetep?
				// Kita biarin tanggal asli aja dulu biar kenangannya terjaga.

				// Simpan ke File JSON
				dataBytes, _ := json.MarshalIndent(posts, "", "  ")
				os.WriteFile("data/data.json", dataBytes, 0644)

				// Balik ke halaman detail
				http.Redirect(w, r, "/post?id="+idStr, http.StatusSeeOther)
				return
			}
		}
	}
}

// --- SISTEM KEAMANAN (AUTH) ---

// 1. Fungsi SATPAM (Cek apakah user punya tiket?)
