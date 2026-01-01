package tulisulang

import (
	"encoding/json"
	"fmt"
	"os"
)

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
