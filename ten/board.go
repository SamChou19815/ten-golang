package ten

/*
 * [move] represents a simple move in TEN.
 */
type move struct {
	a int
	b int
}

/*
 * [Board] is the concrete type of the board game TEN.
 */
type Board struct {
	/*
	 * In variable names, a big square refers to a 3*3 square;
     * a tile refers to a 1*1 square.
     * Each tile is either 1, -1 or 0 (black, white, empty).
	 */
	board [9][9]int
	/*
	 * Keep track of winning progress on big squares.
	 * 1, -1, 0, 2 mean black wins, white wins, inconclusive, and all occupied.
	 */
	bigSquaresStatus [9]int
	/*
	 * The current legal big square to pick as next move.
     * If it's value is -1, that means the user can place the move everywhere.
	 * This variable is important for maintaining the current game state.
	 */
	currentBigSquareLegalPosition int
	/**
	 * The identity of the current player. Must be 1 or -1.
	 */
	currentPlayerIdentity int
}

/*
 * [makeBoard] creates an empty board.
 */
func makeBoard() *Board {
	var rawBoard [9][9]int
	var bigSquaresStatus [9]int
	return &Board{
		board:                         rawBoard,
		bigSquaresStatus:              bigSquaresStatus,
		currentBigSquareLegalPosition: -1,
		currentPlayerIdentity:         1,
	}
}

/*
 * Obtain a deep copy of the board to allow different simulations on the same board without interference.
 */
func (b *Board) copy() *Board {
	var rawBoard [9][9]int
	var bigSquaresStatus [9]int
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			rawBoard[i][j] = b.board[i][j]
		}
		bigSquaresStatus[i] = b.bigSquaresStatus[i]
	}
	return &Board{
		board:                         rawBoard,
		bigSquaresStatus:              bigSquaresStatus,
		currentBigSquareLegalPosition: b.currentBigSquareLegalPosition,
		currentPlayerIdentity:         b.currentPlayerIdentity,
	}
}

/**
 * Obtain a list of all legal moves for AI.
 */
func (b *Board) allLegalMovesForAI() []*move {
	list := make([]*move, 0, 40) // 20 will be enough for most cases.
	if b.currentBigSquareLegalPosition != -1 {
		for j := 0; j < 9; j++ {
			// Can only move in the specified square
			if b.isLegalMove(b.currentBigSquareLegalPosition, j) {
				move := move{a: b.currentBigSquareLegalPosition, b: j}
				list = append(list, &move)
			}
		}
	} else {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if b.isLegalMove(i, j) {
					move := move{a: i, b: j}
					list = append(list, &move)
				}
			}
		}
	}
	return list
}

/*
 * Obtain the [gameStatus] on current board.
 * This happens immediately after a player makes a move, before switching identity.
 * The status must be 1, -1, or 0 (inconclusive).
 */
func (b *Board) gameStatus() int {
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

/*
 * Check whether a move is legal, where move is given by ([a], [c]).
 */
func (b *Board) isLegalMove(a, c int) bool {
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

/**
 * Make a [move] without any check, which can accelerate AI simulation.
 * It should also switch the identity of the current player.
 * The identity must be 1 or -1, so that checking game status can determine who wins.
 *
 * Requires: [move] is a valid int array representation of a move in the game.
 */
func (b *Board) makeMoveWithoutCheck(m *move) {
	b.board[m.a][m.b] = b.currentPlayerIdentity
	b.updateBigSquareStatus(m.a)
	if b.bigSquaresStatus[m.b] == 0 {
		b.currentBigSquareLegalPosition = m.b
	} else {
		b.currentBigSquareLegalPosition = -1
	}
	b.currentPlayerIdentity = -b.currentPlayerIdentity // switch identity
}

/*
 * Perform a naive check on the square [s] about whether the player with identity [id] win the square.
 * It only checks according to the primitive tic-tac-toe rule.
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
 * A function that helps to determine whether a square [square] belongs to black (1) or white (-1).
 * If there is no direct victory, it will return 0.
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

/*
 * Update the big square status for ONE big square of id [bigSquareID].
 */
func (b *Board) updateBigSquareStatus(bigSquareID int) {
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

/*
 * Run a game between two AI with a specified [aiThinkingTime] in milliseconds.
 */
func RunAGameBetweenTwoAIs(aiThinkingTime int64) {
	board := makeBoard()
	status := 0
	for status == 0 {
		move, _ := selectMove(board, aiThinkingTime)
		board.makeMoveWithoutCheck(move)
		status = board.gameStatus()
	}
}
