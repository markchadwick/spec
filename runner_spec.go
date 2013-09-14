package spec

import (
	"github.com/markchadwick/assert"
	"testing"
)

func TestEmptyRunnerInitialization(t *testing.T) {
	assert.That(t, Runner().suites).
		NotNil().
		HasLen(0)
}

func TestRunnerInitialization(t *testing.T) {
	suite1 := Suite("Suite 1", func(c *C) {})
	assert.That(t, Runner(suite1).suites).HasLen(1)

	suite2 := Suite("Suite 2", func(c *C) {})
	assert.That(t, Runner(suite1, suite2).suites).HasLen(1)
}

func TestRunnerAddSuite(t *testing.T) {
	runner := Runner()
	assert.That(t, runner.suites).HasLen(0)

	runner.Add(Suite("Test 1", func(c *C) {}))
	assert.That(t, runner.suites).HasLen(1)

	runner.Add(Suite("Test 2", func(c *C) {}))
	assert.That(t, runner.suites).HasLen(2)
}

func TestDefaultRunner(t *testing.T) {
	t.Skip("TODO")
}
