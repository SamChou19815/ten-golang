package ten

import (
	"../mcts"
	"strconv"
)

// [Board] encapsulates the functionality of the board of TEN.
type Board interface {
	mcts.Board

	// Make a move [move] with legality check and tells whether the move is legal/successful.
	MakeMove(*[]int) bool
}

// [tenBoard] is the concrete type of the board game TEN.
type tenBoard struct {
	Board

	/*
		In variable names, a big square refers to a 3*3 square;
		a tile refers to a 1*1 square.
		Each tile is either 1, -1 or 0 (black, white, empty).
	*/
	board [9][9]int
	/*
		Keep track of winning progress on big squares.
		1, -1, 0, 2 mean black wins, white wins, inconclusive, and all occupied.
	*/
	bigSquaresStatus [9]int
	/*
		The current legal big square to pick as next move.
		If it's value is -1, that means the user can place the move everywhere.
		This variable is important for maintaining the current game state.
	*/
	currentBigSquareLegalPosition int
	// The identity of the current player. Must be 1 or -1.
	currentPlayerIdentity int
}

// [MakeBoard] creates an empty board.
func MakeBoard() Board {
	var rawBoard [9][9]int
	var bigSquaresStatus [9]int
	return &tenBoard{
		board:                         rawBoard,
		bigSquaresStatus:              bigSquaresStatus,
		currentBigSquareLegalPosition: -1,
		currentPlayerIdentity:         1,
	}
}

// [makeBoardFromData] creates an board from existing [data].
func makeBoardFromData(data BoardData) *tenBoard {
	rawBoard := data.board
	var bigSquaresStatus [9]int
	board := &tenBoard{
		board:                         rawBoard,
		bigSquaresStatus:              bigSquaresStatus,
		currentBigSquareLegalPosition: data.currentBigSquareLegalPosition,
		currentPlayerIdentity:         data.currentPlayerIdentity,
	}
	for i := 0; i < 9; i++ {
		board.updateBigSquareStatus(i)
	}
	return board
}

// Decode int [i] stored internally in data structure to player name.
func decode(i int) string {
	switch i {
	case 0:
		return "0"
	case 1:
		return "b"
	case -1:
		return "w"
	default:
		panic("Bad Data in Board!")
	}
}

func (b *tenBoard) Print() {
	var playerString string
	if b.currentPlayerIdentity == 1 {
		playerString = "Black"
	} else {
		playerString = "White"
	}
	println("Current Player: " + playerString)
	println("Printing the board:")
	println("-----------------")
	for row := 0; row < 3; row++ {
		for innerRow := 0; innerRow < 3; innerRow++ {
			print(decode(b.board[row*3][innerRow*3]) + " " +
				decode(b.board[row*3][innerRow*3+1]) + " " +
				decode(b.board[row*3][innerRow*3+2]) + "|")
			print(decode(b.board[row*3+1][innerRow*3]) + " " +
				decode(b.board[row*3+1][innerRow*3+1]) + " " +
				decode(b.board[row*3+1][innerRow*3+2]) + "|")
			print(decode(b.board[row*3+2][innerRow*3]) + " " +
				decode(b.board[row*3+2][innerRow*3+1]) + " " +
				decode(b.board[row*3+2][innerRow*3+2]))
			println()
		}
		if row != 2 {
			println("- - -*- - -*- - -")
		}
	}
	println("-----------------")
}

func (b *tenBoard) CurrentPlayerIdentity() int {
	return b.currentPlayerIdentity
}

func (b *tenBoard) Copy() mcts.Board {
	var rawBoard [9][9]int
	var bigSquaresStatus [9]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			rawBoard[i][j] = b.board[i][j]
		}
		bigSquaresStatus[i] = b.bigSquaresStatus[i]
	}
	return &tenBoard{
		board:                         rawBoard,
		bigSquaresStatus:              bigSquaresStatus,
		currentBigSquareLegalPosition: b.currentBigSquareLegalPosition,
		currentPlayerIdentity:         b.currentPlayerIdentity,
	}
}

func (b *tenBoard) AllLegalMovesForAI() [][]int {
	list := make([][]int, 0, 40) // 40 will be enough for most cases.
	if b.currentBigSquareLegalPosition == -1 {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if b.isLegalMove(i, j) {
					move := []int{i, j}
					list = append(list, move)
				}
			}
		}
	} else {
		for j := 0; j < 9; j++ {
			// Can only move in the specified square
			if b.isLegalMove(b.currentBigSquareLegalPosition, j) {
				move := []int{b.currentBigSquareLegalPosition, j}
				list = append(list, move)
			}
		}
	}
	return list
}

func (b *tenBoard) GameStatus() int {
	for i := 0; i < 9; i++ {
		b.updateBigSquareStatus(i)
	}
	bigSquaresStatus := b.bigSquaresStatus
	simpleStatus := getSimpleStatusFromSquare(bigSquaresStatus)
	if simpleStatus == 1 || simpleStatus == -1 {
		return simpleStatus
	}
	for i := 0; i < 9; i++ {
		if bigSquaresStatus[i] == 0 {
			return 0
		}
	}
	blackBigSquareCounter := 0
	whiteBigSquareCounter := 0
	for i := 0; i < 9; i++ {
		status := bigSquaresStatus[i]
		if status == 1 {
			blackBigSquareCounter++
		} else if status == -1 {
			whiteBigSquareCounter++
		}
	}
	if blackBigSquareCounter > whiteBigSquareCounter {
		return 1
	} else {
		return -1
	}
}

