package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

type BlogPost struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Date    string `json:"date"`
}

func IsLoggedIn(r *http.Request) bool {
    // Coba minta cookie 'session_token' dari browser
    cookie, err := r.Cookie("session_token")
    
    // Kalau gak punya cookie, ATAU isinya bukan "admin_valid", tolak!
    if err != nil || cookie.Value != "admin_valid" {
        return false
    }
    
    // Kalau punya, boleh lewat
    return true
}

// 2. Handler LOGIN (Tempat minta tiket)
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

// 3. Handler LOGOUT (Buang tiket)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Kita timpa tiket lama dengan tiket yang udah kadaluarsa (Expired masa lalu)
    http.SetCookie(w, &http.Cookie{
        Name:    "session_token",
        Value:   "",
        Expires: time.Now().Add(-1 * time.Hour), 
        Path:    "/",
    })
    // Balikin ke home
    http.Redirect(w, r, "/", http.StatusSeeOther)
}



func GetPost() []BlogPost {
    fileBytes, err := os.ReadFile("data/data.json")
    
    if err != nil {
        fmt.Println("Error occured : cannot read file")
        return []BlogPost{}
    }

    var posts []BlogPost
    err = json.Unmarshal(fileBytes, &posts)
    if err != nil {
        fmt.Println("error occured : failed load file", err)
        return []BlogPost{}
    }

    return posts
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
    posts := GetPost()
    fp := "views/index.html"

    tmpl, err := template.ParseFiles(fp)
    if err != nil {
        http.Error(w, "template have error : "+err.Error(), http.StatusInternalServerError)
        return
    }

    err = tmpl.Execute(w, posts)
    if err != nil {
        http.Error(w, "failed render html: "+err.Error(), http.StatusInternalServerError)
    }
}


func SavePost(newPost BlogPost) error {
    posts := GetPost()

    posts = append(posts, newPost)

    dataBytes, err := json.MarshalIndent(posts, "", "  ")
    if err != nil {
        return err
    }

    err = os.WriteFile("data/data.json", dataBytes, 0644)
    if err != nil {
        return err
    }

    return nil
}


func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    if !IsLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("views/create.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }


    if r.Method == "POST" {
        judul := r.FormValue("title")
        isi := r.FormValue("content")

        id := int(time.Now().Unix()) 

        postBaru := BlogPost{
            ID:      id,
            Title:   judul,
            Content: isi,
            Date:    time.Now().Format("2006-01-02"), 
        }

        // 4. Simpen deh!
        err := SavePost(postBaru)
        if err != nil {
            http.Error(w, "Gagal nyimpen curhatan: "+err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Println("Saved making new blog ! ....",http.StatusAccepted)

        // 5. Balikin user ke halaman utama
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }


    
}

func GetPostByID(id int) *BlogPost {
    posts := GetPost() // Ambil semua data dulu
    
    for _, post := range posts {
        if post.ID == id {
            return &post // Ketemu! Balikin alamatnya
        }
    }
    return nil // Gak ketemu, sedih :(
}

func DetailPostHandler(w http.ResponseWriter, r *http.Request) {
    if !IsLoggedIn(r) {
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
    if !IsLoggedIn(r) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    // ---------------------

    // Logika hapus (sama kayak tutorial sebelumnya)
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)
    
    posts := GetPost()
    var newPosts []BlogPost
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
    if !IsLoggedIn(r) {
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
    posts := GetPost()

    // --- LOGIKA MENAMPILKAN FORM (GET) ---
    if r.Method == "GET" {
        // Cari postingan yang mau diedit
        var targetPost *BlogPost
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
                posts[i].Title = r.FormValue("title")
                posts[i].Content = r.FormValue("content")
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


