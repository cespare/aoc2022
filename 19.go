package main

import (
	"math"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/cespare/hasty"
)

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	blueprints := scanSlice(ctx, parseBlueprint)
	ctx.reportLoad()

	var part1 atomic.Int64
	ParDo(blueprints, func(bp blueprint) {
		part1.Add(bp.id * mostGeodes(bp, 24))
	})
	ctx.reportPart1(part1.Load())

	if len(blueprints) > 3 {
		blueprints = blueprints[:3]
	}
	part2 := int64(1)
	var mu sync.Mutex
	ParDo(blueprints, func(bp blueprint) {
		m := mostGeodes(bp, 32)
		mu.Lock()
		part2 *= m
		mu.Unlock()
	})
	ctx.reportPart2(part2)
}

type resource int

const (
	ore resource = iota
	clay
	obsidian
	geode
	numResources
)

func parseResource(s string) resource {
	switch s {
	case "ore":
		return ore
	case "clay":
		return clay
	case "obsidian":
		return obsidian
	case "geode":
		return geode
	default:
		panic("bad")
	}
}

type blueprint struct {
	id    int64
	costs [numResources][numResources]int64
}

var blueprintRegexps = [...]*regexp.Regexp{
	regexp.MustCompile(`Blueprint (?P<ID>\d+): (?P<Rest>.*)$`),
	regexp.MustCompile(`Each (?P<Output>\w+) robot costs (?P<Rest>.*).?$`),
	regexp.MustCompile(`(?P<Amount>\d+) (?P<Resource>\w+)`),
}

func parseBlueprint(line string) blueprint {
	var v0 struct {
		ID   int64
		Rest string
	}
	hasty.MustParse([]byte(line), &v0, blueprintRegexps[0])
	bp := blueprint{id: v0.ID}
	for _, part := range strings.Split(v0.Rest, ". ") {
		var v1 struct {
			Output string
			Rest   string
		}
		hasty.MustParse([]byte(part), &v1, blueprintRegexps[1])
		output := parseResource(v1.Output)
		for _, input := range strings.Split(v1.Rest, " and ") {
			var v2 struct {
				Amount   int64
				Resource string
			}
			hasty.MustParse([]byte(input), &v2, blueprintRegexps[2])
			input := parseResource(v2.Resource)
			bp.costs[output][input] = v2.Amount
		}
	}
	return bp
}

type geodeBotState struct {
	bots      [numResources]int64
	resources [numResources]int64
	minutes   int
}

type geodeSolver struct {
	bp        blueprint
	maxUseful [numResources]int64 // enough to build any bot
	bests     map[geodeBotState]int64
	bestSoFar int64
}

func mostGeodes(bp blueprint, minutes int) int64 {
	solver := &geodeSolver{
		bp:    bp,
		bests: make(map[geodeBotState]int64),
	}
	for _, c := range bp.costs {
		for in, n := range c {
			solver.maxUseful[in] = maxim(solver.maxUseful[in], n)
		}
	}
	solver.maxUseful[geode] = math.MaxInt64 // always useful
	var g geodeBotState
	g.bots[ore] = 1
	g.minutes = minutes
	return solver.solve(g)
}

func (s *geodeSolver) solve(g geodeBotState) (best int64) {
	if b, ok := s.bests[g]; ok {
		return b
	}
	defer func() {
		s.bests[g] = best
		if best > s.bestSoFar {
			s.bestSoFar = best
		}
	}()

	if g.minutes == 0 {
		return g.resources[geode]
	}
	// Compute a loose upper bound on the maximum possible score still
	// available to us (assume we build a geode bot every turn). If that is
	// less than the best score we've seen, toss this one out.
	geodes := g.resources[geode]
	geodeBots := g.bots[geode]
	for m := g.minutes; m > 0; m-- {
		geodes += geodeBots
		geodeBots++
	}
	if geodes <= s.bestSoFar {
		return 0
	}

	next := g
	next.minutes--
	for r, n := range g.bots {
		next.resources[r] += n
	}

	buy := func(out resource) (next1 geodeBotState, ok bool) {
		next1 = next
		for in, n := range s.bp.costs[out] {
			if g.resources[in] < n { // check current resources, not next ones
				return next1, false
			}
			next1.resources[in] -= n
		}
		next1.bots[out]++
		return next1, true
	}
	// Try not buying anything.
	best = s.solve(next)
	// Try buying one of each bot in turn.
	for out := resource(0); out < numResources; out++ {
		if g.bots[out] >= s.maxUseful[out] {
			continue
		}
		if next1, ok := buy(out); ok {
			if b := s.solve(next1); b > best {
				best = b
			}
		}
	}
	return best
}
