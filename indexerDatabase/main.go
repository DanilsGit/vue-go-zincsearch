package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/danilsgit/indexerDatabase/constants"
	"github.com/danilsgit/indexerDatabase/models"
)

func main() {
	// Create a CPU profile
	f, errCpu := os.Create("cpu_profile.prof")
	if errCpu != nil {
		log.Fatal("could not create CPU profile: ", errCpu)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Create a memory profile

	memFile, errMemory := os.Create("mem_profile.prof")
	if errMemory != nil {
		log.Fatal("could not create memory profile: ", errMemory)
	}
	defer memFile.Close()

	// Read the data from the Enron email dataset
	data := "./enron_mail_20110402/maildir"
	// Emails slice
	var emails []models.Email

	// Create a new index in ZincSearch
	createNewIndex()

	// Walk the directory and read the email files
	err := filepath.Walk(data, func(path string, info os.FileInfo, err error) error {
		// If there is an error, return it
		if err != nil {
			return err
		}

		// If the file is a directory, return nil
		if info.IsDir() {
			return nil
		}

		// Read the email file
		email, err := readEmailFile(path)
		if err != nil {
			fmt.Printf("Error reading email file %s: %v\n", path, err)
			return nil
		}

		emails = append(emails, email)
		if len(emails) >= 2000 {
			uploadToDatabase(emails)
			emails = nil
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	// uploadToDatabase(emails)

	runtime.GC() // Force garbage collection to get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	fmt.Println("Done")
}

// readEmailFile given a file path, it reads the email file and returns an Email struct
func readEmailFile(filePath string) (models.Email, error) {

	// Read the file content
	content, err := os.ReadFile(filePath)

	// If there is an error, return it
	if err != nil {
		return models.Email{}, err
	}

	// Find the end of the header, it is marked by a new line followed by a carriage return
	headerEnd := strings.Index(string(content), "\n\r")

	// If the header end is not found, return an error
	if headerEnd == -1 {
		return models.Email{}, fmt.Errorf("invalid email: %s", filePath)
	}

	// Split the header by new line
	headers := strings.Split(string(content[:headerEnd]), "\n")

	// The body of the email starts after the header
	// +2 to skip the new line and carriage return
	body := string(content[headerEnd+2:])

	// Create a new Email struct
	email := models.Email{}

	// for each header line, set the value in the Email struct
	for _, line := range headers {
		setHeader(line, &email)
	}

	// Set the body of the email
	email.Content = body

	// Return the Email struct
	return email, nil
}

// setHeader given a header line and an Email struct, it sets the value in the Email struct
func setHeader(line string, email *models.Email) {

	// Split the line by the colon
	parts := strings.SplitN(line, ":", 2)

	// If the line does not have a colon or the parts are not 2, return
	if len(parts) != 2 {
		return
	}

	// Trim the spaces from the parts
	headerName := strings.TrimSpace(parts[0])
	headerValue := strings.TrimSpace(parts[1])

	// Switch the header name and set the value in the Email struct
	switch headerName {
	case "Message-ID":
		email.MessageID = headerValue
	case "Date":
		email.Date = headerValue
	case "From":
		email.From = headerValue
	case "To":
		email.To = headerValue
	case "Subject":
		email.Subject = headerValue
	case "Mime-Version":
		email.MimeVersion = headerValue
	case "Content-Type":
		email.ContentType = headerValue
	case "Content-Transfer-Encoding":
		email.ContentTransferEncoding = headerValue
	case "X-From":
		email.XFrom = headerValue
	case "X-To":
		email.XTo = headerValue
	case "X-cc":
		email.Xcc = headerValue
	case "X-bcc":
		email.Xbcc = headerValue
	case "X-Folder":
		email.XFolder = headerValue
	case "X-Origin":
		email.XOrigin = headerValue
	case "X-FileName":
		email.XFilename = headerValue
	}
}

// createNewIndex creates a new index in ZincSearch
func createNewIndex() {

	// Read the JSON file
	index, err := os.ReadFile("./data.json")

	// If there is an error, print it and return
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Create a new request
	req, err := http.NewRequest("POST", "http://localhost:4080/api/index", bytes.NewReader(index))
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
	}

	// Set the basic auth and the content type
	req.SetBasicAuth("admin", "admin123")
	req.Header.Set("Content-Type", "application/json")

	// Create a new client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)

	// If there is an error, print it and return
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Close the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)

	// If there is an error, print it and return
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println(string(body))
}

// uploadToDatabase uploads the emails to the database
func uploadToDatabase(emails []models.Email) {

	// Create a new Bulk struct
	// Index is the name of the index in ZincSearch
	// Records is the slice of emails
	emailData := models.Bulk{
		Index:   constants.IndexName,
		Records: emails,
	}

	// Encode the Bulk struct to JSON
	jsonData, err := json.Marshal(emailData)
	if err != nil {
		fmt.Printf("Error encoding to JSON: %v\n", err)
		return
	}

	// Create a new request
	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
	}

	// Set the basic auth and the content type
	req.SetBasicAuth("admin", "admin123")
	req.Header.Set("Content-Type", "application/json")

	// Create a new client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Close the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err, string(body))
		return
	}

	// Print the response body
	// fmt.Println(string(body))
}
