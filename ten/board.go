package ten

import "fmt"

// Move represents a simple move in TEN.
type Move struct {
	a int
	b int
}

/*
	In variable names, a big square refers to a 3*3 square;
	a tile refers to a 1*1 square.
	Each tile is either 1, -1 or 0 (black, white, empty).
*/

// Board is the concrete type of the board game TEN.
type Board struct {
	// board represents the entire board's occupation information.
	tiles [81]int8
	// bigSquaresStatus keeps track of winning progress on big squares.
	// 1, -1, 0, 2 mean black wins, white wins, inconclusive, and all occupied.
	bigSquaresStatus [9]int8
	// bigSquareToPick is the current legal big square to pick as next move.
	// If it's value is -1, that means the user can place the move everywhere.
	// This variable is important for maintaining the current game state.
	bigSquareToPick int
	// blackBigSquaresCounter is the winning big squares counter for black.
	blackBigSquaresCounter int
	// whiteBigSquaresCounter is the winning big squares counter for white.
	whiteBigSquaresCounter int
	// playerIdentity must be 1 or -1.
	playerIdentity int8
}

// BoardData is a struct designed to be sent from the client side that can completely decides a
// game state, although it does not record enough information to efficiently simulate the game.
type BoardData struct {
	Tiles           []int8 `json:"tiles"`
	BigSquareToPick int    `json:"bigSquareToPick"`
	PlayerIdentity  int8   `json:"playerIdentity"`
}

// makeBoard returns an initial empty board.
func makeBoard() *Board {
	var rawTiles [81]int8
	var bigSquaresStatus [9]int8
	return &Board{
		tiles:                  rawTiles,
		bigSquaresStatus:       bigSquaresStatus,
		bigSquareToPick:        -1,
		blackBigSquaresCounter: 0,
		whiteBigSquaresCounter: 0,
		playerIdentity:         1,
	}
}

func makeBoardFromData(data *BoardData) *Board {
	var tiles [81]int8
	copy(tiles[:], data.Tiles)
	var bigSquaresStatus [9]int8
	blackBigSquaresCounter := 0
	whiteBigSquaresCounter := 0
	for i := 0; i < 9; i++ {
		status := computeSquareStatus(tiles[:], i*9)
		bigSquaresStatus[i] = status
		if status == 1 {
			blackBigSquaresCounter++
		} else if status == -1 {
			whiteBigSquaresCounter++
		}
	}
	return &Board{
		tiles:                  tiles,
		bigSquaresStatus:       bigSquaresStatus,
		bigSquareToPick:        data.BigSquareToPick,
		blackBigSquaresCounter: blackBigSquaresCounter,
		whiteBigSquaresCounter: whiteBigSquaresCounter,
		playerIdentity:         data.PlayerIdentity,
	}
}

// playerSimplyWinSquare performs a naive check on the square about whether the player with id wins
// the square. It only checks according to the primitive tic-tac-toe rule.
func playerSimplyWinSquare(square []int8, offset int, id int8) bool {
	return square[offset] == id && square[offset+1] == id && square[offset+2] == id ||
		square[offset+3] == id && square[offset+4] == id && square[offset+5] == id ||
		square[offset+6] == id && square[offset+7] == id && square[offset+8] == id ||
		square[offset] == id && square[offset+3] == id && square[offset+6] == id ||
		square[offset+1] == id && square[offset+4] == id && square[offset+7] == id ||
		square[offset+2] == id && square[offset+5] == id && square[offset+8] == id ||
		square[offset] == id && square[offset+4] == id && square[offset+8] == id ||
		square[offset+2] == id && square[offset+4] == id && square[offset+6] == id
}

// computeSquareStatus helps to determine whether a square with offset belongs to black or white.
// If all tiles are occupied, it returns 2; else (there is no direct victory), it returns 0.
func computeSquareStatus(square []int8, offset int) int8 {
	if playerSimplyWinSquare(square, offset, 1) {
		return 1
	}
	if playerSimplyWinSquare(square, offset, -1) {
		return -1
	}
	for i := 0; i < 9; i++ {
		if square[offset+i] == 0 {
			// There is a space left.
			return 0
		}
	}
	return 2
}

/*
 * Check whether a move is legal, where move is given by ([a], [b]).
 */
func (board *Board) isLegalMove(a, b int) bool {
	if a < 0 || a > 8 || b < 0 || b > 8 {
		// Out of boundary values
		return false
	}
	if board.bigSquareToPick != -1 && board.bigSquareToPick != a {
		// in the wrong big square when it cannot have a free move
		return false
	}
	// not in the occupied big square and on an empty tile
	return board.bigSquaresStatus[a] == 0 && board.tiles[a*9+b] == 0
}

// allLegalMovesForAI returns a list of all legal moves for AI.
func (board *Board) allLegalMovesForAI() []*Move {
	list := make([]*Move, 0, 40) // 40 will be enough for most cases.
	bigSquareToPick := board.bigSquareToPick
	if bigSquareToPick == -1 {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board.isLegalMove(i, j) {
					list = append(list, &Move{a: i, b: j})
				}
			}
		}
	} else {
		for j := 0; j < 9; j++ {
			// Can only move in the specified square
			if board.isLegalMove(bigSquareToPick, j) {
				list = append(list, &Move{a: bigSquareToPick, b: j})
			}
		}
	}
	return list
}

