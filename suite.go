package mint

type Suite interface {
}

type Test []Event

type TestSuite struct {
	Events []Event
	Tests  []Test
}
