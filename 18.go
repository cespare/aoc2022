package main

import (
	"strings"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	var points []vec3
	scanner := ctx.scanner()
	for scanner.scan() {
		v := SliceMap(strings.Split(scanner.text(), ","), parseInt)
		points = append(points, vec3{x: v[0], y: v[1], z: v[2]})
	}
	ctx.reportLoad()

	var s cubeSet
	for _, p := range points {
		s.add(p)
	}
	ctx.reportPart1(s.faces)

	for _, p := range s.interiorPoints() {
		s.add(p)
	}
	ctx.reportPart2(s.faces)
}

type cubeSet struct {
	cubes set.Set[vec3]
	faces int64
}

func (s *cubeSet) add(v vec3) {
	if s.cubes.Contains(v) {
		panic("duplicate")
	}
	s.cubes.Add(v)
	for _, n := range v.neighbors6() {
		if s.cubes.Contains(n) {
			s.faces--
		} else {
			s.faces++
		}
	}
}

func (s *cubeSet) interiorPoints() []vec3 {
	state := &cubeState{
		interior: new(set.Set[vec3]),
		exterior: new(set.Set[vec3]),
	}
	// Find a bounding box.
	s.cubes.Do(func(v vec3) bool {
		state.min = state.min.min(v)
		state.max = state.max.max(v)
		return true
	})
	for z := state.min.z; z <= state.max.z; z++ {
		for y := state.min.y; y <= state.max.y; y++ {
			for x := state.min.x; x <= state.max.x; x++ {
				v := vec3{x, y, z}
				if s.cubes.Contains(v) {
					continue
				}
				s.checkEscape(v, state)
			}
		}
	}
	return state.interior.Slice()
}

type cubeState struct {
	min, max vec3
	interior *set.Set[vec3] // known to be interior
	exterior *set.Set[vec3] // known to be exterior
}

func (s *cubeSet) checkEscape(v vec3, state *cubeState) {
	if state.interior.Contains(v) || state.exterior.Contains(v) {
		return
	}
	// Explore from v until we reach a known state. At that point we put all
	// of wip into one of the two known state sets.
	wip := set.Of(v)
	q := []vec3{v}
	for len(q) > 0 {
		v := SlicePop(&q)
		if v.x < state.min.x || v.y < state.min.y || v.z < state.min.z ||
			v.x > state.max.x || v.y > state.max.y || v.z > state.max.z {
			state.exterior.AddSet(wip)
			return
		}
		if state.interior.Contains(v) {
			panic("not possible")
		}
		if state.exterior.Contains(v) {
			state.exterior.AddSet(wip)
			return
		}
		for _, n := range v.neighbors6() {
			if s.cubes.Contains(n) || wip.Contains(n) {
				continue
			}
			wip.Add(n)
			SlicePush(&q, n)
		}
	}
	state.interior.AddSet(wip)
}
