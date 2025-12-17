package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type blogPost struct {
    ID      int    `json:"id"`     
    Title   string `json:"title"`
    Content string `json:"content"`
    Date    string `json:"date"`
}




func savePost(newPost blogPost) error {
	posts := getPost()

	posts=append(posts, newPost)

	dataBytes, err :=json.MarshalIndent(posts, "","  ")
	if (err != nil ){
		return err
	}

	err = os.WriteFile("data/data.json", dataBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getPost() []blogPost{
	fileBytes, err := os.ReadFile("data/data.json")
	
	if err!= nil {
		fmt.Println("Error occured : cannot read file")
		return []blogPost{}

	}

	var posts []blogPost

	err = json.Unmarshal(fileBytes,&posts)

	if err != nil {
		fmt.Println("error occured : failed load file", err)
		return []blogPost{}
	}

	return posts


}



func postsHandler(w http.ResponseWriter, r *http.Request) {

	posts := getPost()

	jsonResponse, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Waduh, server pusing.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.Write(jsonResponse)
}

func main (){
	fmt.Println("server is running.......")

	fmt.Println("=======testing reading file json, below is just a example==========")
	post :=blogPost{ID: 1, Title: "Tes", Content: "Isi", Date: "2025"}
    fmt.Println(post)
	fmt.Println("=========================================")

	data := getPost()
	fmt.Println("file data.json value =")
	fmt.Println(data)

	//saving file baru 
	fmt.Println()
	fmt.Println()


	artikelBaru := blogPost{
        ID:      6,
        Title:   "udah sampe di karawang seru juga",
        Content: "Ternyata gak sesusah itu kalau paham konsepnya.",
        Date:    "2025-12-17",
    }

	fmt.Println("saving blog to data......")
	err := savePost(artikelBaru)

	if err !=nil {
		fmt.Println("failed saving file ", err)
	}else {
		fmt.Println("succes saving file")
	}

	fmt.Println()
	fmt.Println()

	//setelah pembacaan ulang 
	fmt.Println(getPost())
	fmt.Println()
	fmt.Println()
	fmt.Println("server is now starting ")
	fmt.Println("server is starting in local host http://8484")
	fmt.Println("ctrl+c to stop server")

	http.HandleFunc("/", postsHandler)

	http.ListenAndServe(":8484",nil)
	
}



