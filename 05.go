package main

import (
	"regexp"
	"strings"

	"github.com/cespare/hasty"
	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	var stacks cargoStacks
	var instructions []cargoInstruction
	re := regexp.MustCompile(`^move (?P<N>\d+) from (?P<From>\d+) to (?P<To>\d+)$`)
	scanner := ctx.scanner()
	for scanner.scan() {
		line := scanner.text()
		switch {
		case strings.Contains(line, "["):
			stacks.addLine(line)
		case re.MatchString(line):
			var inst cargoInstruction
			hasty.MustParse([]byte(line), &inst, re)
			instructions = append(instructions, inst)
		}
	}
	stacks.reverse()
	ctx.reportLoad()

	stacks1 := stacks.clone()
	for _, inst := range instructions {
		stacks1.apply1(inst)
	}
	ctx.reportPart1(stacks1.message())

	stacks2 := stacks.clone()
	for _, inst := range instructions {
		stacks2.apply2(inst)
	}
	ctx.reportPart2(stacks2.message())
}

type cargoStacks [][]byte

func (s *cargoStacks) addLine(line string) {
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch c {
		case ' ', '[', ']':
			continue
		}
		idx := i / 4
		for len(*s)-1 < idx {
			*s = append(*s, nil)
		}
		(*s)[idx] = append((*s)[idx], c)
	}
}

func (s cargoStacks) reverse() {
	for _, stk := range s {
		SliceReverse(stk)
	}
}

func (s cargoStacks) clone() cargoStacks {
	return SliceMap(s, slices.Clone[[]byte])
}

func (s cargoStacks) apply1(inst cargoInstruction) {
	for i := 0; i < inst.N; i++ {
		x := SlicePop(&s[inst.From-1])
		SlicePush(&s[inst.To-1], x)
	}
}

func (s cargoStacks) apply2(inst cargoInstruction) {
	from := &s[inst.From-1]
	to := &s[inst.To-1]
	*to = append(*to, (*from)[len(*from)-inst.N:]...)
	*from = (*from)[:len(*from)-inst.N]
}

func (s cargoStacks) message() string {
	return string(SliceMap(s, func(stk []byte) byte { return stk[len(stk)-1] }))
}

type cargoInstruction struct {
	N    int
	From int
	To   int
}
