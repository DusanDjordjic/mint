package mint

type Executor interface {
	Init()
	Deinit()
	Execute(event Event, eventTracker *TrackerEventHandle) // test result ..
}
