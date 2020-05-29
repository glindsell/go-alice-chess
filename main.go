package main

import "fmt"

const (
	startFEN      = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	startFENEmpty = "8/8/8/8/8/8/8/8 w KQkq - 0 1"
)

type Game struct {
	Pos1 *Position
	Pos2 *Position
}

func main() {
	pos1, _ := decodeFEN(startFEN)
	pos2, _ := decodeFEN(startFENEmpty)
	fmt.Printf("%v", pos1.Board().Draw())
	fmt.Printf("%v", pos2.Board().Draw())
	m := &Move{
		s1:    Square(8),
		s2:    Square(16),
		promo: NoPieceType,
		tags:  0,
	}
	g := &Game{
		Pos1: pos1,
		Pos2: pos2,
	}

	g = g.Update(m)
	fmt.Printf("%v", g.Pos1.Board().Draw())
	fmt.Printf("%v", g.Pos2.Board().Draw())
}

func (g *Game) Update(m *Move) *Game {
	b1 := g.Pos1.board
	b2 := g.Pos2.board

	s1BB := bbForSquare(m.s1)
	s2BB := bbForSquare(m.s2)
	p1 := b1.Piece(m.s1)

	// move s1 piece to s2 board1
	bb := b1.bbForPiece(p1)
	// remove what was at s2
	b1.setBBForPiece(p1, bb & ^s1BB)
	// move what was at s1 to s2
	bb = b2.bbForPiece(p1)
	b2.setBBForPiece(p1, (bb | s2BB))
	//}
	g.Pos1.board = b1
	g.Pos2.board = b2
	return &Game{
		Pos1: g.Pos1,
		Pos2: g.Pos2,
	}
}
