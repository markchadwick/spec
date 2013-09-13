package spec

import (
	"fmt"
	"sort"
	"strings"
)

type suite struct {
	Name     string
	Test     Test
	parent   *suite
	children []*suite
	ctx      *C
}

type Test func(c *C)

func Suite(name string, test Test) *suite {
	suite := &suite{
		Name: name,
		Test: test,
	}
	suite.ctx = &C{
		suite: suite,
	}
	return suite
}

func (s *suite) Add(child *suite) {
	child.parent = s
	s.children = append(s.children, child)
	fmt.Printf("[Suite %s] children now %d\n", s.Name, len(s.children))
}

// Run a suite and all of its children. If a suite has no children, it will be
// run exactly once. Otherwise, it will be run before each of its children.
func (s *suite) Run() {
	if len(s.children) == 0 {
		s.Test(s.ctx)
	} else {
		for _, child := range s.children {
			s.Test(s.ctx)
			child.Run()
		}
	}
}

// Simple string represensation of a suite that traces its lineage back to the
// root parent.
func (s *suite) String() string {
	suite := s
	lineage := make([]string, 0)
	for suite != nil {
		lineage = append(lineage, suite.Name)
		suite = suite.parent
	}
	sort.Sort(sort.Reverse(sort.StringSlice(lineage)))
	return strings.Join(lineage, " -> ")
}
