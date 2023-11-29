package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetNewFilms(bearer, language, timeWindow string) ([]byte, error) {
	var url = "https://api.kinopoisk.dev/v1.4/list?page=1&limit=100&selectFields=category&selectFields=slug" +
		"&selectFields=cover&notNullFields=category&notNullFields=slug&notNullFields=cover.url" +
		"&notNullFields=cover.previewUrl&category=%D0%A1%D0%B5%D1%80%D0%B8%D0%B0%D0%BB%D1%8B" +
		"&category=%D0%A4%D0%B8%D0%BB%D1%8C%D0%BC%D1%8B"
	req, err := http.NewRequest("GET", url, nil)
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

	fmt.Println(string(body))
	return body, nil
}
