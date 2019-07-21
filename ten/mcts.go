package ten

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var cores = runtime.NumCPU() / 2

// selection returns a node starting from parent, according to selection rule in MCTS.
func selection(root *node) *node {
	isPlayer := true
	for {
		// Find optimal move and loop down.
		var children = root.children
		length := len(children)
		if length == 0 {
			return root
		}
		var selectedNode = children[0]
		max := selectedNode.getUpperConfidenceBound(isPlayer)
		for i := 1; i < length; i++ {
			var node = children[i]
			ucb := node.getUpperConfidenceBound(isPlayer)
			if ucb > max {
				max = ucb
				selectedNode = node
			}
		}
		isPlayer = !isPlayer // switch player identity
		root = selectedNode
	}
}

// simulation runs a simulation for a specific board and gives back a win value between 0 and 1.
func simulation(playerIdentity int8, board *Board) int {
	currentBoard := board
	status := currentBoard.GameStatus()
	for status == 0 {
		moves := currentBoard.allLegalMovesForAI()
		move := moves[rand.Intn(len(moves))]
		currentBoard = currentBoard.makeMoveWithoutCheck(move)
		status = currentBoard.GameStatus()
	}
	if status == playerIdentity {
		return 1
	}
	return 0
}

/*
 * Small struct to give goroutine job
 */
type job struct {
	board          *Board
	selectedNode   *node
	move           *Move
	playerIdentity int8
}

func computeNewChildNodeWorker(jobs <-chan *job, ch chan *node, wg *sync.WaitGroup) {
	for job := range jobs {
		board := job.board
		move := job.move
		newBoard := board.makeMoveWithoutCheck(move)
		nodeAfterSimulation := &node{
			parent:                 job.selectedNode,
			move:                   move,
			board:                  newBoard,
			children:               emptyNodeList,
			winningProbNumerator:   simulation(job.playerIdentity, newBoard),
			winningProbDenominator: 1,
		}
		// Push result down to channel
		ch <- nodeAfterSimulation
		// We are done.
		wg.Done()
	}
}

/*
 * A method that connected all parts of of MCTS to build an evaluation tree.
 */
func think(root *node, playerIdentity int8, timeLimit float64) int {
	tStart := time.Now()
	simulationCounter := 0
	for time.Now().Sub(tStart).Seconds() < timeLimit {
		var selectedNode = selection(root)
		b := *selectedNode.board
		// Expansion: Get all legal moves from a current board
		allLegalMoves := b.allLegalMovesForAI()
		length := len(allLegalMoves)
		if length == 0 {
			if playerIdentity == b.GameStatus() {
				selectedNode.updateWinningProbability(1, 1)
			} else {
				selectedNode.updateWinningProbability(0, 1)
			}
			simulationCounter++
		} else {
			// Help GC
			selectedNode.board = nil
			// Parallelize the simulation
			newChildren := make([]*node, length)
			jobs := make(chan *job, length)
			ch := make(chan *node, length)
			wg := sync.WaitGroup{}
			wg.Add(length)
			// Use worker pool.
			for w := 0; w < cores; w++ {
				// Setup workers
				go computeNewChildNodeWorker(jobs, ch, &wg)
			}
			// Add jobs to workers.
			for i := 0; i < length; i++ {
				move := allLegalMoves[i]
				job := job{
					board:          &b,
					selectedNode:   selectedNode,
					move:           move,
					playerIdentity: playerIdentity,
				}
				jobs <- &job
			}
			close(jobs)
			wg.Wait()
			close(ch)
			// Collect simulation results
			i := 0
			winCount := 0
			for node := range ch {
				newChildren[i] = node
				winCount += node.winningProbNumerator
				i++
			}
			selectedNode.children = newChildren
			selectedNode.updateWinningProbability(winCount, length)
			simulationCounter += length // Update statistics
		}
	}
	println("Number of simulations: " + strconv.Itoa(simulationCounter))
	return simulationCounter
}

// MctsResponse is the response from the MCTS AI.
type MctsResponse struct {
	Move               []int `json:"move"`
	WinningProbability int   `json:"winningPercentage"`
	SimulationCounter  int   `json:"simulationCounter"`
}

/*
 * Give the final move chosen by AI with the format:
 * (...decided move, winning probability percentage by that move, simulation counter).
 */
func selectMove(board *Board, timeLimit float64) *MctsResponse {
	root := &node{board: board, children: emptyNodeList}
	simulationCounter := think(root, board.playerIdentity, timeLimit)
	children := root.children
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
	return &MctsResponse{
		Move:               []int{move.a, move.b},
		WinningProbability: maxWinningProbability,
		SimulationCounter:  simulationCounter,
	}
}
