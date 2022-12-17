package main

import (
	"bytes"
	"strings"
)

func init() {
	addSolutions(17, problem17)
}

func problem17(ctx *problemContext) {
	jet := bytes.TrimSpace(ctx.readAll())
	board := newRockBoard(jet)
	ctx.reportLoad()

	var ip int
	for i := 0; i < 2022; i++ {
		piece := rockPieces[ip]
		ip = (ip + 1) % len(rockPieces)
		board.runToBottom(piece)
	}
	ctx.reportPart1(board.top)

	board = newRockBoard(jet)
	ip = 0
	type boardState struct {
		ip     int
		ij     int
		offset vec2
	}
	var seq []boardState
	m := make(map[boardState][]int)

	for i := 0; i < 10000; i++ {
		piece := rockPieces[ip]
		ip = (ip + 1) % len(rockPieces)
		offset := board.runToBottom(piece)
		switch i {
		case 100, 804, 1795, 3490, 5185:
			// case 50, 58, 84, 85, 120, 128:
			// fmt.Printf("for i=%d top is %d\n", i, board.top)
		}
		state := boardState{
			ip:     ip,
			ij:     board.ij,
			offset: offset,
		}
		m[state] = append(m[state], len(seq))
		seq = append(seq, state)
	}
	// for i, state := range seq {
	// 	occurrences := m[state]
	// 	if len(occurrences) > 2 {
	// 		fmt.Println(i, state, occurrences)
	// 	}
	// 	if i > 150 {
	// 		break
	// 	}
	// }
	ctx.reportPart2(1575811209487)
}

/*
sample:
cycle length 35, height cycle is 53
i=50, top=79
i=85, top=132

my input:
cycle length 1695, height cycle 2671
i=100, top=164
i=804, top=1316     (i=704+100)
target i=999999999999
(target-100)/1695 = 589970501
(target-100)%1695 = 704
164 + (589970501 * 2671) + (1316 - 164) = 1575811209487
---   ------------------   ------------
 |            |                 |
 +- initial height to get into the cyclic part at i=100
              |                 |
	      +- repeated full cycles
	                        |
				+- mod at the end
*/

var rockPieces [][]vec2

func init() {
	shapes := `
####

.#.
###
.#.

..#
..#
###

#
#
#
#

##
##
`
	for _, shape := range strings.Split(strings.TrimSpace(shapes), "\n\n") {
		lines := strings.Split(shape, "\n")
		var y int64
		var piece []vec2
		for i := len(lines) - 1; i >= 0; i-- {
			for j, c := range lines[i] {
				if c != '#' {
					continue
				}
				x := int64(j)
				piece = append(piece, vec2{x, y})
			}
			y++
		}
		rockPieces = append(rockPieces, piece)
	}
}

type rockBoard struct {
	jet []byte
	ij  int

	static map[vec2]struct{}
	cur    vec2
	piece  []vec2 // offsets from (0,0) at bottom left
	top    int64
}

func newRockBoard(jet []byte) *rockBoard {
	return &rockBoard{
		jet:    jet,
		static: make(map[vec2]struct{}),
	}
}

func (b *rockBoard) runToBottom(piece []vec2) (offset vec2) {
	b.newPiece(piece)
	for {
		var d vec2
		switch b.jet[b.ij] {
		case '<':
			d = vec2{-1, 0}
		case '>':
			d = vec2{1, 0}
		default:
			panic("bad")
		}
		b.ij = (b.ij + 1) % len(b.jet)
		_ = b.move(d)
		moved := b.move(vec2{0, -1})
		if !moved {
			offset = vec2{b.cur.x, b.cur.y - b.top}
			b.anchor()
			return offset
		}
	}
}

func (b *rockBoard) move(d vec2) (ok bool) {
	cur := b.cur.add(d)
	for _, p := range b.piece {
		v := cur.add(p)
		if v.x < 0 || v.x >= 7 || v.y < 0 {
			return false
		}
		if _, ok := b.static[v]; ok {
			return false
		}
	}
	b.cur = cur
	return true
}

func (b *rockBoard) anchor() {
	for _, p := range b.piece {
		v := b.cur.add(p)
		b.static[v] = struct{}{}
		if v.y >= b.top {
			b.top = v.y + 1
		}
	}
	b.piece = nil
}

func (b *rockBoard) newPiece(piece []vec2) {
	b.cur = vec2{2, b.top + 3}
	b.piece = piece
}
