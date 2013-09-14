package spec

import (
	"sync"
)

// TODO Docs
var DefaultRunner = Runner()

// TODO:
//  Flags
//  -----
//    -spec.f filename
type runner struct {
	suites    []*suite
	reporters []Reporter
	runLock   *sync.Mutex
}

func Runner(suites ...*suite) *runner {
	return &runner{
		suites:  suites,
		runLock: new(sync.Mutex),
	}
}

// Adds a suite to this runner
func (r *runner) Add(s *suite) {
	DefaultRunner.Add(s)
	r.suites = append(r.suites, s)
}

// Run this suite reporting test conditions to each of the given reporters.
// Tests will only be run once, and their results broadcast to each reporter.
func (r *runner) Run(reporters ...Reporter) {
	r.runLock.Lock()
	defer r.runLock.Unlock()

	r.reporters = reporters
	defer func() { r.reporters = nil }()

	for _, suite := range r.suites {
		suite.Run(r)
	}
}

func (r *runner) Start(s *suite) {
	for _, r := range r.reporters {
		r.Start(s)
	}
}

func (r *runner) Pass(s *suite) {
	for _, r := range r.reporters {
		r.Pass(s)
	}
}

func (r *runner) Fail(s *suite, errs []error) {
	for _, r := range r.reporters {
		r.Fail(s, errs)
	}
}
