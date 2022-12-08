package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(8, problem8)
}

var nesw = []vec2{
	{0, 1},
	{-1, 0},
	{0, -1},
	{1, 0},
}

func problem8(ctx *problemContext) {
	var forest grid[int]
	scanner := ctx.scanner()
	for scanner.scan() {
		row := SliceMap(
			[]byte(scanner.text()),
			func(c byte) int { return int(c) - '0' },
		)
		forest.addRow(row)
	}
	ctx.reportLoad()

	var visible set.Set[vec2]
	var highest int
	for i := range nesw {
		major := nesw[i]
		minor := nesw[(i+3)%4]
		var v0 vec2
		if major.x < 0 || minor.x < 0 {
			v0.x = forest.cols - 1
		}
		if major.y < 0 || minor.y < 0 {
			v0.y = forest.rows - 1
		}
		for v1 := v0; forest.contains(v1); v1 = v1.add(minor) {
			highest = -1
			for v2 := v1; forest.contains(v2); v2 = v2.add(major) {
				cur := forest.at(v2)
				if cur > highest {
					visible.Add(v2)
					highest = cur
				}
			}
		}
	}
	ctx.reportPart1(visible.Len())

	var best int64
	forest.forEach(func(v vec2, h int) {
		score := int64(1)
		for _, d := range nesw {
			score *= viewingDistance(&forest, v, d, h)
		}
		if score > best {
			best = score
		}
	})
	ctx.reportPart2(best)
}

func viewingDistance(forest *grid[int], v, d vec2, h int) int64 {
	var score int64
	for {
		v = v.add(d)
		if !forest.contains(v) {
			return score
		}
		score++
		if forest.at(v) >= h {
			return score
		}
	}
}
