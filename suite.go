package spec

import (
	"fmt"
	"testing"
	"time"
)

func Run(t *testing.T) {
	if err := DefaultRunner.Run(Console()); err != nil {
		t.Fatal(err)
	}
}

type stats struct {
	Duration time.Duration
}

type suite struct {
	Name     string
	Test     Test
	Stats    *stats
	children []*suite
	ctx      *C
}

type Test func(c *C)

// Create a new suite with the given name and test body.
func Suite(name string, test Test) *suite {
	suite := &suite{
		Name:  name,
		Test:  test,
		Stats: &stats{},
	}
	suite.ctx = &C{
		suite: suite,
	}
	DefaultRunner.Add(suite)
	return suite
}

// Run a suite and all of its children. If a suite has no children, it will be
// run exactly once. Otherwise, it will be run before each of its children.
func (s *suite) Run(reporter Reporter) {
	reporter.Start(s)

	start := time.Now()
	errs, skip := s.run(reporter)
	s.Stats.Duration = time.Now().Sub(start)

	if skip != nil {
		reporter.Skip(s, skip)
	} else if len(errs) != 0 {
		reporter.Fail(s, errs)
	} else {
		reporter.Pass(s)
		reporter.Descend(s)
		for _, child := range s.children {
			s.runChild(child, reporter)
		}
		reporter.Ascend(s)
	}
}

// Run a test without descending into its children. If the children of this test
// have not been seen before, they will be collected as the test runs through
// its first go. Otherwise, they will be ignored (a tests children should be
// considered immutable).
func (s *suite) run(reporter Reporter) (errs []*TestError, skip *TestError) {
	skip = &TestError{}

	// Capture Skip tests, which panic to halt execution of the test. TODO: Also
	// fatal tests.
	defer func() {
		if r := recover(); r != nil {
			if s, ok := r.(*TestError); ok && s.skip {
				skip = s
				return
			}
			panic(r)
		}
	}()

	if s.children == nil {
		s.children = make([]*suite, 0)
		s.ctx.onChild = func(child *suite) {
			s.children = append(s.children, child)
		}
		defer func() { s.ctx.onChild = nil }()
	}

	s.ctx.errors = make([]*TestError, 0)
	defer func() { s.ctx.errors = nil }()

	s.Test(s.ctx)
	return s.ctx.errors, nil
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
	s.run(reporter)

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
