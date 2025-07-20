package mint

type Events []Event

type Event struct {
	Name  string
	Index int
}

// TODO implemenet UnmarshalJSON for Events so we can load them from an array of strings
// and auto-generate the indexes starting from 0
