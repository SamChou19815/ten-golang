/*
 * MCTS stands for Monte Carlo tree search.
 */
package ten

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
	"runtime"
)

/**
 * It's constructed with a initial [board] and the [timeLimit] in milliseconds.
 */
type mcts struct {
	timeLimit             int64
	currentPlayerIdentity int
	tree                  *node
}

var cores = runtime.NumCPU() / 2

/*
 * Select and return a node starting from parent, according to selection rule in MCTS.
 */
func (m *mcts) selection() *node {
	var root = m.tree
	isPlayer := true
	for {
		// Find optimal move and loop down.
		if root.children == nil {
			return root
		}
		var children = *root.children
		length := len(children)
		if length == 0 {
			return root
		}
		var n = children[0]
		max := n.getUpperConfidenceBound(isPlayer)
		for i := 1; i < length; i++ {
			var node = children[i]
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

/*
 * Perform simulation for a specific node [nodeToBeSimulated] and gives back a win value between 0 and 1.
 */
func (m *mcts) simulation(nodeToBeSimulated *node) int {
	var boardBeforeSimulation = *nodeToBeSimulated.board
	b1 := boardBeforeSimulation.copy()
	status := b1.gameStatus()
	for status == 0 {
		moves := b1.allLegalMovesForAI()
		move := moves[rand.Intn(len(moves))]
		b1.makeMoveWithoutCheck(move)
		status = b1.gameStatus()
	}
	if status == m.currentPlayerIdentity {
		return 1
	} else {
		return 0
	}
}

/*
 * Small struct to give goroutine job
 */
type job struct {
	board        *Board
	selectedNode *node
	move         *move
}

/*
 * Small struct to give goroutine result
 */
type result struct {
	node     *node
	winValue int
}

/*
 * Compute the value of a new child node.
 */
func (m *mcts) computeNewChildNodeWorker(jobs <-chan *job, ch chan *result, wg *sync.WaitGroup) {
	for job := range jobs {
		b := *job.board
		b1 := b.copy()
		move := job.move
		b1.makeMoveWithoutCheck(move)
		n := makeNode(job.selectedNode, move, b1)
		winValue := m.simulation(n)
		// Push result down to channel
		ch <- &result{node: n, winValue: winValue}
		// We are done.
		wg.Done()
	}
}

/*
 * A method that connected all parts of of MCTS to build an evaluation tree.
 */
func (m *mcts) think() {
	tStart := time.Now()
	simulationCounter := 0
	for time.Now().Sub(tStart).Nanoseconds() < m.timeLimit {
		var selectedNode = m.selection()
		b := *selectedNode.board
		// Expansion: Get all legal moves from a current board
		allLegalMoves := b.allLegalMovesForAI()
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
				move:         move,
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
 * Give the final move chosen by AI with the format:
 * (...decided move, winning probability percentage by that move).
 */
func selectMove(board *Board, timeLimit int64) (*move, int) {
	mcts := &mcts{
		timeLimit:             timeLimit * 1000000,
		currentPlayerIdentity: (*board).currentPlayerIdentity,
		tree:                  makeRootNode(board),
	}
	mcts.think()
	children := *mcts.tree.children
	length := len(children)
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
	// Fill in information
	return move, maxWinningProbability
}
