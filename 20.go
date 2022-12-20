package main

import (
	"golang.org/x/exp/slices"
)

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var nums []int
	scanner := ctx.scanner()
	for scanner.scan() {
		nums = append(nums, int(parseInt(scanner.text())))
	}
	ctx.reportLoad()

	f := newEncFile(nums)
	f.mix()
	ctx.reportPart1(f.answer())

	f = newEncFile(SliceMap(nums, func(n int) int { return n * 811589153 }))
	for i := 0; i < 10; i++ {
		f.mix()
	}
	ctx.reportPart2(f.answer())
}

type encFile struct {
	nodes []*encFileNode // in original order
	head  *encFileNode
	tail  *encFileNode
}

type encFileNode struct {
	n    int
	next *encFileNode
	prev *encFileNode
}

func newEncFile(nums []int) *encFile {
	f := &encFile{nodes: make([]*encFileNode, len(nums))}
	for i, n := range nums {
		node := &encFileNode{n: n}
		f.nodes[i] = node
		if f.head == nil {
			f.head = node
		}
		if f.tail != nil {
			f.tail.next = node
			node.prev = f.tail
		}
		f.tail = node
	}
	return f
}

func (f *encFile) slice() []int {
	nums := make([]int, 0, len(f.nodes))
	for node := f.head; node != nil; node = node.next {
		nums = append(nums, node.n)
	}
	return nums
}

func (f *encFile) answer() int {
	s := f.slice()
	iz := slices.Index(s, 0)
	if iz < 0 {
		panic("zero gone")
	}
	get := func(afterZero int) int {
		i := (iz + afterZero) % len(f.nodes)
		return s[i]
	}
	return get(1000) + get(2000) + get(3000)
}

func (f *encFile) mix() {
	for _, node := range f.nodes {
		f.mixAt(node)
	}
}

func (f *encFile) mixAt(node *encFileNode) {
	lim := abs(node.n) % (len(f.nodes) - 1)
	for i := 0; i < lim; i++ {
		if node.n < 0 {
			f.moveLeft(node)
		} else {
			f.moveRight(node)
		}
	}
}

// All methods below assume a reasonable list length (>3, at least) for
// simplicity.

func (f *encFile) moveLeft(node *encFileNode) {
	f.delete(node)
	switch {
	case node.prev == nil:
		f.insertAfter(node, f.tail.prev)
	case node.prev.prev == nil:
		f.insertAfter(node, f.tail)
	default:
		f.insertAfter(node, node.prev.prev)
	}
}

func (f *encFile) moveRight(node *encFileNode) {
	f.delete(node)
	switch {
	case node.next == nil:
		f.insertAfter(node, f.head)
	case node.next.next == nil:
		f.insertHead(node)
	default:
		f.insertAfter(node, node.next)
	}
}

func (f *encFile) delete(node *encFileNode) {
	if node.prev == nil {
		f.head = node.next
	} else {
		node.prev.next = node.next
	}
	if node.next == nil {
		f.tail = node.prev
	} else {
		node.next.prev = node.prev
	}
}

func (f *encFile) insertHead(node *encFileNode) {
	node.prev = nil
	node.next = f.head
	f.head.prev = node
	f.head = node
}

func (f *encFile) insertAfter(node, left *encFileNode) {
	if left.next == nil {
		f.tail = node
	} else {
		left.next.prev = node
	}
	node.next = left.next
	left.next = node
	node.prev = left
}
