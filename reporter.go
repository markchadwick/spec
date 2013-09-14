package spec

import (
	"encoding/xml"
	"fmt"
	"github.com/mgutz/ansi"
	"io"
	"regexp"
	"strings"
	"time"
)

type SuiteFailure struct {
	suite  *suite
	errors []*TestError
}

type Reporter interface {
	Start(*suite)
	Pass(*suite)
	Fail(*suite, []*TestError)
	Skip(*suite, *TestError)
	Descend(*suite)
	Ascend(*suite)
	Begin()
	Finish([]*SuiteFailure)
}

// ----------------------------------------------------------------------------
// Console Reporter
// ----------------------------------------------------------------------------

type ConsoleReporter struct {
	depth int

	numSpec int
	numPass int
	numFail int
	numSkip int

	start time.Time
}

var (
	fail     = "\xE2\x98\xA0"
	iconPass = ansi.Color("\xE2\x9c\x93", "green")
	iconFail = ansi.Color(fail, "white+b")
)

func Console() *ConsoleReporter {
	return &ConsoleReporter{}
}

func (c *ConsoleReporter) Start(s *suite) {
	c.numSpec++
	c.status(" ", s.Name, nil)
}

func (c *ConsoleReporter) Pass(s *suite) {
	c.numPass++

	c.status(iconPass, s.Name, &s.Stats.Duration)
	fmt.Println()
}

func (c *ConsoleReporter) Fail(s *suite, errs []*TestError) {
	c.numFail++

	name := ansi.Color(s.Name, "red")
	c.status(iconFail, name, &s.Stats.Duration)
	fmt.Println()
}

func (c *ConsoleReporter) Skip(s *suite, skip *TestError) {
	c.numSkip++

	reason := skip.Error()
	msg := fmt.Sprintf("%s %s", ansi.Color(s.Name, "yellow"), reason)
	c.status(" ", msg, &s.Stats.Duration)
	fmt.Println()
}

func (c *ConsoleReporter) Descend(*suite) {
	c.depth++
}

func (c *ConsoleReporter) Ascend(*suite) {
	c.depth--
}

func (c *ConsoleReporter) Begin() {
	c.start = time.Now()
	fmt.Println()
}

func (c *ConsoleReporter) Finish(errs []*SuiteFailure) {
	duration := time.Now().Sub(c.start)
	fmt.Printf("\n\n----------------------------------------------------\n")
	fmt.Printf("%d PASSED %d FAILED %d SKIPPED\n", c.numPass, c.numFail, c.numSkip)

	for _, err := range errs {
		c.printSuiteFailure(err)
	}

	var status string
	if len(errs) == 0 {
		status = ansi.Color("OK", "green")
	} else {
		status = ansi.Color("FAIL", "red")
	}

	fmt.Printf("%s (%d specs in %s)\n", status, c.numSpec, duration)
}

func (c *ConsoleReporter) status(icon, msg string, duration *time.Duration) {
	fmt.Print("\r")
	c.pad()

	dur := ""
	if duration != nil {
		dur = duration.String()
	}

	fmt.Printf("%s %-10s %s", icon, msg, ansi.Color(dur, "+h"))
}

func (c *ConsoleReporter) statusf(msg string, args ...interface{}) {
	fmt.Print("\r")
	c.pad()
	fmt.Printf(msg, args...)
}

func (c *ConsoleReporter) printSuiteFailure(err *SuiteFailure) {
	fmt.Println()
	fmt.Printf("  FAILURE in '%s'\n", err.suite.Name)

	for _, err := range err.errors {
		fmt.Printf("  %s %s:%d\n", ansi.Color(fail, "red+b"), err.File, err.Line)
		if src, e := err.Source(); e == nil {
			fmt.Printf("    %s\n", ansi.Color(src, "white+b"))
		}

		for _, msg := range strings.Split(err.Error(), "\n") {
			fmt.Printf("      %s\n", msg)
		}
		fmt.Println()
	}
}

func (c *ConsoleReporter) pad() int {
	pad := "  "
	for i := 0; i < c.depth; i++ {
		fmt.Print(pad)
	}
	return c.depth * len(pad)
}

// ----------------------------------------------------------------------------
// JUnit Reporter
// ----------------------------------------------------------------------------

type JunitReporter struct {
	w            io.Writer
	stack        []string
	enc          *xml.Encoder
	invalidName  *regexp.Regexp
	invalidClass *regexp.Regexp
	suite        *testsuite
	current      *testcase
}

type testsuite struct {
	Tests int         `xml:"tests,attr"`
	Cases []*testcase `xml:"testcase"`
}

func (t *testsuite) Add(c *testcase) {
	t.Cases = append(t.Cases, c)
}

type skipped struct {
	Message string `xml:"message,attr"`
}

type testcase struct {
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Failures  []string `xml:"failure"`
	Skipped   *skipped `xml:"skipped,omitempty"`
}

func (t *testcase) Fail(msg string) {
	t.Failures = append(t.Failures, msg)
}

func JUnit(w io.Writer) *JunitReporter {
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")

	suite := &testsuite{
		Cases: make([]*testcase, 0),
	}

	return &JunitReporter{
		w:            w,
		stack:        make([]string, 0),
		enc:          enc,
		invalidClass: regexp.MustCompile(`[^A-Za-z0-9]`),
		invalidName:  regexp.MustCompile(`[^A-Za-z0-9 ]`),
		suite:        suite,
	}
}

func (j *JunitReporter) Start(s *suite) {
	j.suite.Tests++
	stack := append([]string{"test"}, j.stack...)
	j.current = &testcase{
		Classname: strings.Join(stack, "."),
		Name:      j.name(s.Name),
		Failures:  make([]string, 0),
	}
}

func (j *JunitReporter) Pass(*suite) {
	j.suite.Add(j.current)
	j.current = nil
}

func (j *JunitReporter) Fail(s *suite, errs []*TestError) {
	for _, err := range errs {
		j.current.Fail(err.Error())
	}
	j.suite.Add(j.current)
	j.current = nil
}

func (j *JunitReporter) Skip(s *suite, e *TestError) {
	j.current.Skipped = &skipped{e.Error()}
	j.suite.Add(j.current)
	j.current = nil
}

func (j *JunitReporter) Descend(s *suite) {
	j.stack = append(j.stack, j.className(s.Name))
}

func (j *JunitReporter) Ascend(*suite) {
	j.stack = j.stack[0 : len(j.stack)-1]
}

func (j *JunitReporter) Begin() {
}

func (j *JunitReporter) Finish([]*SuiteFailure) {
	fmt.Fprint(j.w, xml.Header)
	j.enc.Encode(j.suite)
}

func (j *JunitReporter) className(s string) string {
	s = strings.Title(s)
	return j.invalidClass.ReplaceAllString(s, "")
}

func (j *JunitReporter) name(s string) string {
	s = strings.Title(s)
	return j.invalidName.ReplaceAllString(s, "")
}
