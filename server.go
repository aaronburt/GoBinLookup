package main

//DUE TO SQLITE3 by Mattn it requires C (GCC installed), because it uses CGO  $env:CGO_ENABLED=1; go build server.go

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func searchBinAndLog(databasePath string, binValue string) (string, error) {

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return "", fmt.Errorf("could not open database: %v", err)
	}

	defer db.Close()

	// Prepare and execute the SQL query, prepared statements!!!!
	query := "SELECT * FROM bin WHERE BIN = ?"
	rows, err := db.Query(query, binValue)
	if err != nil {
		return "", fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("could not retrieve columns: %v", err)
	}

	// Prepare a slice to hold column values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Create a slice to hold the result
	var result []map[string]interface{}

	// Fetch the result
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return "", fmt.Errorf("could not scan row: %v", err)
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			row[colName] = values[i]
		}
		result = append(result, row)
	}

	// Handle no records found
	if len(result) == 0 {
		return "", fmt.Errorf("no record found with BIN = %s", binValue)
	}

	// Convert the result to JSON
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("could not convert result to JSON: %v", err)
	}

	return string(jsonResult), nil
}

func faviconHandler(response http.ResponseWriter, request *http.Request) {
	faviconPath := "static/favicon.ico"
	http.ServeFile(response, request, faviconPath)
}

func issuerHandler(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("bin")
	if query == "" {
		fmt.Fprintf(response, "Enter a value")
	}

	if len(query) != 6 {
		fmt.Fprintf(response, "Expected 6 digits")
		return
	}

	row, err := searchBinAndLog("static/sqlite.db", query)

	if err != nil {
		fmt.Fprintf(response, "Error: %v", err)
		return
	}

	fmt.Fprintf(response, "%s", row)
}

func homeHandler(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "static/index.html")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/issuer", issuerHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)

	port := os.Getenv("PORT")
	fmt.Printf("Server is starting on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
