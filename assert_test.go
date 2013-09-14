package spec

import (
	"testing"
)

var assertSpec = Suite("Suite assertions", func(c *C) {

	c.It("Should integrate with assert checks", func(c *C) {
		c.Assert(nil).
			IsNil().
			HasLen(52).
			NotNil()
	})
})

func TestAssert(t *testing.T) {
	Runner(assertSpec).Run(Console())
}
