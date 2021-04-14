package functions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"developersam.com/ten-golang/ten"
)

// HandleTenAIMoveRequest is an HTTP Cloud Function to handle a client request for an AI move.
func HandleTenAIMoveRequest(writer http.ResponseWriter, request *http.Request) {
	// Set CORS headers for the preflight request
	if request.Method == http.MethodOptions {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		writer.Header().Set("Access-Control-Max-Age", "3600")
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	writer.Header().Set("Access-Control-Allow-Origin", "*")

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
