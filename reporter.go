package spec

import (
	"fmt"
)

type Reporter interface {
	Start(*suite)
	Pass(*suite)
	Fail(*suite, []error)
}

// ----------------------------------------------------------------------------
// ASCII Reporter
// ----------------------------------------------------------------------------

type Ascii struct {
	stack []string
}

func AsciiReporter() Reporter {
	return &Ascii{
		stack: make([]string, 0),
	}
}

func (a *Ascii) Start(s *suite) {
	a.stack = append(a.stack, s.Name)

	a.printPad()
	fmt.Printf("* %s\n", s.Name)
}

func (a *Ascii) Pass(s *suite) {
	a.stack = a.stack[:len(a.stack)-1]
}

func (a *Ascii) Fail(s *suite, errs []error) {
	a.stack = a.stack[:len(a.stack)-1]
}

func (a *Ascii) printPad() {
	for i := 0; i < len(a.stack); i++ {
		fmt.Print("  ")
	}
}

// ----------------------------------------------------------------------------
// Console Reporter
// ----------------------------------------------------------------------------

// Prints curses-style progress to the console. Absolutely cannot be used at the
// same time as other reporters printing to the console.
type Console struct {
}
