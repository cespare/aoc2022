package main

import (
	"math"
	"strings"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var m monkeyMaze
	scanner := ctx.scanner()
	for scanner.scan() {
		if scanner.text() == "" {
			break
		}
		m.add(scanner.text())
	}
	m.finalize()
	var pathLine string
	for scanner.scan() {
		if pathLine != "" {
			panic("multiple path lines")
		}
		pathLine = scanner.text()
	}
	if pathLine == "" {
		panic("no path line")
	}
	path := parseMazePath(pathLine)
	ctx.reportLoad()

	for _, x := range path {
		switch x := x.(type) {
		case byte:
			m.turn(x)
		case int64:
			m.move(x)
		default:
			panic("bad")
		}
	}
	ctx.reportPart1(m.password())

	m.reset()
	for _, x := range path {
		switch x := x.(type) {
		case byte:
			m.turn(x)
		case int64:
			m.cubeMove(x)
		default:
			panic("bad")
		}
	}
	ctx.reportPart2(m.password())
}

type monkeyMaze struct {
	m         map[vec2]byte
	bound     vec2               // bottom-right corner, exclusive
	lrBorders map[int64][2]int64 // inclusive
	udBorders map[int64][2]int64 // inclusive

	edgeLen int // cube size

	cur    vec2
	facing mazeDirection
}

type mazeDirection int

const (
	right mazeDirection = 0 // values also scores
	down  mazeDirection = 1
	left  mazeDirection = 2
	up    mazeDirection = 3
)

func (d mazeDirection) String() string {
	switch d {
	case left:
		return "<"
	case right:
		return ">"
	case up:
		return "^"
	case down:
		return "v"
	default:
		panic("bad")
	}
}

var mazeVecs = []vec2{
	right: {1, 0},
	down:  {0, 1},
	left:  {-1, 0},
	up:    {0, -1},
}

func (m *monkeyMaze) add(line string) {
	y := m.bound.y
	if m.m == nil {
		m.m = make(map[vec2]byte)
	}
	for x := int64(0); x < int64(len(line)); x++ {
		c := line[x]
		if c == ' ' {
			continue
		}
		m.m[vec2{x, y}] = line[x]
	}
	if max := int64(len(line)); max > m.bound.x {
		m.bound.x = max
	}
	m.bound.y++
}

func (m *monkeyMaze) finalize() {
	m.lrBorders = make(map[int64][2]int64)
	m.udBorders = make(map[int64][2]int64)
	for y := int64(0); y < m.bound.y; y++ {
		var border [2]int64
		for x := int64(0); x < m.bound.x; x++ {
			if _, ok := m.m[vec2{x, y}]; ok {
				border[0] = x
				break
			}
		}
		for x := m.bound.x - 1; x >= 0; x-- {
			if _, ok := m.m[vec2{x, y}]; ok {
				border[1] = x
				break
			}
		}
		m.lrBorders[y] = border
	}
	for x := int64(0); x < m.bound.x; x++ {
		var border [2]int64
		for y := int64(0); y < m.bound.y; y++ {
			if _, ok := m.m[vec2{x, y}]; ok {
				border[0] = y
				break
			}
		}
		for y := m.bound.y - 1; y >= 0; y-- {
			if _, ok := m.m[vec2{x, y}]; ok {
				border[1] = y
				break
			}
		}
		m.udBorders[x] = border
	}
	m.edgeLen = int(math.Sqrt(float64(len(m.m) / 6)))
	m.reset()
}

func (m *monkeyMaze) reset() {
	m.cur = vec2{m.lrBorders[0][0], 0}
	m.facing = right
}

func parseMazePath(line string) []any {
	var path []any
	for len(line) > 0 {
		switch line[0] {
		case 'L', 'R':
			path = append(path, line[0])
			line = line[1:]
			continue
		}
		i := strings.IndexFunc(line, func(r rune) bool {
			return r < '0' || r > '9'
		})
		var ns string
		if i >= 0 {
			ns = line[:i]
			line = line[i:]
		} else {
			ns = line
			line = ""
		}
		path = append(path, parseInt(ns))
	}
	return path
}

func (m *monkeyMaze) turn(c byte) {
	var d mazeDirection
	switch c {
	case 'L':
		switch m.facing {
		case left:
			d = down
		case right:
			d = up
		case up:
			d = left
		case down:
			d = right
		}
	case 'R':
		switch m.facing {
		case left:
			d = up
		case right:
			d = down
		case up:
			d = right
		case down:
			d = left
		}
	default:
		panic("bad dir")
	}
	m.facing = d
}

func (m *monkeyMaze) move(n int64) {
	for i := int64(0); i < n; i++ {
		if !m.move1() {
			break
		}
	}
}