// Check whether a move is legal, where move is given by ([a], [c]).
func (b *tenBoard) isLegalMove(a, c int) bool {
	if a < 0 || a > 8 || c < 0 || c > 8 {
		// Out of boundary values
		return false
	}
	if b.currentBigSquareLegalPosition != -1 && b.currentBigSquareLegalPosition != a {
		// in the wrong big square when it cannot have a free move
		return false
	} else {
		// not in the occupied big square and on an empty tile
		return b.bigSquaresStatus[a] == 0 && b.board[a][c] == 0
	}
}

func (b *tenBoard) MakeMoveWithoutCheck(m *[]int) {
	move := *m
	b.board[move[0]][move[1]] = b.currentPlayerIdentity
	b.updateBigSquareStatus(move[0])
	if b.bigSquaresStatus[move[1]] == 0 {
		b.currentBigSquareLegalPosition = move[1]
	} else {
		b.currentBigSquareLegalPosition = -1
	}
	b.switchIdentity()
}

func (b *tenBoard) MakeMove(m *[]int) bool {
	move := *m
	if !b.isLegalMove(move[0], move[1]) {
		return false
	}
	b.MakeMoveWithoutCheck(m)
	return true
}

/*
Perform a naive check on the square [s] about whether the player with identity [id] win the square.
It only checks according to the primitive tic-tac-toe rule.
*/
func playerSimplyWinSquare(s [9]int, id int) bool {
	return s[0] == id && s[1] == id && s[2] == id ||
		s[3] == id && s[4] == id && s[5] == id ||
		s[6] == id && s[7] == id && s[8] == id ||
		s[0] == id && s[3] == id && s[6] == id ||
		s[1] == id && s[4] == id && s[7] == id ||
		s[2] == id && s[5] == id && s[8] == id ||
		s[0] == id && s[4] == id && s[8] == id ||
		s[2] == id && s[4] == id && s[6] == id
}

/*
A function that helps to determine whether a square [square] belongs to black (1) or white (-1).
If there is no direct victory, it will return 0.
*/
func getSimpleStatusFromSquare(square [9]int) int {
	if playerSimplyWinSquare(square, 1) {
		return 1
	} else if playerSimplyWinSquare(square, -1) {
		return -1
	} else {
		return 0
	}
}

// Update the big square status for ONE big square of id [bigSquareID].
func (b *tenBoard) updateBigSquareStatus(bigSquareID int) {
	bigSquare := b.board[bigSquareID]
	bigSquareStatus := getSimpleStatusFromSquare(bigSquare)
	if bigSquareStatus == 1 || bigSquareStatus == -1 {
		b.bigSquaresStatus[bigSquareID] = bigSquareStatus
		return
		// already won by a player
	}
	for i := 0; i < 9; i++ {
		if bigSquare[i] == 0 {
			// there is a space left.
			b.bigSquaresStatus[bigSquareID] = 0
			return
		}
	}
	b.bigSquaresStatus[bigSquareID] = 2 // no space left.
}

// Switches the identity of the current player to complete a move.
func (b *tenBoard) switchIdentity() {
	b.currentPlayerIdentity = -b.currentPlayerIdentity
}

/*
Respond to a [clientMove] represented by a [ClientMove] object and
gives back the formatted [ServerResponse].
*/
func respond(clientMove ClientMove) ServerResponse {
	board := makeBoardFromData(*clientMove.boardBeforeHumanMove)
	board.switchIdentity()
	if !board.MakeMove(clientMove.humanMove) {
		// Stop illegal move from corrupting game data.
		return serverResponseOfIllegalMove
	}
	status := board.GameStatus()
	if status == 1 || status == -1 {
		// Black/White wins before AI move
		return serverResponseOfWinner(status)
	}
	// Let AI think
	var mctsBoard mcts.Board = board
	aiMove := mcts.SelectMove(&mctsBoard, 1500)
	board.MakeMove(&aiMove)
	status = board.GameStatus()
	// A full response.
	return ServerResponse{
		aiMove: [2]int{aiMove[0], aiMove[1]},
		currentBigSquareLegalPosition: board.currentBigSquareLegalPosition,
		status:               status,
		aiWinningProbability: aiMove[2],
	}
}

/*
Run a game between two AI with a specified [aiThinkingTime] in milliseconds.
The user of the function can specify whether to print game status out by [printGameStatus].
*/
func RunAGameBetweenTwoAIs(aiThinkingTime int64, printGameStatus bool) {
	board := MakeBoard()
	moveCounter := 1
	status := 0
	for status == 0 {
		if printGameStatus {
			board.Print()
		}
		var mctsBoard mcts.Board = board
		move := mcts.SelectMove(&mctsBoard, aiThinkingTime)
		board.MakeMoveWithoutCheck(&move)
		status = board.GameStatus()
		if printGameStatus {
			println("Move " + strconv.Itoa(moveCounter) + " finished.")
			var player string
			if moveCounter%2 == 0 {
				player = "White"
			} else {
				player = "Black"
			}
			println("Winning Probability for " + player + " is " + strconv.Itoa(move[2]) + "%.")
		}
		moveCounter++
	}
	if printGameStatus {
		board.Print()
		var player string
		if status == 1 {
			player = "Black"
		} else {
			player = "White"
		}
		println(player + " wins.")
	}
}
