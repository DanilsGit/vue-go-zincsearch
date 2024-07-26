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
	"sync"
	"time"

	"github.com/danilsgit/indexerDatabase/constants"
	"github.com/danilsgit/indexerDatabase/models"
)

func main() {

	// start time
	start := time.Now()
	fmt.Println("Start time is: ", start)

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

	// Create a new WaitGroup
	wg := &sync.WaitGroup{}

	// data is the Enron email dataset
	data := "./enron_mail_20110402/maildir"
	// Emails slice
	var pathEmails []string

	// Create a new index in ZincSearch
	createNewIndex()

	// Walk the directory and read the email paths
	err := filepath.Walk(data, func(path string, info os.FileInfo, err error) error {
		// If there is an error, return it
		if err != nil {
			return err
		}

		// If the file is a directory, return nil
		if info.IsDir() {
			return nil
		}

		// Append the path to the pathEmails slice
		pathEmails = append(pathEmails, path)

		// Read the email file
		// email, err := readEmailFile(path)
		// if err != nil {
		// 	fmt.Printf("Error reading email file %s: %v\n", path, err)
		// 	return nil
		// }

		// Append the email to the emails slice
		// emails = append(emails, email)
		// if len(emails) >= 5000 {
		// 	// Add 1 to the WaitGroup
		// 	wg.Add(1)
		// 	// Upload the emails to the database
		// 	go uploadToDatabase(emails, wg)
		// 	emails = nil
		// }

		// Add 1 to the WaitGroup
		// wg.Add(1)
		// Upload the email to the database
		// go uploadToDatabase(email, wg)

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	division := len(pathEmails) / constants.DivisionOfPaths

	if len(pathEmails)%constants.DivisionOfPaths != 0 {
		fmt.Println("Division of paths is not even")
		// stop main
		os.Exit(1)
	}

	if len(pathEmails)%constants.DivisionOfPaths == 0 {
		for i := 0; i < constants.DivisionOfPaths; i++ {
			start := i * division
			end := start + (division - 1)
			// fmt.Printf("Start: %d End: %d\n", start, end)
			wg.Add(1)
			go readPathEmails(pathEmails, start, end, wg)
		}
	}

	// Wait for all the goroutines to finish
	wg.Wait()

	runtime.GC() // Force garbage collection to get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	// Print the number of emails processed and division
	fmt.Printf("Number of emails processed: %d Division of paths: %d\n", len(pathEmails), constants.DivisionOfPaths)
	// end time
	end := time.Now()
	fmt.Println("End time is: ", end)
	// Print the time taken
	duration := time.Since(start).Seconds()
	fmt.Printf("Time taken: %f seconds\n", duration)
	fmt.Println("Done")
}

// readPathEmails reads the emails from the paths
func readPathEmails(pathEmails []string, start int, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	// var emails []models.Email

	for i := start; i <= end; i++ {
		email, err := readEmailFile(pathEmails[i])
		if err != nil {
			fmt.Printf("Error reading email file %s: %v\n", pathEmails[i], err)
			return
		}
		// emails = append(emails, email)
		wg.Add(1)
		go uploadToDatabase(email, wg)
		// emails = nil
	}

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
	req.SetBasicAuth(constants.User, constants.Password)
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
	if strings.Contains(string(body), "already exists") {
		fmt.Println("Index already exists")
		// stop main
		os.Exit(1)
	}
}

// uploadToDatabase uploads the emails to the database
func uploadToDatabase(email models.Email, wg *sync.WaitGroup) {

	defer wg.Done()

	// Create a new Bulk
	// Index is the name of the index in ZincSearch
	// Records is the slice of emails
	// emailData := models.Single{
	// 	Index:  constants.IndexName,
	// 	Record: email,
	// }

	// Encode the Bulk struct to JSON
	jsonData, err := json.Marshal(email)
	if err != nil {
		fmt.Printf("Error encoding to JSON: %v\n", err)
		return
	}

	// Create a new request
	// API bulk is http://localhost:4080/api/_bulkv2
	// API single is http://localhost:4080/api/name/_doc
	req, err := http.NewRequest("POST", "http://localhost:4080/api/"+constants.IndexName+"/_doc", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
	}

	// Set the basic auth and the content type
	req.SetBasicAuth(constants.User, constants.Password)
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
}
