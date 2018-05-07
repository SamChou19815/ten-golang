package mcts

// This interface represents any game board that supports MCTS AI framework.
type Board interface {
	// Obtain identity of current player.
	CurrentPlayerIdentity() int

	// Obtain a deep copy of the board to allow different simulations on the same board without interference.
	Copy() Board

	/*
		Obtain the [gameStatus] on current board.
		This happens immediately after a player makes a move, before switching identity.
		The status must be 1, -1, or 0 (inconclusive).
	*/
	GameStatus() int

	/*
		Obtain a list of all legal moves for AI.
		DO NOT confuse: a legal move for human is not necessarily one for AI
		because AI needs less moves to save time for computation.
	*/
	AllLegalMovesForAI() [][]int

	/*
		Make a [move] without any check, which can accelerate AI simulation.
		It should also switch the identity of the current player.
		The identity must be 1 or -1, so that checking game status can determine who wins.

		Requires: [move] is a valid int array representation of a move in the game.
	*/
	MakeMoveWithoutCheck(*[]int)

	// Print the board.
	Print()
}
