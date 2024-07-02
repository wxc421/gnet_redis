package iox

type CountWriter struct {
	count int
}

func (c *CountWriter) Write(p []byte) (n int, err error) {
	c.count += len(p)
	return len(p), nil
}

func (c *CountWriter) Count() int {
	return c.count
}
