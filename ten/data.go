package ten

// [BoardData] represents the json data from the client.
type BoardData struct {
	board                         [9][9]int
	currentBigSquareLegalPosition int
	currentPlayerIdentity         int
}

// [ClientMove] represents the move data from the client.
type ClientMove struct {
	boardBeforeHumanMove *BoardData
	humanMove            *[]int
}

// [serverResponse] is the typed response from the server.
type ServerResponse struct {
	aiMove                        [2]int
	currentBigSquareLegalPosition int
	status                        int
	aiWinningProbability          int
}

// A standard placeholder AI move.
var placeholderMove = [2]int{-1, -1}

// A standard response to an illegal move.
var serverResponseOfIllegalMove = ServerResponse{
	aiMove: placeholderMove,
	currentBigSquareLegalPosition: -1,
	status:               2,
	aiWinningProbability: 0,
}

// Create a [ServerResponse] when the [winner] wins before the AI can move.
func serverResponseOfWinner(winner int) ServerResponse {
	return ServerResponse{
		aiMove: placeholderMove,
		currentBigSquareLegalPosition: -1,
		status:               winner,
		aiWinningProbability: 0,
	}
}
