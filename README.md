# CRUD_Golang

This project is a Rest API to manage books.

This API is built with [Golang](https://golang.org) and the used database is [MongoDB](www.mongodb.com/) (To configure the access check the [Enviroment Variables](#EnviromentVariables) section.).

To run it, make sure to install all dependencies first. Enter the project folder and run:

### `go get ./...`

After the installation you can run it with:

### `go run .\main.go`

The API will listen by default on [http://localhost:8000](http://localhost:8000). To change the default port check the [Enviroment Variables](#EnviromentVariables) section.

## <a name="EnviromentVariables"></a>Enviroment Variables

To set the enviroment variables it's necessary to create a `.env` file on the project root with the following optional variables:

- PORT:

  Set the port where the application should run.

  Default = 8000

- MONGOHOST:

  Set the host where the MongoDB instance should be running.

  Default = localhost

- MONGOPORT:

  Set the port where the MongoDB instance should be running.

  Default = 27017

- MONGOUSER:

  Set the username for MongoDB authetication.

  Default = ""

- MONGOPASS:

  Set the password for MongoDB authetication.

  Default = ""
