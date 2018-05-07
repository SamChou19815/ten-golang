/*
MCTS stands for Monte Carlo tree search.
*/
package mcts

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

/**
It's constructed with a initial [board] and the [timeLimit] in milliseconds.
*/
type mcts struct {
	timeLimit             int64
	currentPlayerIdentity int
	tree                  *node
}

// Select and return a node starting from parent, according to selection rule in MCTS.
func (m *mcts) selection() *node {
	var root *node = m.tree
	isPlayer := true
	for {
		// Find optimal move and loop down.
		if root.children == nil {
			return root
		}
		var children []*node = *root.children
		length := len(children)
		if length == 0 {
			return root
		}
		var n *node = children[0]
		max := n.getUpperConfidenceBound(isPlayer)
		for i := 1; i < length; i++ {
			var node *node = children[i]
			ucb := node.getUpperConfidenceBound(isPlayer)
			if ucb > max {
				max = ucb
				n = node
			}
		}
		isPlayer = !isPlayer // switch player identity
		root = n
	}
}

// Perform simulation for a specific node [nodeToBeSimulated] and gives back a win value between 0 and 1.
func (m *mcts) simulation(nodeToBeSimulated *node) int {
	var boardBeforeSimulation Board = *nodeToBeSimulated.board
	b1 := boardBeforeSimulation.Copy()
	status := b1.GameStatus()
	for status == 0 {
		moves := b1.AllLegalMovesForAI()
		move := moves[rand.Intn(len(moves))]
		b1.MakeMoveWithoutCheck(&move)
		status = b1.GameStatus()
	}
	if status == m.currentPlayerIdentity {
		return 1
	} else if status == -1*m.currentPlayerIdentity {
		return 0
	} else {
		panic("Bad Status")
	}
}

// Small struct to give goroutine job
type job struct {
	board        *Board
	selectedNode *node
	move         *[]int
}

// Small struct to give goroutine result
type result struct {
	node     *node
	winValue int
}

// Compute the value of a new child node.
func (m *mcts) computeNewChildNodeWorker(jobs <-chan *job, ch chan *result, wg *sync.WaitGroup) {
	for job := range jobs {
		b := *job.board
		b1 := b.Copy()
		move := job.move
		b1.MakeMoveWithoutCheck(move)
		n := makeNode(job.selectedNode, move, &b1)
		winValue := m.simulation(n)
		// Push result down to channel
		ch <- &result{node: n, winValue: winValue}
		// We are done.
		wg.Done()
	}
}

// A method that connected all parts of of MCTS to build an evaluation tree.
func (m *mcts) think() {
	tStart := time.Now()
	cores := runtime.NumCPU() / 2
	simulationCounter := 0
	for time.Now().Sub(tStart).Nanoseconds() < m.timeLimit {
		var selectedNode *node = m.selection()
		b := *selectedNode.board
		// Expansion: Get all legal moves from a current board
		allLegalMoves := b.AllLegalMovesForAI()
		length := len(allLegalMoves)
		if length > 0 {
			// Help GC
			selectedNode.dereferenceBoard()
		}
		simulationCounter += length // Update statistics
		newChildren := make([]*node, length)
		jobs := make(chan *job, length)
		ch := make(chan *result, length)
		wg := sync.WaitGroup{}
		wg.Add(length)
		// Use worker pool.
		for w := 0; w < cores; w++ {
			// Setup workers
			go m.computeNewChildNodeWorker(jobs, ch, &wg)
		}
		for i := 0; i < length; i++ {
			move := allLegalMoves[i]
			job := job{
				board:        &b,
				selectedNode: selectedNode,
				move:         &move,
			}
			// Add jobs to workers.
			jobs <- &job
		}
		close(jobs)
		wg.Wait()
		close(ch)
		i := 0
		for res := range ch {
			node := res.node
			node.winningStatisticsPlusOne(res.winValue)
			newChildren[i] = node
			i++
		}
		selectedNode.children = &newChildren
	}
	println("Number of simulations: " + strconv.Itoa(simulationCounter))
}

/*
Give the final move chosen by AI with the format:
(...decided move, winning probability percentage by that move).
*/
func SelectMove(board *Board, timeLimit int64) []int {
	mcts := &mcts{
		timeLimit:             timeLimit * 1000000,
		currentPlayerIdentity: (*board).CurrentPlayerIdentity(),
		tree: makeRootNode(board),
	}
	mcts.think()
	children := *mcts.tree.children
	length := len(children)
	if length == 0 {
		panic("Impossible Length!")
	}
	nodeChosen := children[0]
	maxWinningProbability := nodeChosen.winningProbability()
	for i := 1; i < length; i++ {
		n := children[i]
		value := n.winningProbability()
		if value > maxWinningProbability {
			maxWinningProbability = value
			nodeChosen = n
		}
	}
	move := nodeChosen.move
	winningProbPercentage := nodeChosen.winningProbabilityInPercentage()
	// Fill in information
	return append(*move, winningProbPercentage)
}
