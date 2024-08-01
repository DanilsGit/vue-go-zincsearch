# VUE - GO - ZINCSEARCH - DOCKER

## About
This project consists of a Zincsearch database, a backend in Go, and a frontend in Vue. A Go script is used to index the data into Zincsearch and display the content in the frontend.

## Performance Documentation
[See Performance Improvements](ImprovePerformance.pdf)

## Test the App

1. Clone the project.
2. Go to `/indexerDatabase/emails` (create the `emails` folder if it does not exist).
3. Download and unzip the data from [http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz).
4. Go to `/frontend` and create a `.env` file with the following content: `VITE_API_URL=http://localhost:8080`.
5. Go to the root directory.
6. Run `docker compose build`.
7. Run `docker compose up -d`.
8. Go to `/indexerDatabase/constants/constants.go` and configure the workers and subdivisions according to your machine.
9. Go to `/indexerDatabase` and execute go run main.go. Please note, this process may take some time. If an error occurs during execution, try reducing the subdivisions and workers in constants.go.
10. Test the app.

Frontend is available at port 5173  
Backend is available at port 8080  
Zincsearch is available at port 4080
