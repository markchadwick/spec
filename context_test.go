package spec

import (
	"errors"
	"github.com/markchadwick/assert"
	"testing"
)

func TestFailCallStack(t *testing.T) {
	err := errors.New("Worst. Error. Ever.")

	Suite("Fail callstack", func(c *C) {
		c.Fail(err) // Nice work
	}).Run(nilReporter)

	assert.That(t, nilReporter.lastErrors).HasLen(1)
	e := nilReporter.lastErrors[0]

	assert.That(t, e.Err).Equals(err)
	assert.That(t, e.File != "").IsTrue()
	assert.That(t, e.Line > 0).IsTrue()
	assert.That(t, e.skip).IsFalse()

	src, err2 := e.Source()
	assert.That(t, err2).IsNil()
	assert.That(t, src).Equals("c.Fail(err) // Nice work")
}

func TestFailDirectCall(t *testing.T) {
	c := (&C{}).Fail(errors.New("Pants"))
	assert.That(t, c.errors).HasLen(1)

	src, _ := c.errors[0].Source()
	assert.That(t, src).Equals(`c := (&C{}).Fail(errors.New("Pants"))`)
}

func TestFailf(t *testing.T) {
	c := &C{}
	c.Failf("%s %d times", "count", 3)
	err := c.errors[0]

	assert.That(t, err.Error()).Equals("count 3 times")

	src, _ := err.Source()
	assert.That(t, src).Equals(`c.Failf("%s %d times", "count", 3)`)
}

func TestSkipPanic(t *testing.T) {
	skipCalled := false
	defer func() {
		assert.That(t, skipCalled).IsTrue()
	}()

	defer func() {
		if r := recover(); r != nil {
			if skip, ok := r.(*TestError); ok && skip.skip {
				skipCalled = true
				assert.That(t, skip.Err.Error()).Equals("what am I doing?")
				src, _ := skip.Source()
				assert.That(t, src).Equals(`c.Skip("what am I doing?")`)
				return
			}
			panic(r)
		}
	}()

	c := &C{}
	c.Skip("what am I doing?")
}
