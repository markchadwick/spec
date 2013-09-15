package doc

import (
	spec ".."
	"errors"
	"testing"
)

var _ = spec.Suite("An Array", func(c *spec.C) {
	array := []string{}

	c.It("should be initialized with 0 length", func(c *spec.C) {
		c.Assert(array).
			NotNil().
			HasLen(0)
	})

	c.It("should hold a name", func(c *spec.C) {
		array = append(array, "bob")
		c.Assert(array).HasLen(1)
	})

	c.It("should sort itself", func(c *spec.C) {
		c.Skip("I don't know how to do this yet")
	})
})

var _ = spec.Suite("Failing suite", func(c *spec.C) {
	c.It("can fail with the assert syntax", func(c *spec.C) {
		c.Assert(nil).NotNil()
	})

	c.It("can fail with an error", func(c *spec.C) {
		err := errors.New("Creative!")
		c.Fail(err)
	})

	c.It("can fail with a message", func(c *spec.C) {
		c.Failf("Not happening %s", "today")
	})
})

var _ = spec.Suite("Execution order", func(c *spec.C) {
	grandparent := 0

	c.It("should reset for the first child", func(c *spec.C) {
		grandparent++
		c.Assert(grandparent).Equals(1)
	})

	c.It("should reset for the second child", func(c *spec.C) {
		grandparent++
		c.Assert(grandparent).Equals(1)
	})

	c.It("is no longer the direct parent", func(c *spec.C) {
		parent := 0

		c.It("should reset the first child", func(c *spec.C) {
			grandparent++
			parent++

			c.Assert(grandparent).Equals(1)
			c.Assert(parent).Equals(1)
		})

		c.It("should reset the second child", func(c *spec.C) {
			grandparent++
			parent++

			// HEY! Grandparent is now two because the test above has modified above
			// its parent scope!
			c.Assert(grandparent).Equals(2)
			c.Assert(parent).Equals(1)
		})
	})
})

func Test(t *testing.T) {
	spec.Run(t)
}
