package spec

var nilReporter = &testReporter{}

type testReporter struct {
	lastErrors []*TestError
	lastSkip   *TestError
}

func (t *testReporter) Start(s *suite) {
}

func (t *testReporter) Pass(s *suite) {
}

func (t *testReporter) Fail(s *suite, errs []*TestError) {
	t.lastErrors = errs
}

func (t *testReporter) Skip(s *suite, skip *TestError) {
	t.lastSkip = skip
}

func (t *testReporter) Descend(s *suite) {
}

func (t *testReporter) Ascend(s *suite) {
}

func (t *testReporter) Begin() {
	t.lastErrors = nil
	t.lastSkip = nil
}

func (t *testReporter) Finish(errs []*SuiteFailure) {
}
