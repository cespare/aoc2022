package main

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	nums := scanSlice(ctx, func(line string) int { return int(parseInt(line)) })
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
}

type encFileNode struct {
	n    int
	next *encFileNode
	prev *encFileNode
}

func newEncFile(nums []int) *encFile {
	f := &encFile{nodes: make([]*encFileNode, len(nums))}
	var prev *encFileNode
	for i, n := range nums {
		node := &encFileNode{n: n}
		f.nodes[i] = node
		node.prev = prev
		if prev != nil {
			prev.next = node
		}
		prev = node
	}
	f.nodes[0].prev = f.nodes[len(f.nodes)-1]
	f.nodes[len(f.nodes)-1].next = f.nodes[0]
	return f
}

func (f *encFile) answer() int {
	var zero *encFileNode
	for _, node := range f.nodes {
		if node.n == 0 {
			zero = node
			break
		}
	}

	var i int
	ordered := make([]int, len(f.nodes))
	for node := zero; i == 0 || node != zero; node = node.next {
		ordered[i] = node.n
		i++
	}
	get := func(afterZero int) int {
		i := afterZero % len(ordered)
		return ordered[i]
	}
	return get(1000) + get(2000) + get(3000)
}

func (f *encFile) mix() {
	for _, node := range f.nodes {
		f.mixAt(node)
	}
}

func (f *encFile) mixAt(node *encFileNode) {
	node.remove()
	lim := abs(node.n) % (len(f.nodes) - 1)
	node1 := node
	if node.n <= 0 {
		for i := 0; i < lim+1; i++ {
			node1 = node1.prev
		}
	} else {
		for i := 0; i < lim; i++ {
			node1 = node1.next
		}
	}
	node.insertAfter(node1)
}

func (n *encFileNode) remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (n *encFileNode) insertAfter(left *encFileNode) {
	left.next.prev = n
	n.next = left.next
	left.next = n
	n.prev = left
}
