package spec

import (
	"fmt"
	"github.com/markchadwick/assert"
	"testing"
)

// func TestSuiteString(t *testing.T) {
// 	parent := Suite("Parent", func(c *C) {})
// 	assert.That(t, parent.String()).Equals("Parent")
//
// 	child := Suite("Child", func(c *C) {})
// 	assert.That(t, child.String()).Equals("Child")
// 	child.parent = parent
// 	assert.That(t, child.String()).Equals("Parent -> Child")
// }
//
// func TestBasicSuiteConstruction(t *testing.T) {
// 	calls := 0
// 	suite := Suite("Syntax", func(c *C) {
// 		calls++
// 	})
//
// 	assert.That(t, suite).NotNil()
// 	assert.That(t, calls).Equals(0)
//
// 	suite.Run()
// 	assert.That(t, calls).Equals(1)
// }
//
// func TestBasicNestedSuite(t *testing.T) {
// 	parentCalled := 0
// 	child1Called := 0
// 	child2Called := 0
//
// 	parent := Suite("Parent", func(c *C) {
// 		parentCalled++
// 	})
//
// 	parent.children = []*suite{
// 		Suite("Child 1", func(c *C) { child1Called++ }),
// 		Suite("Child 2", func(c *C) { child2Called++ }),
// 	}
//
//   assert.That(t, parentCalled).Equals(0)
//   assert.That(t, child1Called).Equals(0)
//   assert.That(t, child2Called).Equals(0)
//
//   parent.Run()
//   assert.That(t, parentCalled).Equals(2)
//   assert.That(t, child1Called).Equals(1)
//   assert.That(t, child2Called).Equals(1)
// }

func TestDeclarativeNestedSuite(t *testing.T) {
	parentCalled := 0
	child1Called := 0
	child2Called := 0

	s := Suite("Array", func(c *C) {
		fmt.Println("* calling parent")
		parentCalled++
		arr := []string{}

		c.It("should have size zero", func(c *C) {
			fmt.Println("  * calling child 1")
			child1Called++
			assert.That(t, arr).HasLen(0)
			arr = append(arr, "woah!")
			assert.That(t, arr).HasLen(1)
		})

		c.It("should still have size zero", func(c *C) {
			fmt.Println("  * calling child 2")
			child2Called++
			assert.That(t, arr).HasLen(0)
			arr = append(arr, "woah!")
		})
	})

	assert.That(t, parentCalled).Equals(0)
	assert.That(t, child1Called).Equals(0)
	assert.That(t, child2Called).Equals(0)

	s.Run()
	assert.That(t, parentCalled).Equals(2)
	assert.That(t, child1Called).Equals(1)
	assert.That(t, child2Called).Equals(1)
}
