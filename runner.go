package spec

import (
	"fmt"
	"sync"
)

var DefaultRunner = Runner()

// TODO:
//  - Flags
//    -spec.f filename
type runner struct {
	suites    []*suite
	reporters []Reporter
	runLock   *sync.Mutex
	errors    []*SuiteFailure
}

func Runner(suites ...*suite) *runner {
	return &runner{
		suites:  suites,
		runLock: new(sync.Mutex),
	}
}

// Adds a suite to this runner
func (r *runner) Add(s *suite) {
	r.suites = append(r.suites, s)
}

// Run this suite reporting test conditions to each of the given reporters.
// Tests will only be run once, and their results broadcast to each reporter.
func (r *runner) Run(reporters ...Reporter) error {
	r.runLock.Lock()
	defer r.runLock.Unlock()

	r.reporters = reporters
	r.errors = make([]*SuiteFailure, 0)

	r.Begin()

	for _, suite := range r.suites {
		suite.Run(r)
	}

	r.Finish(r.errors)

	if len(r.errors) == 0 {
		return nil
	}
	return fmt.Errorf("%d test failures", len(r.errors))
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

func (r *runner) Fail(s *suite, errs []*TestError) {
	r.errors = append(r.errors, &SuiteFailure{s, errs})
	for _, r := range r.reporters {
		r.Fail(s, errs)
	}
}

func (r *runner) Skip(s *suite, skip *TestError) {
	for _, r := range r.reporters {
		r.Skip(s, skip)
	}
}

func (r *runner) Descend(s *suite) {
	for _, r := range r.reporters {
		r.Descend(s)
	}
}

func (r *runner) Ascend(s *suite) {
	for _, r := range r.reporters {
		r.Ascend(s)
	}
}

func (r *runner) Begin() {
	for _, r := range r.reporters {
		r.Begin()
	}
}

func (r *runner) Finish(e []*SuiteFailure) {
	for _, r := range r.reporters {
		r.Finish(e)
	}
}
