package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/danilsgit/indexerDatabase/constants"
	"github.com/danilsgit/indexerDatabase/models"
)

var client = &http.Client{} // client is a pointer to an http.Client

// RangeOfParts given a number and the initial number of parts, it returns a slice with the range of parts
func RangeOfParts(number, initParts int) []int {

	// First | divide the number by the initial number of parts
	equalDivision := make([]int, initParts)
	// Calculate the division and the remainder
	part := number / initParts
	rem := number % initParts

	// For each part, set the value
	for i := 0; i < initParts; i++ {
		equalDivision[i] = part
	}

	// Add the remainder to the first parts
	for i := 0; i < rem; i++ {
		equalDivision[i]++
	}

	//  Second | create a response slice with the range of parts
	response := make([]int, len(equalDivision)+1)
	response[0] = 0
	for i := 1; i < len(equalDivision)+1; i++ {
		response[i] = response[i-1] + equalDivision[i-1]
	}
	return response
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

// readPathEmails reads the emails from the paths
func ReadPathEmails(pathEmails []string, start int, end int, wg *sync.WaitGroup, sem chan struct{}) {

	defer wg.Done()

	var emails []models.Email

	for i := start; i <= end; i++ {
		email, err := readEmailFile(pathEmails[i])
		if err != nil {
			fmt.Printf("Error reading email file %s: %v\n", pathEmails[i], err)
			return
		}
		emails = append(emails, email)
		// wg.Add(1)
		// go uploadToDatabase([]models.Email{email}, wg)
	}

	// // Add
	// wg.Add(1)
	// go uploadToDatabase(emails, wg)
	sem <- struct{}{}
	wg.Add(1)
	go func() {
		defer func() {
			<-sem // Release
			wg.Done()
		}()
		uploadToDatabase(emails)
	}()
}

// uploadToDatabase uploads the emails to the database
func uploadToDatabase(emails []models.Email) {

	// Create a new Bulk
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
	// API bulk is http://localhost:4080/api/_bulkv2
	// API single is http://localhost:4080/api/"+constants.IndexName+"/_doc"
	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
	}

	// Set the basic auth and the content type
	req.SetBasicAuth(constants.Admin, constants.Password)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// Close the response body
	defer resp.Body.Close()
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
	header := strings.TrimSpace(parts[0])
	header_value := strings.TrimSpace(parts[1])

	// Switch the header name and set the value in the Email struct
	switch header {
	case "Message-ID":
		email.MessageID = header_value
	case "Date":
		parsedDate := parseDateWithFormats(header_value)
		email.Date = parsedDate
	case "From":
		email.From = header_value
	case "To":
		email.To = header_value
	case "Subject":
		email.Subject = header_value
	case "Mime-Version":
		email.MimeVersion = header_value
	case "Content-Type":
		email.ContentType = header_value
	case "Content-Transfer-Encoding":
		email.ContentTransferEncoding = header_value
	case "X-From":
		email.XFrom = header_value
	case "X-To":
		email.XTo = header_value
	case "X-cc":
		email.Xcc = header_value
	case "X-bcc":
		email.Xbcc = header_value
	case "X-Folder":
		email.XFolder = header_value
	case "X-Origin":
		email.XOrigin = header_value
	case "X-FileName":
		email.XFilename = header_value
	}
}

// parseDateWithFormats given a date string, it tries to parse it with multiple formats
func parseDateWithFormats(date string) time.Time {
	formats := []string{
		"Mon, _2 Jan 2006 15:04:05 -0700 (MST)",
		"Monday, January 2, 2006",
		"Monday, March 12",
	}

	var parsedDate time.Time
	var err error
	for _, layout := range formats {
		if len(date) == len("Monday, March 12") {
			date = date + ", 2000"
		}
		parsedDate, err = time.Parse(layout, date)
		if err == nil {
			return parsedDate
		}
	}

	fmt.Printf("Error parsing date: %v\n", err)
	fmt.Printf("Date: %s\n", date)
	// stop main
	os.Exit(1)
	return time.Time{}
}
