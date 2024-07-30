package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	_type := r.URL.Query().Get("type")
	_search := r.URL.Query().Get("search")
	_from := r.URL.Query().Get("from")
	_max := r.URL.Query().Get("max")
	_sort := r.URL.Query().Get("sort")
	// Check if the type is empty
	if _type == "" {
		_type = "match"
	}

	// Create the query
	query := fmt.Sprintf(
		`{
			"search_type": "%s",
			"query": {
				"term": "%s"
			},
			"sort_fields": ["%s"],
			"from": %s,
			"max_results": %s,
			"_source": []
		}`,
		_type, _search, _sort+"date", _from, _max)

	// create the request
	req, err := http.NewRequest("POST", "http://localhost:4080/api/indexer-database/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}

	// Set the headers
	user := os.Getenv("ADMIN")
	pass := os.Getenv("ADMIN_PASS")
	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")

	// Create the client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// close the response
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Fatal(err)
	}

}

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	_from := r.URL.Query().Get("from")
	_max := r.URL.Query().Get("max")
	_sort := r.URL.Query().Get("sort")

	// Create the query
	query := fmt.Sprintf(
		`{
			"search_type": "matchall",
			"query": {
				"term": ""
			},
			"sort_fields": ["%sdate"],
			"from": %s,
			"max_results": %s,
			"_source": []
		}`,
		_sort, _from, _max)

	// create the request
	req, err := http.NewRequest("POST", "http://localhost:4080/api/indexer-database/_search", strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}

	// Set the headers
	user := os.Getenv("ADMIN")
	pass := os.Getenv("ADMIN_PASS")
	req.SetBasicAuth(user, pass)
	req.Header.Set("Content-Type", "application/json")

	// Create the client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// close the response
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Fatal(err)
	}

}
