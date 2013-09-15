package spec

import (
	"bytes"
	"log"
	"os"
	"testing"
)

var emptySuite = Suite("Empty JUnit Suite", func(c *C) {})

var jUnitSpec = Suite("JUnit spec output", func(c *C) {

	c.It("should sanitize class names", func(c *C) {
		junit := JUnit(nil)

		c.It("valid class names", func(c *C) {
			c.Assert(junit.className("ClassName")).Equals("ClassName")
		})

		c.It("names with spaces", func(c *C) {
			c.Assert(junit.className("Class Name")).Equals("ClassName")
		})

		c.It("names with numbers", func(c *C) {
			c.Assert(junit.className("Class9 Name")).Equals("Class9Name")
		})

		c.It("names starting with numbers", func(c *C) {
			c.Skip("figure out the expr")
			c.Assert(junit.className("6ClassName")).Equals("ClassName")
		})

		c.It("names with symbols", func(c *C) {
			c.Assert(junit.className(" Class! -Name_  ")).Equals("ClassName")
		})
	})

	c.It("In an empty suite", func(c *C) {
		buf := new(bytes.Buffer)
		junit := JUnit(buf)

		log.Printf("Junit: %v", junit)
	})

	c.It("fails", func(c *C) {
		c.Failf("nope")
	})
})

func TestJunit(t *testing.T) {
	t.Skip("Visual inspection time")
	Runner(emptySuite, jUnitSpec).Run(JUnit(os.Stdout))
}
