package main

import (
	"fmt"
	"net/http"

	"valeth-personal-blog/handlers"
)


func main (){
    fmt.Println("server is running.......")

    fmt.Println("=======testing reading file json, below is just a example==========")
    
    post := handlers.BlogPost{ID: 1, Title: "Tes", Content: "Isi", Date: "2025"}
    fmt.Println(post)
    fmt.Println("=========================================")

    data := handlers.GetPost()
    fmt.Println("file data.json value =")
    fmt.Println(data)

    fmt.Println()
    fmt.Println()
    fmt.Println("=========================================")


    // artikelBaru := handlers.BlogPost{
    //     ID:      10,
    //     Title:   "ayam",
    //     Content: "Ternyayamlakdjf;lakjd f;kadjsa.",
    //     Date:    "2025-12-17",
    // }

    // fmt.Println("saving blog to data......")
    
    // err := handlers.SavePost(artikelBaru)

    // if err !=nil {
    //     fmt.Println("failed saving file ", err)
    // }else {
    //     fmt.Println("succes saving file")
    // }

    fmt.Println()
    fmt.Println()


    fmt.Println(handlers.GetPost())
    
    fmt.Println()
    fmt.Println()
    fmt.Println("server is now starting ")
    fmt.Println("server is starting in local host http://localhost:8484") 
    fmt.Println("ctrl+c to stop server")


    http.HandleFunc("/", handlers.PostsHandler)
	http.HandleFunc("/create", handlers.CreatePostHandler)

    http.ListenAndServe(":8484",nil)   
}