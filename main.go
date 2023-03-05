package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const apiUrl = "https://api.textgears.com/grammar?key="

type GrammarResult struct {
	Status   bool `json:"status"`
	Response struct {
		Errors []struct {
			ID     string   `json:"id"`
			Offset int      `json:"offset"`
			Length int      `json:"length"`
			Bad    string   `json:"bad"`
			Better []string `json:"better"`
			Type   string   `json:"type"`
		} `json:"errors"`
	} `json:"response"`
}

func main() {
	word := getWord()
	article := "A"
	vowels := "aeiou"

	if strings.ContainsAny(vowels, string(word[0])) {
		article = "An"
	}

	apiKey := getApiKey()
	correctArticle := checkArticleWithApi(apiKey, word, article)

	fmt.Printf("The word '%s' should be preceded by '%s'.\n", word, correctArticle)
}

func getWord() string {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a word as an argument.")
		os.Exit(1)
	}
	return os.Args[1]
}

func getApiKey() string {
	apiKey := os.Getenv("TEXT_GEARS_API_KEY")
	return apiKey
}

func checkArticleWithApi(apiKey, word, article string) string {
	var urlBuffer bytes.Buffer
	urlBuffer.WriteString(apiUrl)
	urlBuffer.WriteString(apiKey)
	urlBuffer.WriteString("&text=")
	urlBuffer.WriteString(article)
	urlBuffer.WriteString("+")
	urlBuffer.WriteString(word)
	urlBuffer.WriteString("&language=en-GB")

	response, err := http.Get(urlBuffer.String())
	if err != nil {
		fmt.Println("Error calling the external API:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	var result GrammarResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding the API response:", err)
		os.Exit(1)
	}

	if result.Status {
		for _, err := range result.Response.Errors {
			if err.Bad == article {
				article = err.Better[0]
				break
			}
		}
	}

	return strings.ToLower(article)
}
