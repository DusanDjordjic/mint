package mint

type Executor interface {
	Init()
	Deinit()
	Execute(event Event, tracker *Tracker) // test result ..
}
