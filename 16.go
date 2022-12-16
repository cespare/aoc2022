package main

import (
	"math/bits"
	"regexp"
	"sort"
	"strings"

	"github.com/cespare/hasty"
	"github.com/cespare/next/container/heap"
	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	var valves []rawValve
	scanner := ctx.scanner()
	for scanner.scan() {
		valves = append(valves, parseValve(scanner.text()))
	}
	t := buildValveTunnels(valves)
	ctx.reportLoad()

	ctx.reportPart1(t.bestRelease())
	ctx.reportPart2(t.bestRelease2())
}

var valveRegexp = regexp.MustCompile(`Valve (?P<Source>[A-Z]+) has flow rate=(?P<Flow>\d+); tunnels? leads? to valves? (?P<Dests>[A-Z, ]+)$`)

func parseValve(line string) rawValve {
	var v struct {
		Source string
		Flow   int64
		Dests  string
	}
	hasty.MustParse([]byte(line), &v, valveRegexp)
	return rawValve{
		src:  v.Source,
		flow: v.Flow,
		dsts: strings.Split(v.Dests, ", "),
	}
}

type rawValve struct {
	src  string
	flow int64
	dsts []string
}

type valveTunnels struct {
	edges     [][]int
	flow      []int64
	bestFlows []flowValve // non-zero flows sorted largest -> smallest
}

type flowValve struct {
	flow  int64
	valve int
}

func buildValveTunnels(valves []rawValve) *valveTunnels {
	nodes := make([]string, len(valves))
	for i, v := range valves {
		nodes[i] = v.src
	}
	sort.Strings(nodes)
	if nodes[0] != "AA" {
		panic("unexpected")
	}
	byName := make(map[string]int)
	for i, src := range nodes {
		byName[src] = i
	}
	t := &valveTunnels{
		edges: make([][]int, len(valves)),
		flow:  make([]int64, len(valves)),
	}
	for _, v := range valves {
		si := byName[v.src]
		t.flow[si] = v.flow
		for _, dst := range v.dsts {
			di := byName[dst]
			t.edges[si] = append(t.edges[si], di)
		}
		if v.flow > 0 {
			t.bestFlows = append(t.bestFlows, flowValve{v.flow, si})
		}
	}
	slices.SortFunc(t.bestFlows, func(fv0, fv1 flowValve) bool {
		return fv0.flow > fv1.flow
	})
	return t
}

type tunnelState struct {
	rem   int
	pos0  int
	pos1  int // unused for part 1
	open  uint64
	total int64
	max   int64
}

// maxPossible computes a loose upper bound on the best outcome of the current
// state by assuming we visit as many remaining valves as we can, largest to
// smallest, with a single walk in between.
func (t *valveTunnels) maxPossible(ts tunnelState, numWorkers int) int64 {
	if ts.rem < 2 { // need rem>=2 to benefit from opening a valve
		return ts.total
	}
	w := numWorkers
	for _, fv := range t.bestFlows {
		if ts.open&(1<<fv.valve) > 0 {
			continue
		}
		if w == numWorkers {
			ts.rem-- // 1 to open
		}
		ts.total += int64(ts.rem) * fv.flow
		w--
		if w == 0 {
			ts.rem -= 1 // 1 to move
			if ts.rem < 2 {
				break
			}
			w = numWorkers
		}
	}
	return ts.total
}

func (t *valveTunnels) bestRelease() int64 {
	var best int64
	bestByState := make(map[tunnelState]int64) // state{total: 0} -> best total
	q := heap.New(func(ts0, ts1 tunnelState) bool {
		return ts0.max > ts1.max
	})
	pushState := func(ts tunnelState) {
		vs := ts
		vs.total = 0
		if b, ok := bestByState[vs]; !ok || ts.total > b {
			bestByState[vs] = ts.total
			ts.max = t.maxPossible(ts, 1)
			q.Push(ts)
		}
	}
	pushState(tunnelState{rem: 30})
	for {
		ts := q.Pop()
		if ts.max < best {
			return best
		}
		if ts.total > best {
			best = ts.total
		}
		if ts.rem == 0 || bits.OnesCount64(ts.open) == len(t.edges) {
			continue
		}
		flow := t.flow[ts.pos0]
		if flow > 0 && ts.open&(1<<ts.pos0) == 0 {
			ts1 := ts
			ts1.rem--
			ts1.open |= 1 << ts.pos0
			ts1.total += int64(ts1.rem) * flow
			pushState(ts1)
		}
		for _, dst := range t.edges[ts.pos0] {
			ts1 := ts
			ts1.rem--
			ts1.pos0 = dst
			pushState(ts1)
		}
	}
}

func (t *valveTunnels) bestRelease2() int64 {
	var best int64
	bestByState := make(map[tunnelState]int64) // state{pos0<=pos1, total=0} -> best total
	q := heap.New(func(ts0, ts1 tunnelState) bool {
		return ts0.max > ts1.max
	})
	pushState := func(ts tunnelState) {
		vs := ts
		vs.total = 0
		if vs.pos0 > vs.pos1 {
			// We can swap the positions WLOG.
			vs.pos0, vs.pos1 = vs.pos1, vs.pos0
		}
		if b, ok := bestByState[vs]; !ok || ts.total > b {
			bestByState[vs] = ts.total
			ts.max = t.maxPossible(ts, 2)
			q.Push(ts)
		}
	}
	pushState(tunnelState{rem: 26})
	for {
		ts := q.Pop()
		if ts.max < best {
			return best
		}
		if ts.total > best {
			best = ts.total
		}
		if ts.rem == 0 || bits.OnesCount64(ts.open) == len(t.edges) {
			continue
		}
		flow0, flow1 := t.flow[ts.pos0], t.flow[ts.pos1]
		canOpen0 := flow0 > 0 && ts.open&(1<<ts.pos0) == 0
		canOpen1 := flow1 > 0 && ts.open&(1<<ts.pos1) == 0 && ts.pos0 != ts.pos1
		if canOpen0 {
			ts1 := ts
			ts1.rem--
			ts1.open |= 1 << ts.pos0
			ts1.total += int64(ts1.rem) * flow0
			if canOpen1 {
				// Open both
				ts2 := ts1
				ts2.open |= 1 << ts.pos1
				ts2.total += int64(ts1.rem) * flow1
				pushState(ts2)
			}
			// Open pos0, move pos1
			for _, dst := range t.edges[ts.pos1] {
				ts2 := ts1
				ts2.pos1 = dst
				pushState(ts2)
			}
		}
		if canOpen1 {
			// Move pos0, open pos1
			ts1 := ts
			ts1.rem--
			ts1.open |= 1 << ts.pos1
			ts1.total += int64(ts1.rem) * flow1
			for _, dst := range t.edges[ts.pos0] {
				ts2 := ts1
				ts2.pos0 = dst
				pushState(ts2)
			}
		}
		// Move both
		for _, dst0 := range t.edges[ts.pos0] {
			for _, dst1 := range t.edges[ts.pos1] {
				ts1 := ts
				ts1.rem--
				ts1.pos0 = dst0
				ts1.pos1 = dst1
				pushState(ts1)
			}
		}
	}
}
