package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	startFEN      = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	startFENEmpty = "8/8/8/8/8/8/8/8 w KQkq - 0 1"
	startFENA     = "1r6/8/8/8/8/8/8/Q7 w KQkq - 0 1"
	startFENB     = "8/8/8/8/8/8/N7/8 w KQkq - 0 1"
)

// TODO: move turn to game struct
type Game struct {
	Pos1 *Position
	Pos2 *Position
}

func (g *Game) copy() *Game {
	return &Game{
		Pos1: g.Pos1.copy(),
		Pos2: g.Pos2.copy(),
	}
}

func main() {
	pos1, _ := decodeFEN(startFEN)
	pos2, _ := decodeFEN(startFENEmpty)
	g := &Game{
		Pos1: pos1,
		Pos2: pos2,
	}
	fmt.Printf("Move Count: %v\n", 0)
	fmt.Printf("%v", g.Pos1.Board().Draw())
	fmt.Printf("%v\n------\n", g.Pos2.Board().Draw())
	/*m := &Move{
		s1:    Square(8),
		s2:    Square(16),
		promo: NoPieceType,
		tags:  0,
	}*/
	for i := 0; i < 1000; i++ {
		fmt.Printf("Move Count: %v\n", g.Pos1.moveCount)
		fmt.Printf("Turn: %v\n", g.Pos1.Turn().Name())
		g.MoveRand()
		fmt.Printf("\n --- Board A ---\n%v\n", g.Pos1.Board().Draw())
		fmt.Printf(" --- Board B ---\n%v\n *--------------*\n", g.Pos2.Board().Draw())
	}
}

func (g *Game) MoveRand() {
	standardMoves := []*Move{}
	standardMovesA := g.StandardMovesA(false)
	fmt.Printf("Moves A: %v\n", standardMovesA)
	standardMovesB := g.StandardMovesB(false)
	fmt.Printf("Moves B: %v\n", standardMovesB)
	standardMoves = append(standardMoves, standardMovesA...)
	standardMoves = append(standardMoves, standardMovesB...)
	board := ""
	if len(standardMoves) > 0 {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(standardMoves))
		if r < len(standardMovesA) {
			g.UpdateA(standardMoves[r])
			board = "A"
		} else {
			g.UpdateB(standardMoves[r])
			board = "B"
		}
		fmt.Printf("Chosen: %v\nBoard: %v\n", standardMoves[r], board)
		g.Pos1.turn = g.Pos1.turn.Other()
		g.Pos2.turn = g.Pos2.turn.Other()
		g.Pos1.moveCount++
		g.Pos2.moveCount++
	} else {
		panic("No possible moves")
	}
}

func (g *Game) UpdateB(m *Move) {
	s1BB := bbForSquare(m.s1)
	s2BB := bbForSquare(m.s2)
	p1 := g.Pos2.board.Piece(m.s1)

	// move s1 piece to s2 board1
	bb := g.Pos2.board.bbForPiece(p1)
	// remove what was at s1
	g.Pos2.board.setBBForPiece(p1, bb & ^s1BB)

	if m.HasTag(Capture) {
		p2 := g.Pos2.board.Piece(m.s2)
		bb := g.Pos2.board.bbForPiece(p2)
		// remove what was at s2
		g.Pos2.board.setBBForPiece(p2, bb & ^s2BB)
	}

	// move what was at s1 to s2
	bb = g.Pos1.board.bbForPiece(p1)
	g.Pos1.board.setBBForPiece(p1, (bb | s2BB))
	g.Pos2.board.calcConvienceBBs(m)
	g.Pos1.board.calcConvienceBBs(m)
}

func (g *Game) UpdateA(m *Move) {
	s1BB := bbForSquare(m.s1)
	s2BB := bbForSquare(m.s2)
	p1 := g.Pos1.board.Piece(m.s1)

	// move s1 piece to s2 board1
	bb := g.Pos1.board.bbForPiece(p1)
	// remove what was at s1
	g.Pos1.board.setBBForPiece(p1, bb & ^s1BB)

	if m.HasTag(Capture) {
		p2 := g.Pos1.board.Piece(m.s2)
		bb := g.Pos1.board.bbForPiece(p2)
		// remove what was at s2
		g.Pos1.board.setBBForPiece(p2, bb & ^s2BB)
	}

	// move what was at s1 to s2
	bb = g.Pos2.board.bbForPiece(p1)
	g.Pos2.board.setBBForPiece(p1, (bb | s2BB))
	g.Pos1.board.calcConvienceBBs(m)
	g.Pos2.board.calcConvienceBBs(m)
}
