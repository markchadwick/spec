package spec

var nilReporter = &testReporter{}

type testReporter struct{}

func (t *testReporter) Start(s *suite) {
}

func (t *testReporter) Pass(s *suite) {
}

func (t *testReporter) Fail(s *suite, errs []error) {
}
