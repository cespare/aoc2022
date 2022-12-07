package main

import (
	"fmt"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	scanner := ctx.scanner()
	root := new(fsDir)
	dir := root
outer:
	for scanner.scan() {
		line := scanner.text()
		if line == "$ ls" {
			continue
		}
		if d, ok := cutPrefix(line, "$ cd "); ok {
			if d == "/" {
				continue
			}
			if d == ".." {
				dir = dir.parent
				if dir == nil {
					panic(".. of root")
				}
				continue
			}
			for _, child := range dir.dirs {
				if child.name == d {
					dir = child
					continue outer
				}
			}
			panic(fmt.Sprintf("no child %s", d))
		}
		if d, ok := cutPrefix(line, "dir "); ok {
			dir.dirs = append(dir.dirs, &fsDir{name: d, parent: dir})
			continue
		}
		var file fsFile
		if _, err := fmt.Sscanf(line, "%d %s", &file.size, &file.name); err != nil {
			panic(err)
		}
		dir.files = append(dir.files, &file)
	}
	root.fillSizes()
	ctx.reportLoad()

	ctx.reportPart1(root.sizeLimit(100000))

	free := 70000000 - root.size
	need := 30000000 - free
	ctx.reportPart2(root.nearest(need))
}

type fsFile struct {
	name string
	size int64
}

type fsDir struct {
	name   string
	parent *fsDir
	dirs   []*fsDir
	files  []*fsFile
	size   int64
}

func (d *fsDir) fillSizes() {
	for _, child := range d.dirs {
		child.fillSizes()
		d.size += child.size
	}
	for _, f := range d.files {
		d.size += f.size
	}
}

func (d *fsDir) sizeLimit(lim int64) int64 {
	var total int64
	if d.size <= lim {
		total += d.size
	}
	for _, child := range d.dirs {
		total += child.sizeLimit(lim)
	}
	return total
}

func (d *fsDir) nearest(thresh int64) int64 {
	r := int64(-1)
	if d.size >= thresh {
		r = d.size
	}
	for _, child := range d.dirs {
		rc := child.nearest(thresh)
		if rc > 0 && rc < r {
			r = rc
		}
	}
	return r
}
