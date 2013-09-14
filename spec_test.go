package spec

import (
	"./github.com/markchadwick/assert" // TODO
	"fmt"
	"testing"
)

func TestBasicSuiteConstruction(t *testing.T) {
	calls := 0
	suite := Suite("Syntax", func(c *C) {
		calls++
	})

	assert.That(t, suite).NotNil()
	assert.That(t, calls).Equals(0)

	suite.Run(nilReporter)
	assert.That(t, calls).Equals(1)
}

func TestSuiteChildren(t *testing.T) {
	suiteRuns := 0
	child1Runs := 0
	child2Runs := 0

	suite := Suite("Children suite", func(c *C) {
		suiteRuns++
		c.It("child 1", func(c *C) { child1Runs++ })
		c.It("child 2", func(c *C) { child2Runs++ })
	})

	assert.That(t, suite.children).IsNil()
	suite.run(nilReporter)
	assert.That(t, suite.children).NotNil()

	assert.That(t, suiteRuns).Equals(1)
	assert.That(t, child1Runs).Equals(0)
	assert.That(t, child2Runs).Equals(0)

	children := suite.children
	assert.That(t, children).HasLen(2)
	assert.That(t, children[0].Name).Equals("child 1")
	assert.That(t, children[1].Name).Equals("child 2")
}

func TestRunChild(t *testing.T) {
	suiteRuns := 0
	child1Runs := 0
	child2Runs := 0

	var child1, child2 *suite

	suite := Suite("Children suite", func(c *C) {
		suiteRuns++
		child1 = c.It("child 1", func(c *C) { child1Runs++ })
		child2 = c.It("child 2", func(c *C) { child2Runs++ })
	})
	suite.run(nilReporter)
	assert.That(t, suiteRuns).Equals(1)
	assert.That(t, child1Runs).Equals(0)
	assert.That(t, child2Runs).Equals(0)

	assert.That(t, suite.runChild(child1, nilReporter)).IsNil()
	assert.That(t, suiteRuns).Equals(2)
	assert.That(t, child1Runs).Equals(1)
	assert.That(t, child2Runs).Equals(0)

	assert.That(t, suite.runChild(child2, nilReporter)).IsNil()
	assert.That(t, suiteRuns).Equals(3)
	assert.That(t, child1Runs).Equals(1)
	assert.That(t, child2Runs).Equals(1)

	nonsense := Suite("nonsense", func(c *C) {})
	err := suite.runChild(nonsense, nilReporter)
	assert.That(t, err).NotNil()

	assert.That(t, suiteRuns).Equals(4)
	assert.That(t, child1Runs).Equals(1)
	assert.That(t, child2Runs).Equals(1)
}

func TestEqualChildren(t *testing.T) {
	suiteRuns := 0
	child1Runs := 0
	child2Runs := 0

	var child1, child2 *suite

	Suite("Equal children", func(c *C) {
		suiteRuns++
		child1 = c.It("child", func(c *C) { child1Runs++ })
		child2 = c.It("child", func(c *C) { child2Runs++ })
	}).Run(nilReporter)

	t.Skip("Can't differentiate between children with the same name")
	assert.That(t, suiteRuns).Equals(3)
	assert.That(t, child1Runs).Equals(1)
	assert.That(t, child2Runs).Equals(1)
}

func TestTestErrors(t *testing.T) {
	Suite("Suite errors", func(c *C) {
		c.Fail(fmt.Errorf("First error"))
		c.Fail(fmt.Errorf("Second error"))
	}).Run(nilReporter)

	errs := nilReporter.lastErrors
	assert.That(t, errs).HasLen(2)
	assert.That(t, errs[0].Error()).Equals("First error")
	assert.That(t, errs[1].Error()).Equals("Second error")
}

func TestSuiteSkipHalts(t *testing.T) {
	started := false
	finished := false
	Suite("Skip should halt", func(c *C) {
		started = true
		c.Skip("Changed my mind!")
		finished = true
	}).Run(nilReporter)

	assert.That(t, started).IsTrue()
	assert.That(t, finished).IsFalse()

	src, _ := nilReporter.lastSkip.Source()
	assert.That(t, src).Equals(`c.Skip("Changed my mind!")`)
}
