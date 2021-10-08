package main

import (
	"crud-golang/bookStruct"
	"crud-golang/errorStruct"
	"crud-golang/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var bookMap = make(map[string]bookStruct.Book)

var id int = 1

func getBooks(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	queryParams := req.URL.Query()
	limit := utils.InitializeQueryNumber(queryParams, "limit", len(bookMap))
	skip := utils.InitializeQueryNumber(queryParams, "skip", 0)
	startIdx := skip
	endIdx := skip + limit
	data := make([]bookStruct.Book, 0, limit)
	idx := 0
	for _, book := range bookMap {
		if idx >= startIdx && idx < endIdx {
			data = append(data, book)
		}
		idx++
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(data)
	return
}

func getBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	book, bookExists := bookMap[bookId]
	if !bookExists {
		utils.SetUrlError(&urlError, http.StatusNotFound, `{"error":"Book not found"}`)
		hasError = true
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(book)
	}
	return
}

func createBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	var newBook bookStruct.Book
	err := utils.GetBookDataFromReq(req, &newBook)
	if err != nil {
		utils.SetUrlError(&urlError, http.StatusBadRequest, `{"error":"Data informed doesn't match the contract"}`)
		hasError = true
		return
	}
	utils.DefaultBookAttributes(&newBook, &id)
	bookMap[newBook.ID] = newBook
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(newBook)
	return
}

func updateBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	_, bookExists := bookMap[bookId]
	if !bookExists {
		utils.SetUrlError(&urlError, http.StatusNotFound, `{"error":"Book not found"}`)
		hasError = true
		return
	} else {
		var newBook bookStruct.Book
		err := utils.GetBookDataFromReq(req, &newBook)
		if err != nil {
			utils.SetUrlError(&urlError, http.StatusBadRequest, `{"error":"Data informed doesn't match the contract"}`)
			hasError = true
			return
		}
		newBook.ID = bookId
		bookMap[bookId] = newBook
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(newBook)
	}
	return
}

func deleteBook(res http.ResponseWriter, req *http.Request) (urlError errorStruct.Error, hasError bool) {
	bookId := utils.GetReqParam(req, "bookId")
	book, bookExists := bookMap[bookId]
	if !bookExists {
		utils.SetUrlError(&urlError, http.StatusNotFound, `{"error":"Book not found"}`)
		hasError = true
		return
	} else {
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(book)
		delete(bookMap, bookId)
	}
	return
}

func main() {
	commandLineArgs := os.Args[1:]
	var PORT string
	if len(commandLineArgs) > 0 {
		PORT = commandLineArgs[0]
	} else {
		PORT = "8000"
	}
	log.Println("Starting the HTTP server on port " + PORT)
	router := mux.NewRouter().StrictSlash(true)
	utils.CreateBaseBook(&id, bookMap)
	router.HandleFunc("/books", utils.BaseUrlHandler(getBooks)).Methods("GET")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(getBook)).Methods("GET")
	router.HandleFunc("/books", utils.BaseUrlHandler(createBook)).Methods("POST")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(updateBook)).Methods("PUT")
	router.HandleFunc("/books/{bookId}", utils.BaseUrlHandler(deleteBook)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
