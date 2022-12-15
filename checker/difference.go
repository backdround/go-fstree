package checker

// Difference type describes specific difference between filesystem
// and expectation
type Difference struct {
	Path string
	Expectation string
	Real string
}
