package tulisulang

//fully sendiri
import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
)

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