// makeMoveWithoutCheck returns a new board with the given move applied, without doing any check.
func (board *Board) makeMoveWithoutCheck(move *Move) *Board {
	// Directly make the move
	var newTiles [81]int8
	copy(newTiles[:], board.tiles[:])
	a := move.a
	b := move.b
	newTiles[a*9+b] = board.playerIdentity

	// Update big squares
	var newBigSquareStatusArray [9]int8
	copy(newBigSquareStatusArray[:], board.bigSquaresStatus[:])
	newBigSquareStatus := computeSquareStatus(newTiles[:], a*9)
	newBigSquareStatusArray[a] = newBigSquareStatus
	var newBigSquareToPick int
	if newBigSquareStatusArray[b] == 0 {
		newBigSquareToPick = b
	} else {
		newBigSquareToPick = -1
	}

	// Compute counter
	blackBigSquaresCounter := board.blackBigSquaresCounter
	whiteBigSquaresCounter := board.whiteBigSquaresCounter
	if newBigSquareStatus == 1 {
		blackBigSquaresCounter++
	} else if newBigSquareStatus == -1 {
		whiteBigSquaresCounter++
	}

	return &Board{
		tiles:                  newTiles,
		bigSquaresStatus:       newBigSquareStatusArray,
		bigSquareToPick:        newBigSquareToPick,
		blackBigSquaresCounter: blackBigSquaresCounter,
		whiteBigSquaresCounter: whiteBigSquaresCounter,
		playerIdentity:         -board.playerIdentity,
	}
}

// MakeMove returns a new board after move is applied, or nil if such move is illegal.
func (board *Board) MakeMove(move *Move) *Board {
	if board.isLegalMove(move.a, move.b) {
		return board.makeMoveWithoutCheck(move)
	}
	return nil
}

// GameStatus returns the game status: 1 (black wins), -1 (white wins), 0 (inconclusive).
func (board *Board) GameStatus() int8 {
	simpleStatus := computeSquareStatus(board.bigSquaresStatus[:], 0)
	if simpleStatus != 2 {
		return simpleStatus
	}
	if board.blackBigSquaresCounter > board.whiteBigSquaresCounter {
		return 1
	}
	return -1
}

func (board *Board) printTileContent(index int) {
	content := board.tiles[index]
	if content == 1 {
		print("b")
	} else if content == -1 {
		print("w")
	} else if content == 0 {
		print("0")
	} else {
		panic("Bad board!")
	}
}

func (board *Board) print() {
	if board.playerIdentity == 1 {
		println("Current Player: Black")
	} else {
		println("Current Player: White")
	}
	println("Printing the board:")
	println("-----------------")
	for row := 0; row < 3; row++ {
		for innerRow := 0; innerRow < 3; innerRow++ {
			board.printTileContent(row*27 + innerRow*3)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 1)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 2)
			print("|")
			board.printTileContent(row*27 + innerRow*3 + 9)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 10)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 11)
			print("|")
			board.printTileContent(row*27 + innerRow*3 + 18)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 19)
			print(" ")
			board.printTileContent(row*27 + innerRow*3 + 20)
			print("\n")
		}
		if row != 2 {
			println("-----*-----*-----")
		}
	}
	println("-----------------")
}

// RunAGameBetweenTwoAIs runs a game between two AI with a specified [aiThinkingTime] in ms.
func RunAGameBetweenTwoAIs(aiThinkingTime float64) {
	board := makeBoard()
	moveCounter := 1
	var status int8
	for status == 0 {
		board.print()
		response := selectMove(board, aiThinkingTime)
		responseMove := response.Move
		board = board.makeMoveWithoutCheck(&Move{a: responseMove[0], b: responseMove[1]})
		fmt.Printf("Move %d finished.\n", moveCounter)
		status = board.GameStatus()
		var player string
		if board.playerIdentity == -1 {
			player = "Black"
		} else {
			player = "White"
		}
		fmt.Printf("Winning probability for %s is %d%%.\n", player, response.WinningProbability)
		moveCounter++
	}
	board.print()
	if status == 1 {
		println("Black wins!")
	} else {
		println("White wins!")
	}
}

// RespondToClient returns the MctsResponse from to a client board.
// It assumes that the client information is legal.
func RespondToClient(clientBoard *BoardData) *MctsResponse {
	println("Handling client request...")
	var player string
	if clientBoard.PlayerIdentity == 1 {
		player = "Black"
	} else {
		player = "White"
	}
	fmt.Printf("Player is %s.", player)
	squareID := clientBoard.BigSquareToPick
	if squareID == -1 {
		println("We can move anywhere.")
	} else {
		fmt.Printf("We have to move in big square %d.", squareID)
	}
	response := selectMove(makeBoardFromData(clientBoard), 1.5)
	println("Decided AI response.")
	return response
}
