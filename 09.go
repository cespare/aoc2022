package main

import (
	"fmt"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(9, problem9)
}

func problem9(ctx *problemContext) {
	var insts []ropeInstruction
	scanner := ctx.scanner()
	for scanner.scan() {
		var inst ropeInstruction
		line := scanner.text()
		if _, err := fmt.Sscanf(line, "%c %d", &inst.dir, &inst.n); err != nil {
			panic(err)
		}
		insts = append(insts, inst)
	}
	ctx.reportLoad()

	g := newRopeGrid(2)
	for _, inst := range insts {
		g.do(inst)
	}
	ctx.reportPart1(g.tailVisited.Len())

	g = newRopeGrid(10)
	for _, inst := range insts {
		g.do(inst)
	}
	ctx.reportPart2(g.tailVisited.Len())
}

type ropeGrid struct {
	knots       []vec2
	tailVisited *set.Set[vec2]
}

type ropeInstruction struct {
	dir byte
	n   int
}

func newRopeGrid(numKnots int) *ropeGrid {
	return &ropeGrid{
		knots:       make([]vec2, numKnots),
		tailVisited: set.Of(vec2{0, 0}),
	}
}

func (g *ropeGrid) do(inst ropeInstruction) {
	var d vec2
	switch inst.dir {
	case 'U':
		d = vec2{0, -1}
	case 'D':
		d = vec2{0, 1}
	case 'L':
		d = vec2{-1, 0}
	case 'R':
		d = vec2{1, 0}
	default:
		panic("bad dir")
	}
	for i := 0; i < inst.n; i++ {
		g.knots[0] = g.knots[0].add(d)
		for j := 1; j < len(g.knots); j++ {
			g.knots[j] = followKnot(g.knots[j], g.knots[j-1])
		}
		g.tailVisited.Add(g.knots[len(g.knots)-1])
	}
}

func followKnot(p0, targ vec2) vec2 {
	d := targ.sub(p0)
	if maxim(abs(d.x), abs(d.y)) <= 1 {
		return p0
	}
	adjust := func(n int64) int64 {
		switch n {
		case -2:
			return -1
		case 2:
			return 1
		default:
			return n
		}
	}
	d.x = adjust(d.x)
	d.y = adjust(d.y)
	return p0.add(d)
}
