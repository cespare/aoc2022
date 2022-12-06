package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	line := ctx.readAll()
	ctx.reportLoad()
	find := func(n int) int {
		for i := range line {
			if set.Of(line[i:i+n]...).Len() == n {
				return i + n
			}
		}
		panic("fail")
	}
	ctx.reportPart1(find(4))
	ctx.reportPart2(find(14))
}
