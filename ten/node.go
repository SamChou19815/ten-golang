package ten

import (
	"math"
)

/*
 * A node in the simulation tree.
 *
 * It is initialized by an optional parent (used to track back to update winning probability),
 * a move of the user or AI, the board bind to the node.
 */
type node struct {
	parent                 *node
	move                   *Move
	children               []*node
	winningProbNumerator   int
	winningProbDenominator int
	board                  *Board
}

// emptyNodeList is used for dummy children list.
var emptyNodeList = make([]*node, 0)

// winningProbability returns winning probability of the current node.
func (node *node) winningProbability() int {
	return 100 * node.winningProbNumerator / node.winningProbDenominator
}

// updateWinningProbability recursively update the winning probability of the node from the
// given node to the root.
func (node *node) updateWinningProbability(winCount int, totalCount int) {
	currentNode := node
	for currentNode != nil {
		currentNode.winningProbNumerator += winCount
		currentNode.winningProbDenominator += totalCount
		currentNode = currentNode.parent
	}
}

// getUpperConfidenceBound returns upper confidence bound in MCTS, which needs a isPlayer parameter
// to tell whether to calculate in favor or against the player.
// The node must not be the root.
func (node *node) getUpperConfidenceBound(isPlayer bool) int {
	denominator := node.parent.winningProbDenominator
	floatDenominator := float64(denominator)
	var lnt = math.Log(floatDenominator)
	var winningProb int
	if isPlayer {
		winningProb = node.winningProbability()
	} else {
		winningProb = 100 - node.winningProbability()
	}
	var secondPart = math.Sqrt(2*lnt/floatDenominator) * 100
	return winningProb + int(secondPart)
}
