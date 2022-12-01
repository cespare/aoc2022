package main

import "golang.org/x/exp/slices"

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	var elves [][]int64
	scanner := ctx.scanner()
	var elf []int64
	for scanner.scan() {
		if scanner.text() != "" {
			elf = append(elf, scanner.int64())
			continue
		}
		if len(elf) > 0 {
			elves = append(elves, elf)
			elf = nil
		}
	}
	if len(elf) > 0 {
		elves = append(elves, elf)
	}
	ctx.reportLoad()

	calories := SliceMap(elves, SliceSum[[]int64, int64])
	part1 := SliceReduce(calories, 0, maxim[int64])
	ctx.reportPart1(part1)

	slices.SortFunc(calories, func(x, y int64) bool { return x > y })
	part2 := SliceSum(calories[:3])
	ctx.reportPart2(part2)
}
