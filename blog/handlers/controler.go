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
    if r.Method == "GET" {
        tmpl, err := template.ParseFiles("views/create.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tmpl.Execute(w, nil)
        return
    }

    // Kalau user NGIRIM data (POST)
    if r.Method == "POST" {
        // 1. Ambil data dari form HTML
        judul := r.FormValue("title")
        isi := r.FormValue("content")

        // 2. Bikin ID acak (biar gampang dulu, nanti kita rapihin)
        // Note: Aslinya jangan gini ya, ini cara males tapi jalan :D
        id := int(time.Now().Unix()) 

        // 3. Masukin ke struct
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

