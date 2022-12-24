package main

import (
	"strings"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(24, problem24)
}

func problem24(ctx *problemContext) {
	var m snowMaze
	scanner := ctx.scanner()
	for scanner.scan() {
		m.addLine(scanner.text())
	}
	ctx.reportLoad()

	t0 := m.solve(m.start, m.end, 0)
	ctx.reportPart1(t0)

	t1 := m.solve(m.end, m.start, t0)
	t2 := m.solve(m.start, m.end, t1)
	ctx.reportPart1(t2)
}

type snowMaze struct {
	blizzards []map[vec2][]byte // by minute
	start     vec2
	end       vec2
	bound     vec2 // lower-right wall corner, inclusive
	lcm       int  // for blizzard state
	y         int64
}

func (m *snowMaze) addLine(line string) {
	if len(m.blizzards) == 0 {
		m.blizzards = []map[vec2][]byte{make(map[vec2][]byte)}
		m.bound.x = int64(len(line)) - 1
		if line[1] != '.' {
			panic("bad")
		}
		m.start = vec2{1, 0}
		return
	}
	m.y++
	if int64(len(line)) != m.bound.x+1 {
		panic("bad")
	}
	if line[1] == '#' {
		i := strings.IndexByte(line, '.')
		if i < 0 {
			panic("bad")
		}
		m.end = vec2{int64(i), m.y}
		m.bound.y = m.y
		w := m.bound.x - 1
		h := m.bound.y - 1
		for k := int64(1); ; k++ {
			if x := k * w; x%h == 0 {
				m.lcm = int(x)
				break
			}
		}
		return
	}
	for x := int64(0); x < int64(len(line)); x++ {
		switch line[x] {
		case '#':
			if x != 0 && x != int64(len(line))-1 {
				panic("bad")
			}
		case '.':
		case '<', '^', '>', 'v':
			m.blizzards[0][vec2{x, m.y}] = []byte{line[x]}
		default:
			panic("bad")
		}
	}
}

func (m *snowMaze) solve(start, end vec2, t0 int) int {
	type state struct {
		p      vec2
		minute int
	}
	q := []state{{p: start, minute: t0}}
	seen := set.Of(q[0])
	enqueue := func(s state) {
		s1 := s
		s1.minute = s1.minute % m.lcm
		if seen.Contains(s1) {
			return
		}
		seen.Add(s1)
		q = append(q, s)
	}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]

		blizzards := m.getBlizzards(s.minute + 1)
		candidates := append([]vec2{s.p}, s.p.neighbors4()...)
		next := s
		next.minute++
		for _, c := range candidates {
			if c == end {
				return next.minute
			}
			if c.x <= 0 || c.x >= m.bound.x || c.y <= 0 || c.y >= m.bound.y {
				if c != start {
					continue
				}
			}
			if _, ok := blizzards[c]; ok {
				continue
			}
			next.p = c
			enqueue(next)
		}
	}
	panic("no solution")
}

func (m *snowMaze) getBlizzards(minute int) map[vec2][]byte {
	minute = minute % m.lcm
	for len(m.blizzards) < minute+1 {
		b := m.blizzards[len(m.blizzards)-1]
		m.blizzards = append(m.blizzards, m.nextBlizzard(b))
	}
	return m.blizzards[minute]
}

func (m *snowMaze) nextBlizzard(b0 map[vec2][]byte) map[vec2][]byte {
	b1 := make(map[vec2][]byte)
	for v, cs := range b0 {
		for _, c := range cs {
			v1 := v
			switch c {
			case '<':
				v1.x--
				if v1.x == 0 {
					v1.x = m.bound.x - 1
				}
			case '^':
				v1.y--
				if v1.y == 0 {
					v1.y = m.bound.y - 1
				}
			case '>':
				v1.x++
				if v1.x == m.bound.x {
					v1.x = 1
				}
			case 'v':
				v1.y++
				if v1.y == m.bound.y {
					v1.y = 1
				}
			default:
				panic("bad")
			}
			b1[v1] = append(b1[v1], c)
		}
	}
	return b1
}
