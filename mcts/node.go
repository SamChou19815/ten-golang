package mcts

import "math"

/*
A node in the simulation tree.

It is initialized by an optional parent [parent] (which is used to track back to update winning probability),
a move [move] of the user or AI, the board [board] bind on the node.
*/
type node struct {
	parent         *node
	move           *[]int
	children       *[]*node
	winningProbNum int
	winningProbDen int
	board          *Board
}

/*
[makeRootNode] creates a node with the given [board].

It construct a node without a parent and without a move, with only a starting board.
This node can only be root node.
*/
func makeRootNode(board *Board) *node {
	return &node{board: board}
}

// [makeNode] creates a normal node with the given [parent], [move] and [board].
func makeNode(parent *node, move *[]int, board *Board) *node {
	return &node{parent: parent, move: move, board: board}
}

func (n *node) winningProbability() float64 {
	return float64(n.winningProbNum) / float64(n.winningProbDen)
}

// Obtain winning probability in percentage.
func (n *node) winningProbabilityInPercentage() int {
	return int(n.winningProbability() * 100)
}

/**
Plus one for winning probability denominator and plus the [winValue] for
the numerator. This method does this iteratively until reaching the root.
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

Get upper confidence bound in MCTS, which needs a [isPlayer] parameter to tell whether to calculate
in favor or against the player.

Requires: the node is not the root.
*/
func (n *node) getUpperConfidenceBound(isPlayer bool) float64 {
	if n.parent == nil {
		panic("Cannot be called on root element!")
	}
	den := n.parent.winningProbDen
	var lnt = math.Log(float64(den))
	const c = 1.0
	var winningProb float64
	if isPlayer {
		winningProb = n.winningProbability()
	} else {
		winningProb = 1 - n.winningProbability()
	}
	return winningProb + math.Sqrt(2*lnt/float64(den))*c
}

// Remove the board from the node to allow it to be garbage collected.
func (n *node) dereferenceBoard() {
	n.board = nil
}

// Print itself at [level].
func (n *node) print(level int) {
	printSpace := func() {
		for i := 0; i < level; i++ {
			print(" ")
		}
	}
	printSpace()
	println("Winning Probability:", n.winningProbNum, n.winningProbDen)
	if n.children != nil {
		printSpace()
		println("Children:")
		for _, node := range *n.children {
			node.print(level + 1)
		}
	}
}
