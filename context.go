package spec

type onChild func(*suite)

type C struct {
	suite   *suite
	onChild onChild
}

func (c *C) It(name string, test Test) *suite {
	s := Suite(name, test)
	if c.onChild != nil {
		c.onChild(s)
	}
	return s
}
