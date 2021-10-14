package errorStruct

import "net/http"

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func InternalError() Error {
	return Error{Message: `{"error":"Data informed doesn't match the contract"}`, Status: http.StatusInternalServerError}
}

func NotFoundError() Error {
	return Error{Message: `{"error":"Book not found"}`, Status: http.StatusNotFound}
}

func WrongContractError() Error {
	return Error{Message: `{"error":"Data informed doesn't match the contract"}`, Status: http.StatusBadRequest}
}
