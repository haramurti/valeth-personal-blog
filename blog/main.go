package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type blogPost struct {
    ID      int    `json:"id"`     
    Title   string `json:"title"`
    Content string `json:"content"`
    Date    string `json:"date"`
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

func main (){
	fmt.Println("server is running.......")

	fmt.Println("=======testing reading file json, below is just a example==========")
	post :=blogPost{ID: 1, Title: "Tes", Content: "Isi", Date: "2025"}
    fmt.Println(post)
	fmt.Println("=========================================")

	data := getPost()
	fmt.Println("file data.json value =")
	fmt.Println(data)
	
}