func (m *monkeyMaze) move1() bool {
	v := m.cur.add(mazeVecs[m.facing])
	c, ok := m.m[v]
	if !ok {
		switch m.facing {
		case left:
			v.x = m.lrBorders[v.y][1]
		case right:
			v.x = m.lrBorders[v.y][0]
		case up:
			v.y = m.udBorders[v.x][1]
		case down:
			v.y = m.udBorders[v.x][0]
		}
		c, ok = m.m[v]
		if !ok {
			panic("bad")
		}
	}
	switch c {
	case '#':
		return false
	case '.':
		m.cur = v
		return true
	default:
		panic("bad")
	}
}

func (m *monkeyMaze) cubeMove(n int64) {
	for i := int64(0); i < n; i++ {
		if !m.cubeMove1() {
			break
		}
	}
}

func (m *monkeyMaze) cubeMove1() bool {
	v := m.cur.add(mazeVecs[m.facing])
	newFacing := m.facing
	c, ok := m.m[v]
	if !ok {
		v, newFacing = m.crossCubeEdge()
		c, ok = m.m[v]
		if !ok {
			panic("bad")
		}
	}
	switch c {
	case '#':
		return false
	case '.':
		m.cur = v
		m.facing = newFacing
		return true
	default:
		panic("bad")
	}
}

func (m *monkeyMaze) crossCubeEdge() (vec2, mazeDirection) {
	// Better living through hard coding
	switch m.edgeLen {
	case 4:
		return m.crossCubeEdge4()
	case 50:
		return m.crossCubeEdge50()
	default:
		panic("bad")
	}
}

func (m *monkeyMaze) crossCubeEdge4() (vec2, mazeDirection) {
	fx, fy := m.cur.x/4, m.cur.y/4
	switch [2]int64{fx, fy} {
	case [2]int64{2, 0}:
		switch m.facing {
		case left:
			panic("unimplemented")
		case right:
			panic("unimplemented")
		case up:
			panic("unimplemented")
		default:
			panic("bad")
		}
	case [2]int64{0, 1}:
		switch m.facing {
		case left:
			panic("unimplemented")
		case up:
			panic("unimplemented")
		case down:
			panic("unimplemented")
		default:
			panic("bad")
		}
	case [2]int64{1, 1}:
		switch m.facing {
		case up:
			return vec2{8, 3 - (7 - m.cur.x)}, right
		case down:
			panic("unimplemented")
		default:
			panic("bad")
		}
	case [2]int64{2, 1}:
		switch m.facing {
		case right:
			return vec2{7 - m.cur.y + 12, 8}, down
		default:
			panic("bad")
		}
	case [2]int64{2, 2}:
		switch m.facing {
		case left:
			panic("unimplemented")
		case down:
			return vec2{11 - m.cur.x + 0, 7}, up
		default:
			panic("bad")
		}
	case [2]int64{3, 2}:
		switch m.facing {
		case right:
			panic("unimplemented")
		case up:
			panic("unimplemented")
		case down:
			panic("unimplemented")
		default:
			panic("bad")
		}
	default:
		panic("bad")
	}
}

func (m *monkeyMaze) crossCubeEdge50() (vec2, mazeDirection) {
	x, y := m.cur.x, m.cur.y
	fx, fy := x/50, y/50
	switch [2]int64{fx, fy} {
	case [2]int64{1, 0}: // A
		switch m.facing {
		case left: // to D
			return vec2{0, 100 + (49 - y)}, right
		case up: // to F
			return vec2{0, 150 + (x - 50)}, right
		default:
			panic("bad")
		}
	case [2]int64{2, 0}: // B
		switch m.facing {
		case right: // to E
			return vec2{99, 149 - (y)}, left
		case up: // to F
			return vec2{0 + (x - 100), 199}, up
		case down: // to C
			return vec2{99, 99 - (149 - x)}, left
		default:
			panic("bad")
		}
	case [2]int64{1, 1}: // C
		switch m.facing {
		case left: // to D
			return vec2{0 + (y - 50), 100}, down
		case right: // to B
			return vec2{149 - (99 - y), 49}, up
		default:
			panic("bad")
		}
	case [2]int64{0, 2}: // D
		switch m.facing {
		case left: // to A
			return vec2{50, 49 - (y - 100)}, right
		case up: // to C
			return vec2{50, 50 + (x)}, right
		default:
			panic("bad")
		}
	case [2]int64{1, 2}: // E
		switch m.facing {
		case right: // to B
			return vec2{149, 0 + (149 - y)}, left
		case down: // to F
			return vec2{49, 150 + (x - 50)}, left
		default:
			panic("bad")
		}
	case [2]int64{0, 3}: // F
		switch m.facing {
		case left: // to A
			return vec2{50 + (y - 150), 0}, down
		case right: // to E
			return vec2{50 + (y - 150), 149}, up
		case down: // to B
			return vec2{100 + (x), 0}, down
		default:
			panic("bad")
		}
	default:
		panic("bad")
	}
}

func (m *monkeyMaze) password() int64 {
	row := m.cur.y + 1
	col := m.cur.x + 1
	return 1000*row + 4*col + int64(m.facing)
}
