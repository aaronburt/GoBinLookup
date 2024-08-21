package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func searchCSV(filePath, queryBin string) (string, error) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("could not read file: %v", err)
	}

	// Check if there are any records
	if len(records) == 0 {
		return "", fmt.Errorf("no records found in CSV file")
	}

	// Assuming BIN is the first column, skip header row
	header := records[0]
	if len(header) == 0 || strings.ToUpper(header[0]) != "BIN" {
		return "", fmt.Errorf("the first column is not 'BIN' as expected")
	}

	// Find the index of the 'Issuer' column
	issuerIndex := -1
	for i, colName := range header {
		if strings.ToUpper(colName) == "ISSUER" {
			issuerIndex = i
			break
		}
	}
	if issuerIndex == -1 {
		return "", fmt.Errorf("ERROR, couldn't find 'ISSUER' column? Strange")
	}

	// Create a map to store the results
	result := make([]map[string]string, 0)

	// Iterate through the records and find the matching rows
	for _, record := range records[1:] {
		if len(record) > 0 && record[0] == queryBin {
			row := make(map[string]string)
			for i, col := range record {
				if i < len(header) {
					row[header[i]] = col
				}
			}
			result = append(result, row)
		}
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

	row, err := searchCSV("static/bin-list-data.csv", query)

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
