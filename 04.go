package main

import (
	"regexp"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	re := regexp.MustCompile(`^(?P<A>\d+)-(?P<B>\d+),(?P<C>\d+)-(?P<D>\d+)$`)
	var pairs [][2][2]int
	scanner := ctx.scanner()
	for scanner.scan() {
		var x struct {
			A, B, C, D int
		}
		hasty.MustParse([]byte(scanner.text()), &x, re)
		pair := [2][2]int{{x.A, x.B}, {x.C, x.D}}
		pairs = append(pairs, pair)
	}
	ctx.reportLoad()

	var part1 int
	for _, pair := range pairs {
		if contains(pair[0], pair[1]) || contains(pair[1], pair[0]) {
			part1++
		}
	}
	ctx.reportPart1(part1)

	var part2 int
	for _, pair := range pairs {
		if overlaps(pair[0], pair[1]) {
			part2++
		}
	}
	ctx.reportPart2(part2)
}

func contains(r0, r1 [2]int) bool {
	return r0[0] <= r1[0] && r0[1] >= r1[1]
}

func overlaps(r0, r1 [2]int) bool {
	if r0[0] > r1[0] {
		r0, r1 = r1, r0
	}
	return r0[1] >= r1[0]
}
