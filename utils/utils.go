package utils

import (
	"crud-golang/bookStruct"
	"crud-golang/errorStruct"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type urlHandleFunc func(res http.ResponseWriter, req *http.Request) (errorStruct.Error, bool)

func CreateBaseBook(actualId *int, bookMap map[string]bookStruct.Book) {
	idToString := strconv.Itoa(*actualId)
	bookMap[idToString] = bookStruct.Book{ID: idToString, Title: "Harry Potter", Subtitle: "", AuthorName: "J.K.Rowling", ReleaseDate: time.Now().UTC().String(), Price: 10.00, CreatedAt: time.Now().UTC().String()}
	*actualId++
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

func DefaultBookAttributes(book *bookStruct.Book, actualId *int) {
	book.ID = strconv.Itoa(*actualId)
	book.CreatedAt = time.Now().UTC().String()
	*actualId++
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
