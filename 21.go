package main

import (
	"strconv"
	"strings"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	var w monkeyWorkers
	scanner := ctx.scanner()
	for scanner.scan() {
		w.add(parseMonkey2(scanner.text()))
	}
	ctx.reportLoad()

	w.populate()

	ctx.reportPart1(w.root.v)
	ctx.reportPart2(w.reverseSolve())
}

type monkey2 struct {
	name      string
	v         int64
	op        string
	leftName  string
	rightName string

	// The remaining fields filled in by populate (along with v, if op!="").
	left    *monkey2
	right   *monkey2
	tainted bool // depends on the value of humn
}

func parseMonkey2(line string) *monkey2 {
	name, rest, ok := strings.Cut(line, ": ")
	if !ok {
		panic("bad")
	}
	m := &monkey2{name: name}
	if n, err := strconv.ParseInt(rest, 10, 64); err == nil {
		m.v = n
		return m
	}
	parts := strings.Split(rest, " ")
	if len(parts) != 3 || len(parts[1]) != 1 {
		panic("bad")
	}
	m.leftName = parts[0]
	m.op = parts[1]
	m.rightName = parts[2]
	return m
}

type monkeyWorkers struct {
	byName map[string]*monkey2
	root   *monkey2
}

func (w *monkeyWorkers) add(m *monkey2) {
	if w.byName == nil {
		w.byName = make(map[string]*monkey2)
	}
	w.byName[m.name] = m
	if m.name == "root" {
		w.root = m
	}
}

func (w *monkeyWorkers) populate() {
	for _, m := range w.byName {
		if m.op != "" {
			m.left = w.byName[m.leftName]
			m.right = w.byName[m.rightName]
		}
	}
	_ = w.root.computeValue()
	_ = w.root.computeTainted()
}

func (m *monkey2) computeValue() int64 {
	if m.op == "" {
		return m.v
	}
	left := m.left.computeValue()
	right := m.right.computeValue()
	switch m.op {
	case "+":
		m.v = left + right
	case "*":
		m.v = left * right
	case "-":
		m.v = left - right
	case "/":
		m.v = left / right
	default:
		panic("bad")
	}
	return m.v
}

func (m *monkey2) computeTainted() bool {
	if m.name == "humn" {
		m.tainted = true
		return true
	}
	if m.op == "" {
		return false
	}
	m.tainted = m.left.computeTainted() || m.right.computeTainted()
	return m.tainted
}

func (w *monkeyWorkers) reverseSolve() int64 {
	if w.root.left.tainted {
		if w.root.right.tainted {
			panic("two tainted branches")
		}
		return w.root.left.reverseSolve(w.root.right.v)
	} else {
		if !w.root.right.tainted {
			panic("no tainted branches")
		}
		return w.root.right.reverseSolve(w.root.left.v)
	}
}

func (m *monkey2) reverseSolve(want int64) int64 {
	if m.name == "humn" {
		return want
	}
	if m.op == "" {
		panic("bad")
	}
	if m.left.tainted {
		if m.right.tainted {
			panic("two tainted branches")
		}
		switch m.op {
		case "+":
			want -= m.right.v
		case "*":
			want /= m.right.v
		case "-":
			want += m.right.v
		case "/":
			want *= m.right.v
		}
		return m.left.reverseSolve(want)
	} else {
		if !m.right.tainted {
			panic("no tainted branches")
		}
		switch m.op {
		case "+":
			want -= m.left.v
		case "*":
			want /= m.left.v
		case "-":
			want = m.left.v - want
		case "/":
			want = m.left.v / want
		}
		return m.right.reverseSolve(want)
	}
}
