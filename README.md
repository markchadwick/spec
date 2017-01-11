# spec
Test runner for Go tests.

## Usage
Assertions defined by the [assert](https://github.com/markchadwick/assert)
project can be used to check values.

```go
import (
  "github.com/markchadwick/spec"
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

// Run all the suites
func Test(t *testing.T) {
  spec.Run(t)
}
```

## jUnit reporting
A basic [jUnit](http://junit.org/) runner is built in. To drop a jUnit-formatted
XMl file in your root while your tests run, you can bind to the testing
framework as follows:

```go
import (
  "github.com/markchadwick/spec"
  "os"
  "testing"
)

func TestSpecs(t *testing.T) {
  out, err := os.Create("spec.xml")
  if err != nil {
    t.Fatal(err)
  }
  defer out.Close()

  if err = spec.DefaultRunner.Run(spec.Console(), spec.JUnit(out)); err != nil {
    t.Fatal(err)
  }
}

```

## Execution order
Tests run in the order in which they are declared.

For each test, it can be assumed that the body of the parent suite before the
corrent suite has been run in isolation and the body of the parent will run when
the current is complete. For suites with just one level of nesting, this is
obvious, but it becomes fuzzier as you move down.

In general, it is safe to assume that no other suites have messed with what your
parent declares, but its parent may be fair game. This is largely a matter of
speeding up tests focused on local changes, but is a departure from how, say,
[mocha](http://mochajs.org/) works.

This spec illustrates the matter.

```go
spec.Suite("Execution order", func(c *spec.C) {
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
```
