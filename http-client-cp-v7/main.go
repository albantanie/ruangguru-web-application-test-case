package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Quotes struct {
	Tags   []string `json:"tags"`
	Author string   `json:"author"`
	Quote  string   `json:"content"`
}

func ClientGet() ([]Quotes, error) {
	resp, err := http.Get("https://api.quotable.io/quotes/random?limit=3")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	// fmt.Println(string(body))
	quote := []Quotes{}
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return nil, err
	}
	// fmt.Println(quote[0].Quote)

	// Hit API https://api.quotable.io/quotes/random?limit=3 with method GET:
	return quote, nil // TODO: replace this
}

type data struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Postman struct {
	Data data   `json:"data"`
	Url  string `json:"url"`
}

func ClientPost() (Postman, error) {
	postBody, _ := json.Marshal(map[string]string{
		"name":  "Dion",
		"email": "dionbe2022@gmail.com",
	})
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("https://postman-echo.com/post", "application/json", requestBody)
	if err != nil {
		return Postman{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Postman{}, err
	}

	postman := Postman{}
	err = json.Unmarshal(body, &postman)
	if err != nil {
		return Postman{}, err
	}

	// Hit API https://postman-echo.com/post with method POST:
	return postman, nil // TODO: replace this
}

func main() {
	get, _ := ClientGet()
	fmt.Println(get)

	post, _ := ClientPost()
	fmt.Println(post)
}
