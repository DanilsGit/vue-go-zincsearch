package main

import (
	"bytes"
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
	utils "github.com/danilsgit/indexerDatabase/utils"
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
	data := constants.Directory

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

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
	}

	// create semaphore
	sem := make(chan struct{}, constants.Workers)

	// Slice of emails divided in n parts
	rangeOfParts := utils.RangeOfParts(len(pathEmails), constants.Subdivisions)

	// For each part, read the emails
	for i := 0; i < len(rangeOfParts)-1; i++ {
		wg.Add(1)
		go utils.ReadPathEmails(pathEmails, rangeOfParts[i], rangeOfParts[i+1]-1, wg, sem)
	}

	// Wait for all the goroutines to finish
	wg.Wait()

	runtime.GC() // Force garbage collection to get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	// Print the number of emails processed and division
	fmt.Printf("Number of emails processed: %d\n", len(pathEmails))
	// end time
	end := time.Now()
	fmt.Println("End time is: ", end)
	// Print the time taken
	duration := time.Since(start).Seconds()
	fmt.Printf("Time taken: %f seconds\n", duration)
	fmt.Println("Done")
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

	// Get the credentials from the environment variables
	admin := os.Getenv("ADMIN")
	admin_pass := os.Getenv("ADMIN_PASS")

	fmt.Println(admin, admin_pass)

	if admin == "" || admin_pass == "" {
		fmt.Println("Please set the ADMIN and ADMIN_PASS environment")
		admin = "admin"
		admin_pass = "admin123"
	}

	// Set the basic auth and the content type
	req.SetBasicAuth(admin, admin_pass)
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

		// get in console if the user wants to delete the index
		fmt.Println("Do you want to delete the index? (yes/no)")
		var input string
		fmt.Scanln(&input)

		// if the user wants to delete the index
		if input == "yes" {
			deleteIndex()
		}
		// stop main
		os.Exit(1)
	}
}

func deleteIndex() {
	// Get the credentials from the environment variables
	admin := os.Getenv("ADMIN")
	admin_pass := os.Getenv("ADMIN_PASS")

	// Create a new request
	req, err := http.NewRequest("DELETE", "http://localhost:4080/api/index/"+constants.IndexName, nil)
	if err != nil {
		fmt.Printf("Error creating the request: %v\n", err)
	}

	// Set the basic auth
	req.SetBasicAuth(admin, admin_pass)

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
