package main

import (
	"regexp"
	"strings"

	"github.com/cespare/hasty"
	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	var monkeys0, monkeys1 []*monkey
	for _, text := range strings.Split(string(ctx.readAll()), "\n\n") {
		monkeys0 = append(monkeys0, parseMonkey(text))
		monkeys1 = append(monkeys1, parseMonkey(text))
	}
	gang0 := newMonkeyGang(monkeys0)
	gang1 := newMonkeyGang(monkeys1)
	ctx.reportLoad()

	for i := 0; i < 20; i++ {
		gang0.round(1)
	}
	ctx.reportPart1(gang0.level())

	for i := 0; i < 10000; i++ {
		gang1.round(2)
	}
	ctx.reportPart2(gang1.level())
}

type monkeyGang struct {
	monkeys []*monkey
	div     int64
}

func newMonkeyGang(monkeys []*monkey) *monkeyGang {
	div := int64(1)
	for i, m := range monkeys {
		if m.id != i {
			panic("unexpected monkey ID order")
		}
		div *= m.div
	}
	return &monkeyGang{monkeys: monkeys, div: div}
}

func (g *monkeyGang) round(part int) {
	for i := range g.monkeys {
		g.do(i, part)
	}
}

func (g *monkeyGang) do(i, part int) {
	m := g.monkeys[i]
	for len(m.items) > 0 {
		item := m.items[0]
		m.items = m.items[1:]
		m.numInspects++
		item = m.op(item)
		if part == 1 {
			item /= 3
		} else {
			item = item % g.div
		}
		var next *monkey
		if item%m.div == 0 {
			next = g.monkeys[m.throwTrue]
		} else {
			next = g.monkeys[m.throwFalse]
		}
		next.items = append(next.items, item)
	}
}

func (g *monkeyGang) level() int {
	numInspects := SliceMap(g.monkeys, func(m *monkey) int { return m.numInspects })
	slices.SortFunc(numInspects, func(x, y int) bool { return x > y })
	return numInspects[0] * numInspects[1]
}

type monkey struct {
	id          int
	items       []int64
	op          func(int64) int64
	div         int64
	throwTrue   int
	throwFalse  int
	numInspects int
}

var monkeyRegexp = regexp.MustCompile(`Monkey (?P<ID>\d+):
  Starting items: (?P<Items>.*)
  Operation: new = old (?P<Op>.) (?P<Arg>.+)
  Test: divisible by (?P<Div>\d+)
    If true: throw to monkey (?P<ThrowTrue>\d+)
    If false: throw to monkey (?P<ThrowFalse>\d+)`)

func parseMonkey(text string) *monkey {
	var v struct {
		ID         int
		Items      string
		Op         string
		Arg        string
		Div        int64
		ThrowTrue  int
		ThrowFalse int
	}
	hasty.MustParse([]byte(text), &v, monkeyRegexp)
	m := &monkey{
		id:         v.ID,
		items:      SliceMap(strings.Split(v.Items, ", "), parseInt),
		div:        v.Div,
		throwTrue:  v.ThrowTrue,
		throwFalse: v.ThrowFalse,
	}
	switch v.Op {
	case "*":
		if v.Arg == "old" {
			m.op = func(n int64) int64 { return n * n }
		} else {
			x := parseInt(v.Arg)
			m.op = func(n int64) int64 { return n * x }
		}
	case "+":
		if v.Arg == "old" {
			m.op = func(n int64) int64 { return n + n }
		} else {
			x := parseInt(v.Arg)
			m.op = func(n int64) int64 { return n + x }
		}
	default:
		panic("unhandled op")
	}
	return m
}
