package main

import (
	"regexp"

	"github.com/cespare/hasty"
	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	var g sensorGrid
	scanner := ctx.scanner()
	for scanner.scan() {
		s, b := parseSensorLine(scanner.text())
		g.add(s, b)
	}
	ctx.reportLoad()

	ctx.reportPart1(g.part1(2000000))
	const size = 4_000_000
	ctx.reportPart2(g.part2(size))
}

var sensorRegexp = regexp.MustCompile(`Sensor at x=(?P<SX>-?\d+), y=(?P<SY>-?\d+): closest beacon is at x=(?P<BX>-?\d+), y=(?P<BY>-?\d+)`)

func parseSensorLine(line string) (s, b vec2) {
	var v struct {
		SX int64
		SY int64
		BX int64
		BY int64
	}
	hasty.MustParse([]byte(line), &v, sensorRegexp)
	return vec2{v.SX, v.SY}, vec2{v.BX, v.BY}
}

type sensorGrid struct {
	m map[vec2]vec2 // sensor -> beacon
}

func (g *sensorGrid) add(s, b vec2) {
	if g.m == nil {
		g.m = make(map[vec2]vec2)
	}
	g.m[s] = b
}

func (g *sensorGrid) part1(y int64) int64 {
	var ivs intervalSet
	beacons := make(map[int64]struct{})
	for s, b := range g.m {
		if b.y == y {
			beacons[b.x] = struct{}{}
		}
		d := b.sub(s).mag()
		dy := abs(s.y - y)
		dx := d - dy
		if dx < 0 {
			continue
		}
		ivs.add(interval{s.x - dx, s.x + dx})
	}
	return ivs.countOccupied() - int64(len(beacons))
}

func (g *sensorGrid) find(bounds interval, y int64) int64 {
	var ivs intervalSet
	for s, b := range g.m {
		d := b.sub(s).mag()
		dy := abs(s.y - y)
		dx := d - dy
		if dx < 0 {
			continue
		}
		iv := interval{
			maxim(s.x-dx, bounds.min),
			minim(s.x+dx, bounds.max),
		}
		ivs.add(iv)
	}
	return ivs.unoccupied(bounds)
}

func (g *sensorGrid) part2(size int64) int64 {
	ivs := make([]intervalSet, size+1)
	for s, b := range g.m {
		d := b.sub(s).mag()
		ymin := maxim(s.y-d, 0)
		ymax := minim(s.y+d, size)
		for y := ymin; y <= ymax; y++ {
			dy := abs(s.y - y)
			iv := interval{
				maxim(s.x-(d-dy), 0),
				minim(s.x+(d-dy), size),
			}
			ivs[y].add(iv)
		}
	}
	for y, set := range ivs {
		if x := set.unoccupied(interval{0, size}); x != -1 {
			return x*4000000 + int64(y)
		}
	}
	panic("not found")
}

type interval struct {
	min int64
	max int64 // inclusive
}

type intervalSet []interval // non-overlapping

func (s *intervalSet) add(iv interval) {
	if iv.max-iv.min < 0 {
		panic("bad")
	}
	if iv.max-iv.min == 0 {
		return
	}
	for i, iv0 := range *s {
		if iv.min > iv0.max {
			continue
		}
		if iv.max < iv0.min {
			*s = slices.Insert(*s, i, iv)
			return
		}
		// Overlapping
		iv.min = minim(iv.min, iv0.min)
		iv.max = maxim(iv.max, iv0.max)
		// It might overlap with subsequent intervals as well.
		j := i + 1
		for ; j < len(*s); j++ {
			iv0 := (*s)[j]
			if iv.max < iv0.min {
				break
			}
			iv.max = maxim(iv.max, iv0.max)
		}
		*s = slices.Replace(*s, i, j, iv)
		return
	}
	*s = append(*s, iv)
}

func (s intervalSet) countOccupied() int64 {
	var n int64
	for _, iv := range s {
		n += iv.max - iv.min + 1
	}
	return n
}

func (s intervalSet) unoccupied(bounds interval) int64 {
	if len(s) == 1 && s[0] == bounds {
		return -1
	}
	if len(s) == 0 {
		panic("bad")
	}
	if len(s) > 2 {
		panic("bad")
	}
	if len(s) == 2 {
		d := s[1].min - s[0].max
		switch d {
		case 0, 1:
			panic("weird")
		case 2:
			return s[0].max + 1
		default:
			panic("bad")
		}
	}
	d := s[0].min - bounds.min
	switch d {
	case 0:
	case 1:
		return bounds.min
	default:
		panic("bad")
	}
	d = bounds.max - s[0].max
	switch d {
	case 0:
	case 1:
		return bounds.max
	default:
		panic("bad")
	}
	panic("bad")
}
