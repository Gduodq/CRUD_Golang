package utils

import (
	"context"
	"crud-golang/bookStruct"
	"crud-golang/errorStruct"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type urlHandleFunc func(res http.ResponseWriter, req *http.Request) (errorStruct.Error, bool)

var GetID = uuid.New().String

var MongoConnTimeout = 5 * time.Second

func CreateBaseBook(collection *mongo.Collection, ctx context.Context) {
	book := bookStruct.Book{ID: "3f082af6-d773-4350-9854-1327c91211c7", Title: "Harry Potter", AuthorName: "J.K.Rowling", ReleaseDate: time.Now().UTC().String(), Price: 10.00}
	if _, err := collection.InsertOne(ctx, book); err != nil {
		log.Println("The base book could not be created, skipping the creation")
	}
}

func InitializeQueryNumber(queryParams map[string][]string, key string, defaultValue int) int {
	stringArr, stringExists := queryParams[key]
	if !stringExists {
		return defaultValue
	}
	value, err := strconv.Atoi(stringArr[0])
	if err != nil {
		log.Println("The parsing falied, using default", defaultValue)
		value = defaultValue
	}
	return value
}

func GetBookDataFromReq(req *http.Request, newBook *bookStruct.Book) error {
	err := json.NewDecoder(req.Body).Decode(newBook)
	return err
}

func DefaultBookAttributes(book *bookStruct.Book) {
	book.ID = GetID()
}

func GetReqParam(req *http.Request, key string) string {
	params := mux.Vars(req)
	return params[key]
}

func SetUrlError(urlError *errorStruct.Error, status int, message string) {
	urlError.Status = status
	urlError.Message = message
}

func BaseUrlHandler(handleFunc urlHandleFunc) func(res http.ResponseWriter, req *http.Request) {
	handlerFuncReturn := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		handleError, hasError := handleFunc(res, req)
		if hasError {
			res.WriteHeader(handleError.Status)
			json.NewEncoder(res).Encode(handleError.Message)
		}
	}
	return handlerFuncReturn
}

func ConnectToDB(mongoDBURI string, booksCollection *mongo.Collection) {
	log.Println("Connecting to MongoDB instance on URI " + mongoDBURI)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Println("The API could not connect to the MongoDB instance and can't start up")
		panic(err)
	}
	log.Println("Successfully connected")
	booksCollection = client.Database(`booksGolang`).Collection(`books`)
	ctx, cancel := context.WithTimeout(context.Background(), MongoConnTimeout)
	defer cancel()
	CreateBaseBook(booksCollection, ctx)
}
