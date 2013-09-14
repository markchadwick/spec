package spec

import (
	"testing"
	"time"
)

var _ = Suite("Spec syntax", func(c *C) {

	c.It("Should have a first child", func(c *C) {
	})

	c.It("May take a while!", func(c *C) {
		c.Skip("ok")
		time.Sleep(500 * time.Millisecond)
	})

	c.It("Should have a second child", func(c *C) {
	})

	c.It("Should fail", func(c *C) {
		c.Failf("errrp!")
	})

	c.It("Should skip a test", func(c *C) {
		c.Skip("not now...")
	})

})

func TestSuite(t *testing.T) {
	t.Skip("Not now")
	Run(t)
}
