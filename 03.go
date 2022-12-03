package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	var bags []string
	scanner := ctx.scanner()
	for scanner.scan() {
		bags = append(bags, scanner.text())
	}
	ctx.reportLoad()

	var score int
	for _, bag := range bags {
		items := set.Of([]byte(bag[:len(bag)/2])...)
		for i := len(bag) / 2; i < len(bag); i++ {
			if item := bag[i]; items.Contains(item) {
				score += itemPriority(item)
				break
			}
		}
	}
	ctx.reportPart1(score)

	score = 0
groupLoop:
	for i := 0; i < len(bags); i += 3 {
		group := bags[i : i+3]
		counts := make(map[int]uint)
		for i, bag := range group {
			for j := range bag {
				priority := itemPriority(bag[j])
				counts[priority] |= (1 << i)
				if counts[priority] == 7 {
					score += priority
					continue groupLoop
				}
			}
		}
	}
	ctx.reportPart2(score)
}

func itemPriority(c byte) int {
	if c >= 'a' && c <= 'z' {
		return int(c - 'a' + 1)
	}
	return int(c - 'A' + 27)
}
