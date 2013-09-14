package spec

import (
	"testing"
	"time"
)

var syntax = Suite("Spec syntax", func(c *C) {

	c.It("Should have a first child", func(c *C) {
	})

	c.It("Should have a second child", func(c *C) {
	})

	c.It("May take a while!", func(c *C) {
		time.Sleep(500 * time.Millisecond)
	})

})

func TestSuite(t *testing.T) {
	Runner(syntax).Run(AsciiReporter())
}
