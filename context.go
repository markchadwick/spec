package spec

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type onChild func(*suite)

type C struct {
	suite   *suite
	onChild onChild
	errors  []*TestError
}

func (c *C) It(name string, test Test) *suite {
	s := Suite(name, test)
	if c.onChild != nil {
		c.onChild(s)
	}
	return s
}

func (c *C) Fail(err error) *C {
	return c.fail(err, 1)
}

func (c *C) Failf(f string, args ...interface{}) *C {
	return c.fail(fmt.Errorf(f, args...), 1)
}

func (c *C) Skip(msg string, args ...interface{}) *C {
	err := &TestError{
		Err:  fmt.Errorf(msg, args...),
		skip: true,
	}
	err.inspect(1)
	panic(err)
}

func (c *C) fail(err error, depth int) *C {
	testError := &TestError{Err: err}
	testError.inspect(depth + 1)

	c.errors = append(c.errors, testError)
	return c
}

// ----------------------------------------------------------------------------
// Test Error
// ----------------------------------------------------------------------------

type TestError struct {
	Err  error
	File string
	Line int
	skip bool
}

func (t *TestError) Error() string {
	return t.Err.Error()
}

func (t *TestError) Source() (string, error) {
	line, err := readLine(t.File, t.Line)
	line = strings.TrimSpace(line)
	return line, err
}

// Try to determine what the failing line of code is. It must be in the call
// stack when this is called.
func (t *TestError) inspect(depth int) {
	_, t.File, t.Line, _ = runtime.Caller(depth + 1)
}

// Read a single line of a file.
func readLine(fname string, lineNo int) (string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 1
	for ; scanner.Scan(); i++ {
		if i == lineNo {
			return scanner.Text(), nil
		}
	}
	return "", scanner.Err()
}
