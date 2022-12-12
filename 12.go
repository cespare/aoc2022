package main

import (
	"math"
)

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	var m heightmap
	scanner := ctx.scanner()
	for scanner.scan() {
		row := []byte(scanner.text())
		for i, c := range row {
			v := vec2{x: int64(i), y: m.g.rows}
			switch c {
			case 'S':
				m.start = v
				row[i] = 'a'
			case 'E':
				m.end = v
				row[i] = 'z'
			}
		}
		m.g.addRow(row)
	}
	ctx.reportLoad()

	m.fill()
	ctx.reportPart1(m.dists[m.start])

	best := math.MaxInt
	m.g.forEach(func(v vec2, h byte) {
		if h != 'a' {
			return
		}
		if d, ok := m.dists[v]; ok && d < best {
			best = d
		}
	})
	ctx.reportPart2(best)
}

type heightmap struct {
	g     grid[byte]
	start vec2
	end   vec2
	dists map[vec2]int
}

func (m *heightmap) fill() {
	m.dists = map[vec2]int{m.end: 0}
	q := []vec2{m.end}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		d := m.dists[v]
		for _, n := range m.neighbors(v) {
			if _, ok := m.dists[n]; ok {
				continue
			}
			m.dists[n] = d + 1
			q = append(q, n)
		}
	}
}

func (m *heightmap) neighbors(v vec2) []vec2 {
	h0 := m.g.at(v)
	var neighbors []vec2
	for _, n := range v.neighbors4() {
		if !m.g.contains(n) {
			continue
		}
		h := m.g.at(n)
		if h < h0-1 {
			continue
		}
		neighbors = append(neighbors, n)
	}
	return neighbors
}
