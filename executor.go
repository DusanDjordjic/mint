package mint

type Executor interface {
	Init()
	Deinit()
	Execute(event Event) // test result ..
}
