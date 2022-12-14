package main

import (
	"strings"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	var c sandCave
	scanner := ctx.scanner()
	for scanner.scan() {
		vs := parseCaveWalls(scanner.text())
		v := vs[0]
		for _, v1 := range vs[1:] {
			d := v1.sub(v)
			if d.x == 0 {
				d = vec2{0, d.y / abs(d.y)}
			} else {
				d = vec2{d.x / abs(d.x), 0}
			}
			for {
				c.add(v, '#')
				if v == v1 {
					break
				}
				v = v.add(d)
			}
		}
	}
	c.floor = c.max.y + 2
	ctx.reportLoad()

	var grains int
	for c.addSand() != sandOnFloor {
		grains++
	}
	ctx.reportPart1(grains)

	for {
		grains++
		if c.addSand() == sandFull {
			break
		}
	}
	ctx.reportPart2(grains)
}

func parseCaveWalls(line string) []vec2 {
	var vs []vec2
	for _, pair := range strings.Split(line, " -> ") {
		x, y, _ := strings.Cut(pair, ",")
		vs = append(vs, vec2{parseInt(x), parseInt(y)})
	}
	return vs
}

type sandCave struct {
	m     map[vec2]byte
	min   vec2
	max   vec2
	floor int64
}

func (c *sandCave) add(v vec2, b byte) {
	if c.m == nil {
		c.m = make(map[vec2]byte)
		c.min = v
		c.max = v
	} else {
		if v.x < c.min.x {
			c.min.x = v.x
		}
		if v.y < c.min.y {
			c.min.y = v.y
		}
		if v.x > c.max.x {
			c.max.x = v.x
		}
		if v.y > c.max.y {
			c.max.y = v.y
		}
	}
	c.m[v] = b
}

func (c *sandCave) at(v vec2) byte {
	if v.y >= c.floor {
		return '_'
	}
	if b, ok := c.m[v]; ok {
		return b
	}
	return '.'
}

type sandResult int

const (
	sandAtRest sandResult = iota
	sandOnFloor
	sandFull
)

func (c *sandCave) addSand() sandResult {
	v := vec2{500, 0}
	if c.at(v) != '.' {
		return sandFull
	}
	result := sandAtRest
fall:
	for {
	try:
		for _, d := range []vec2{
			{0, 1},
			{-1, 1},
			{1, 1},
		} {
			v1 := v.add(d)
			switch c.at(v1) {
			case '.':
				v = v1
				continue fall
			case '_':
				result = sandOnFloor
				break try
			}
		}
		c.add(v, 'o')
		return result
	}
}
