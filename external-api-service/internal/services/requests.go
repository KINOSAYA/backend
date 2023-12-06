package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type ResponseJSON struct {
	Docs []struct {
		Category    string `json:"category"`
		Slug        string `json:"slug"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		MoviesCount int    `json:"moviesCount"`
		Cover       struct {
			URL        string `json:"url"`
			PreviewURL string `json:"previewUrl"`
		} `json:"cover"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	} `json:"docs"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

func GetSlugsByCategory(page, limit int, category string) ([]byte, error) {
	encodedCategory := url.QueryEscape(category)

	var URL = fmt.Sprintf("https://api.kinopoisk.dev/v1.4/list?page=%d&limit=%d&category=%s",
		page, limit, encodedCategory)

	return makeReq(URL)
}

func GetFilmsBySlug(page, limit int, slug string) ([]byte, error) {
	// I apologize for this infinite string
	URL := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie?page=%d&limit=%d&selectFields=id&selectFields=name&selectFields=shortDescription&selectFields=genres&selectFields=poster&selectFields=backdrop&selectFields=videos&selectFields=persons&lists=%s", page, limit, slug)
	return makeReq(URL)
}

type KinopoiskResponse struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	MoviesCount int    `json:"moviesCount"`
	Cover       struct {
		URL        string `json:"url"`
		PreviewURL string `json:"previewUrl"`
	} `json:"cover"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func makeReq(URL string) ([]byte, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	fmt.Println(os.Getenv("apiKey"))
	req.Header.Set("X-API-KEY", os.Getenv("apiKey"))
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}
