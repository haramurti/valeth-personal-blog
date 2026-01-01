// package tulisUlang

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// type blogPost struct {
// 	ID int `json:"id"`
// }

// func getPost() []blogPost {

// 	fileBytes, err := os.ReadFile("data/data.json")

// 	if err != nil {
// 		fmt.Println("file cannot be read ")
// 		return []blogPost{}
// 	}

// 	var posts []blogPost
// 	err = json.Unmarshal(fileBytes, &posts)
// 	if err != nil{
// 		fmt.Println("failed to load file")
// 		return []blogPost{}
// 	}

// 	return posts
// }

package tulisulang

import "os"

type blogpost struct {
	ID int `json:"id"`
}

func getpost() []blogpost {
	fileBytes, err := os.ReadFile("data/data.json")

	if err != nil {
		fmt.println
	}

}
