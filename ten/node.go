package ten

import (
	"math"
)

/*
 * A node in the simulation tree.
 *
 * It is initialized by an optional parent [parent] (which is used to track back to update winning probability),
 * a move [move] of the user or AI, the board [board] bind on the node.
*/
type node struct {
	parent         *node
	move           *move
	children       *[]*node
	winningProbNum int
	winningProbDen int
	board          *Board
}

/*
 * [makeRootNode] creates a node with the given [board].
 *
 * It construct a node without a parent and without a move, with only a starting board.
 * This node can only be root node.
 */
func makeRootNode(board *Board) *node {
	return &node{board: board}
}

/*
 * [makeNode] creates a normal node with the given [parent], [move] and [board].
 */
func makeNode(parent *node, move *move, board *Board) *node {
	return &node{parent: parent, move: move, board: board}
}

/*
 * Obtain winning probability of the current node.
 */
func (n *node) winningProbability() int {
	return 100 * n.winningProbNum / n.winningProbDen
}

/**
 * Plus one for winning probability denominator and plus the [winValue] for
 * the numerator. This method does this iteratively until reaching the root.
 */
func (n *node) winningStatisticsPlusOne(winValue int) {
	var node = n
	for node != nil {
		node.winningProbNum += winValue
		node.winningProbDen += 1
		node = node.parent
	}
}

/**
 *
 * Get upper confidence bound in MCTS, which needs a [isPlayer] parameter to tell whether to calculate
 * in favor or against the player.
 *
 * Requires: the node is not the root.
 */
func (n *node) getUpperConfidenceBound(isPlayer bool) int {
	den := n.parent.winningProbDen
	var lnt = math.Log(float64(den))
	var winningProb int
	if isPlayer {
		winningProb = n.winningProbability()
	} else {
		winningProb = 100 - n.winningProbability()
	}
	var secondPart = math.Sqrt(2*lnt/float64(den)) * 100
	return winningProb + int(secondPart)
}

/*
 * Remove the board from the node to allow it to be garbage collected.
 */
func (n *node) dereferenceBoard() {
	n.board = nil
}
