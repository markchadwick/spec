package spec

import (
	"./github.com/markchadwick/assert" // TODO
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

func TestDeclarativeNestedSuite(t *testing.T) {
	parentCalled := 0
	child1Called := 0
	child2Called := 0

	s := Suite("Array", func(c *C) {
		parentCalled++
		arr := []string{}

		c.It("should have size zero", func(c *C) {
			child1Called++
			assert.That(t, arr).HasLen(0)
			arr = append(arr, "woah!")
			assert.That(t, arr).HasLen(1)
		})

		c.It("should still have size zero", func(c *C) {
			child2Called++
			assert.That(t, arr).HasLen(0)
			arr = append(arr, "woah!")
		})
	})

	assert.That(t, parentCalled).Equals(0)
	assert.That(t, child1Called).Equals(0)
	assert.That(t, child2Called).Equals(0)

	s.Run(nilReporter)
	t.Skip("TODO")
	assert.That(t, parentCalled).Equals(3)
	assert.That(t, child1Called).Equals(1)
	assert.That(t, child2Called).Equals(1)
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
	suite.run(nilReporter, false)
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
	suite.run(nilReporter, false)
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

func TestSuiteBeforeAfter(t *testing.T) {
	childCalled := false
	s := Suite("Before and After", func(c *C) {
		value := 1
		c.It("should execute in-order", func(c *C) {
			assert.That(t, value).Equals(1)
			childCalled = true
		})
		value++
	})
	s.Run(nilReporter)
	t.Skip("TODO")
	assert.That(t, childCalled).IsTrue()
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

func TestFailingChildrenShouldNotRun(t *testing.T) {
	t.Skip("Failing children should not run")
}

func TestReporterCalledOnStart(t *testing.T) {
	t.Skip("TODO")
}

func TestReporterCalledOnPass(t *testing.T) {
	t.Skip("TODO")
}

func TestReporterCalledWithErrsOnFail(t *testing.T) {
	t.Skip("TODO")
}
