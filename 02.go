package main

import "strings"

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	var strats [][2]byte
	scanner := ctx.scanner()
	for scanner.scan() {
		parts := strings.Split(scanner.text(), " ")
		strat := [2]byte{parts[0][0], parts[1][0]}
		strats = append(strats, strat)
	}
	ctx.reportLoad()

	var s int
	for _, strat := range strats {
		s += int(strat[1] - 'X' + 1)
		us := int(strat[1] - 'X')
		them := int(strat[0] - 'A')
		switch {
		case us == them:
			s += 3
		case (us+1)%3 == them:
		default:
			s += 6
		}
	}
	ctx.reportPart1(s)

	s = 0
	for _, strat := range strats {
		them := int(strat[0] - 'A')
		var us int
		switch strat[1] {
		case 'X':
			us = (them + 2) % 3
		case 'Y':
			us = them
			s += 3
		case 'Z':
			us = (them + 1) % 3
			s += 6
		}
		s += us + 1
	}
	ctx.reportPart2(s)
}
