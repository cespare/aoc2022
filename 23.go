package main

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	var g elfGrid
	scanner := ctx.scanner()
	for scanner.scan() {
		g.addLine(scanner.text())
	}
	ctx.reportLoad()

	for round := 1; ; round++ {
		moved := g.move()
		if round == 10 {
			ctx.reportPart1(g.answer())
		}
		if moved == 0 {
			ctx.reportPart2(round)
			return
		}
	}
}

type elfGrid struct {
	m            map[vec2]vec2 // pos -> proposed
	targetCounts map[vec2]int  // scratch; reused
	y            int64
	startDir     compassDir
}

func (g *elfGrid) addLine(line string) {
	if g.m == nil {
		g.m = make(map[vec2]vec2)
	}
	for x := int64(0); x < int64(len(line)); x++ {
		if line[x] != '#' {
			continue
		}
		v := vec2{x, g.y}
		g.m[v] = v
	}
	g.y++
}

type compassDir int

const (
	north compassDir = iota
	south
	west
	east
)

func (g *elfGrid) answer() int64 {
	var bounds [2]vec2
	first := true
	for v := range g.m {
		if first || v.x < bounds[0].x {
			bounds[0].x = v.x
		}
		if first || v.y < bounds[0].y {
			bounds[0].y = v.y
		}
		if first || v.x > bounds[1].x {
			bounds[1].x = v.x
		}
		if first || v.y > bounds[1].y {
			bounds[1].y = v.y
		}
		first = false
	}

	var res int64
	for y := bounds[0].y; y <= bounds[1].y; y++ {
		for x := bounds[0].x; x <= bounds[1].x; x++ {
			if _, ok := g.m[vec2{x, y}]; !ok {
				res++
			}
		}
	}
	return res
}

var elfDirs = [...]vec2{
	{-1, -1},
	{0, -1},
	{1, -1},
	{1, 0},
	{1, 1},
	{0, 1},
	{-1, 1},
	{-1, 0},
}

var elfDirChecks = [...][3]int{
	north: {0, 1, 2},
	south: {6, 5, 4},
	west:  {6, 7, 0},
	east:  {2, 3, 4},
}

func (g *elfGrid) move() (moved int64) {
	for v := range g.m {
		var hasNeighbor bool
		// Use a consistent ordering to refer to the 8 neighbors:
		// start at northwest; circle CW.
		occ := make([]bool, 8)
		for i, d := range elfDirs {
			n := v.add(d)
			if _, ok := g.m[n]; ok {
				hasNeighbor = true
				occ[i] = true
			}
		}
		if !hasNeighbor {
			continue
		}

		dir := g.startDir
		for i := 0; i < 4; i++ {
			toCheck := elfDirChecks[dir]
			var blocked bool
			for _, j := range toCheck {
				if occ[j] {
					blocked = true
					break
				}
			}
			if !blocked {
				g.m[v] = v.add(elfDirs[toCheck[1]])
				break
			}
			dir = (dir + 1) % 4
		}
	}
	if g.targetCounts == nil {
		g.targetCounts = make(map[vec2]int)
	}
	for v := range g.targetCounts {
		delete(g.targetCounts, v)
	}
	for _, t := range g.m {
		g.targetCounts[t]++
	}
	for v, t := range g.m {
		if g.targetCounts[t] > 1 {
			continue
		}
		if t != v {
			delete(g.m, v)
			g.m[t] = t
			moved++
		}
	}
	g.startDir = (g.startDir + 1) % 4
	return moved
}
