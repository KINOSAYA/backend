package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetNewFilms(bearer, language, timeWindow string) {
	var url = fmt.Sprintf("https://api.themoviedb.org/3/trending/movie/%s?language=%s", timeWindow, language)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))
}
