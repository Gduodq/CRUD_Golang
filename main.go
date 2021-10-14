package main

import (
	"context"
	"crud-golang/bookStruct"
	"crud-golang/errorStruct"
	"crud-golang/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
)

var booksCollection *mongo.Collection

func getBooks(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	queryParams := req.URL.Query()
	limit := utils.InitializeQueryNumber(queryParams, "limit", 0)
	skip := utils.InitializeQueryNumber(queryParams, "skip", 0)
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if skip > 0 {
		opts.SetSkip(int64(skip))
	}
	var booksFound []bookStruct.Book
	ctx, cancel := context.WithTimeout(req.Context(), utils.MongoConnTimeout)
	defer cancel()
	findCursor, err := booksCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		urlError = errorStruct.InternalError()
	}
	if err = findCursor.All(ctx, &booksFound); err != nil {
		urlError = errorStruct.InternalError()
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(booksFound)
	return
}

func getBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	var book bookStruct.Book
	ctx, cancel := context.WithTimeout(req.Context(), utils.MongoConnTimeout)
	defer cancel()
	if err := booksCollection.FindOne(ctx, bson.M{"_id": bookId}).Decode(&book); err != nil {
		if err == mongo.ErrNoDocuments {
			urlError = errorStruct.NotFoundError()
		} else {
			urlError = errorStruct.InternalError()
		}
		hasError = true
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(book)
	return
}

func createBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	var newBook bookStruct.Book
	if err := utils.GetBookDataFromReq(req, &newBook); err != nil {
		urlError = errorStruct.WrongContractError()
		hasError = true
		return
	}
	utils.DefaultBookAttributes(&newBook)
	ctx, cancel := context.WithTimeout(req.Context(), utils.MongoConnTimeout)
	defer cancel()
	if _, err := booksCollection.InsertOne(ctx, newBook); err != nil {
		urlError = errorStruct.InternalError()
		hasError = true
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(newBook)
	return
}

func replaceBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	var newBook bookStruct.Book
	if err := utils.GetBookDataFromReq(req, &newBook); err != nil {
		urlError = errorStruct.WrongContractError()
		hasError = true
		return
	}
	newBook.ID = bookId
	ctx, cancel := context.WithTimeout(req.Context(), utils.MongoConnTimeout)
	defer cancel()
	err := booksCollection.FindOneAndReplace(ctx, bson.M{"_id": bookId}, newBook).Decode(&bookStruct.Book{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			urlError = errorStruct.NotFoundError()
		} else {
			urlError = errorStruct.InternalError()
		}
		hasError = true
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(newBook)
	return
}

func deleteBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	var book bookStruct.Book
	ctx, cancel := context.WithTimeout(req.Context(), utils.MongoConnTimeout)
	defer cancel()
	err := booksCollection.FindOneAndDelete(ctx, bson.M{"_id": bookId}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			urlError = errorStruct.NotFoundError()
		} else {
			urlError = errorStruct.InternalError()
		}
		hasError = true
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(book)
	return
}

func initializeEnvVars(envVarsKeyDefault [][]string, envVarsPointers []*string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	var envKeyValue string
	for idx, envVarKeyDefault := range envVarsKeyDefault {
		key := envVarKeyDefault[0]
		defaultValue := envVarKeyDefault[1]
		if value := os.Getenv(key); value != "" {
			envKeyValue = value
		} else {
			envKeyValue = defaultValue
		}
		*envVarsPointers[idx] = envKeyValue
	}
}

func startAPI(PORT string) {
	log.Println("Starting the HTTP server...")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/books", utils.BaseUrlHandler(getBooks)).Methods("GET")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(getBook)).Methods("GET")
	router.HandleFunc("/books", utils.BaseUrlHandler(createBook)).Methods("POST")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(replaceBook)).Methods("PUT")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(deleteBook)).Methods("DELETE")
	log.Println("HTTP server listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}

func main() {
	var PORT, MONGOHOST, MONGOPORT, MONGOUSER, MONGOPASS string
	envVarsKeyDefault := [][]string{{"PORT", "8000"}, {"MONGOHOST", "localhost"}, {"MONGOPORT", "27017"}, {"MONGOUSER", ""}, {"MONGOPASS", ""}}
	envVarsPointers := []*string{&PORT, &MONGOHOST, &MONGOPORT, &MONGOUSER, &MONGOPASS}
	initializeEnvVars(envVarsKeyDefault, envVarsPointers)
	if MONGOUSER != "" {
		MONGOUSER = MONGOUSER + ":"
	}
	if MONGOPASS != "" {
		MONGOPASS = MONGOPASS + "@"
	}
	mongoDBURI := fmt.Sprintf(`mongodb://%v%v%v:%v`, MONGOUSER, MONGOPASS, MONGOHOST, MONGOPORT)
	utils.ConnectToDB(mongoDBURI, booksCollection)
	startAPI(PORT)
}
