package main

import (
	"bytes"

	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	var pairs [][2]packetVal
	text := ctx.readAll()
	for _, group := range bytes.Split(text, []byte("\n\n")) {
		lines := bytes.Split(group, []byte("\n"))
		pair := [2]packetVal{
			parsePacket(string(lines[0])),
			parsePacket(string(lines[1])),
		}
		pairs = append(pairs, pair)
	}
	ctx.reportLoad()

	var part1 int
	for i, pair := range pairs {
		if comparePackets(pair[0], pair[1]) < 0 {
			part1 += i + 1
		}
	}
	ctx.reportPart1(part1)

	p2 := parsePacket("[[2]]")
	p6 := parsePacket("[[6]]")
	allPackets := []packetVal{p2, p6}
	for _, pair := range pairs {
		allPackets = append(allPackets, pair[0], pair[1])
	}
	slices.SortFunc(allPackets, func(p0, p1 packetVal) bool {
		return comparePackets(p0, p1) < 0
	})
	part2 := 1
	for i, p := range allPackets {
		if comparePackets(p, p2) == 0 {
			part2 *= (i + 1)
		}
		if comparePackets(p, p6) == 0 {
			part2 *= (i + 1)
		}
	}
	ctx.reportPart2(part2)
}

type packetVal interface{} // int or []packetVal

func parsePacket(text string) packetVal {
	p, tail := parsePacket2(text)
	if tail != "" {
		panic("extra")
	}
	return p
}

func parsePacket2(text string) (p packetVal, tail string) {
	if text[0] != '[' {
		for i := 0; ; i++ {
			if i == len(text) || text[i] < '0' || text[i] > '9' {
				return int(parseInt(text[:i])), text[i:]
			}
		}
	}
	text = text[1:]
	var vals []packetVal
	for {
		if text[0] == ']' {
			return vals, text[1:]
		}
		var child packetVal
		child, text = parsePacket2(text)
		vals = append(vals, child)
		tok := text[0]
		text = text[1:]
		switch tok {
		case ']':
			return vals, text
		case ',':
		default:
			panic("bad packet")
		}
	}
}

func comparePackets(p0, p1 packetVal) int {
	n0, int0 := p0.(int)
	n1, int1 := p1.(int)
	if int0 && int1 {
		switch {
		case n0 < n1:
			return -1
		case n0 > n1:
			return 1
		default:
			return 0
		}
	}
	if int0 {
		return comparePackets([]packetVal{n0}, p1)
	}
	if int1 {
		return comparePackets(p0, []packetVal{n1})
	}
	l0 := p0.([]packetVal)
	l1 := p1.([]packetVal)
	for i, v0 := range l0 {
		if i >= len(l1) {
			return 1
		}
		if cmp := comparePackets(v0, l1[i]); cmp != 0 {
			return cmp
		}
	}
	if len(l1) > len(l0) {
		return -1
	}
	return 0
}
