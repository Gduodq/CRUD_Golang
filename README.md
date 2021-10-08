# CRUD_Golang

This project is a Rest API to manage books. The database is abstracted as a map.

This API is built with [Golang](https://golang.org).

To run it, make sure to install all dependencies first. Enter the project folder and run:

### `go get ./...`

After the installation you can run it with:

### `go run .\main.go`

The command above have an optional argument which is the port where the API will listen. To change it, simply run:

### `go run .\main.go {PORT}`

If no port is provided the API will listen by default on port 8000.
