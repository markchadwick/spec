package spec

import (
	"fmt"
)

type suite struct {
	Name     string
	Test     Test
	children []*suite
	ctx      *C
}

type Test func(c *C)

// Create a new suite with the given name and test body.
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

// Run a suite and all of its children. If a suite has no children, it will be
// run exactly once. Otherwise, it will be run before each of its children.
func (s *suite) Run(reporter Reporter) {
	errs := s.run(reporter, true)

	if len(errs) != 0 {
		reporter.Fail(s, errs)
	} else {
		for _, child := range s.children {
			s.runChild(child, reporter)
		}
		reporter.Pass(s)
	}
}

// Run a test without descending into its children. If the children of this test
// have not been seen before, they will be collected as the test runs through
// its first go. Otherwise, they will be ignored (a tests children should be
// considered immutable).
func (s *suite) run(reporter Reporter, report bool) []error {
	if s.children == nil {
		s.children = make([]*suite, 0)
		s.ctx.onChild = func(child *suite) {
			s.children = append(s.children, child)
		}
		defer func() { s.ctx.onChild = nil }()
	}

	if report {
		reporter.Start(s)
	}
	s.Test(s.ctx)
	return nil
}

// Run this suite only descending into the given child. It should be run in
// order such that the preamble to a child will be executed first, then the body
// of the child, then the postamble of the parent.
//
// If this finds multiple children that look "equal" (see `equals` below), it
// will only run the first.
func (s *suite) runChild(c *suite, reporter Reporter) error {
	childRan := false
	s.ctx.onChild = func(child *suite) {
		if child.equals(c) {
			child.Run(reporter)
			childRan = true
		}
	}
	defer func() { s.ctx.onChild = nil }()
	s.run(reporter, false)

	if !childRan {
		return fmt.Errorf("Found no child named '%s'", c.Name)
	}
	return nil
}

// Super-hacky approximation of if two suites are equal. For the moment, this
// will only compare their names, which is obviously not going to work in all
// cases.
// TODO: Figure out how to accurately compare suites
func (s *suite) equals(s1 *suite) bool {
	if s == nil || s1 == nil {
		return false
	}
	return s.Name == s1.Name
}
