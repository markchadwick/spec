package spec

type C struct {
	suite *suite
}

func (c *C) It(name string, test Test) {
	c.suite.Add(Suite(name, test))
}
