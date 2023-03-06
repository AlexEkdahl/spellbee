package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	apiUrl = "https://api.textgears.com/grammar?key="
	vowels = "aeiou"
	dbFile = "cache.db"
	dbDir  = "spellbee"
)

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
	initFlag := flag.Bool("init", false, "Initialize the database")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting the user's config directory:", err)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".config")
	dbPath := filepath.Join(configDir, dbDir, dbFile)
	if *initFlag {
		initializeDatabase(dbPath)
		return
	}

	word := getWord()
	article := getArticle(word)
	apiKey := getApiKey()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening the database:", err)
		os.Exit(1)
	}
	defer db.Close()

	correctArticle := checkArticle(db, apiKey, word, article)

	fmt.Printf("The word '%s' should be preceded by '%s'.\n", word, correctArticle)
}

func getWord() string {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a word as an argument.")
		os.Exit(1)
	}
	return os.Args[1]
}

func getArticle(word string) string {
	article := "A"
	if strings.ContainsAny(vowels, string(word[0])) {
		article = "An"
	}
	return article
}

func getApiKey() string {
	apiKey := os.Getenv("TEXT_GEARS_API_KEY")
	return apiKey
}

func initializeDatabase(dbPath string) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		fmt.Println("Error creating the database directory:", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening the database:", err)
		os.Exit(1)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cache (word TEXT PRIMARY KEY, article TEXT)`)
	if err != nil {
		fmt.Println("Error creating the cache table:", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully.")
}

func checkArticle(db *sql.DB, apiKey, word, article string) string {
	var cachedArticle string
	err := db.QueryRow(`SELECT article FROM cache WHERE word = ?`, word).Scan(&cachedArticle)
	if err == nil {
		return cachedArticle
	}

	correctArticle := checkArticleWithApi(apiKey, word, article)
	_, err = db.Exec(`INSERT INTO cache (word, article) VALUES (?, ?)`, word, correctArticle)
	if err != nil {
		fmt.Println("Error inserting into the cache table:", err)
		os.Exit(1)
	}

	return correctArticle
}

func checkArticleWithApi(apiKey, word, article string) string {
	urlBuffer := bytes.NewBufferString(apiUrl)
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
