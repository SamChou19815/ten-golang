package functions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"developersam.com/ten-golang/ten"
)

// HandleTenAIMoveRequest is an HTTP Cloud Function to handle a client request for an AI move.
func HandleTenAIMoveRequest(writer http.ResponseWriter, request *http.Request) {
	var clientBoard ten.BoardData
	if err := json.NewDecoder(request.Body).Decode(&clientBoard); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "INVALID REQUEST: %s.\n", err.Error())
		return
	}
	response := ten.RespondToClient(&clientBoard)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
