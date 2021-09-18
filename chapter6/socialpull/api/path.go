package main

import (
	"strings"
)

const PathSeparater = "/"

type Path struct {
	Path string
	ID string
}

// NewPath separates last part and others, then returns them as Path.
func NewPath(p string) *Path {
	var id string
	p = strings.Trim(p, PathSeparater)
	s := strings.Split(p, PathSeparater)
	if len(s) > 1 {
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparater)
	}
	return &Path{Path: p, ID: id}
}

// HasID returns true when the path contains ID.
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}

