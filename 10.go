package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	var insts []clockInstruction
	scanner := ctx.scanner()
	for scanner.scan() {
		var inst clockInstruction
		name, rest, _ := strings.Cut(scanner.text(), " ")
		inst.name = name
		switch name {
		case "noop":
		case "addx":
			inst.v = parseInt(rest, 10, 64)
		default:
			panic("bad instruction")
		}
		insts = append(insts, inst)
	}
	ctx.reportLoad()

	cpu := newClockCPU()
	var signalSum int64
	var part2 strings.Builder
	var pix int
	for cycle := 1; cycle <= 240; cycle++ {
		switch cycle {
		case 20, 60, 100, 140, 180, 220:
			signalSum += int64(cycle) * cpu.x
		}
		if abs(int64(pix)-cpu.x) <= 1 {
			part2.WriteByte('#')
		} else {
			part2.WriteByte('.')
		}
		pix++
		if pix == 40 {
			part2.WriteByte('\n')
			pix = 0
		}

		var inst clockInstruction
		if cycle <= len(insts) {
			inst = insts[cycle-1]
		} else {
			inst = clockInstruction{name: "noop"}
		}
		cpu.tick(inst)
	}
	ctx.reportPart1(signalSum)
	fmt.Println(strings.TrimSpace(part2.String()))
}

type clockInstruction struct {
	name string
	v    int64
}

func (i clockInstruction) cycles() int {
	switch i.name {
	case "noop":
		return 1
	case "addx":
		return 2
	default:
		panic("bad instruction")
	}
}

type clockCPU struct {
	rem   int // cycles left to retire queue[0]
	x     int64
	queue []clockInstruction
}

func newClockCPU() *clockCPU {
	return &clockCPU{x: 1}
}

func (c *clockCPU) tick(inst clockInstruction) {
	c.queue = append(c.queue, inst)
	if c.rem == 0 { // initial state
		c.rem = c.queue[0].cycles()
	}
	c.rem--
	if c.rem > 0 {
		return
	}
	inst, c.queue = c.queue[0], c.queue[1:]
	if inst.name == "addx" {
		c.x += inst.v
	}
}
