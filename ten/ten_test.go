package ten

import (
	"math/rand"
	"testing"
)

// Test on legal moves
func TestLegality(t *testing.T) {
	for i := 0; i < 1<<14; i++ {
		board := MakeBoard()
		for board.GameStatus() == 0 {
			legalMoves := board.AllLegalMovesForAI()
			length := len(legalMoves)
			if length == 0 {
				t.Fatal("Bad!")
			}
			move := legalMoves[rand.Intn(length)]
			successful := board.MakeMove(&move)
			if !successful {
				t.Fatal("Illegal Move!")
			}
			// board.Print()
		}
	}
}

// Test on performance
func BenchmarkPerformance(_ *testing.B) {
	RunAGameBetweenTwoAIs(50, false)
}
