package tests

import (
	"fmt"
	"testing"

	"github.com/DusanDjordjic/mint"
	mintTestSuite "github.com/DusanDjordjic/mint/suites/test"
)

type Executor struct{}

func (e *Executor) Init() {
	fmt.Println("Init")
}
func (e *Executor) Deinit() {
	fmt.Println("Deinit")
}
func (e *Executor) Execute(event mint.Event) {

	fmt.Println("Executing...", event.Name, event.Index)
}

func TestTracker(t *testing.T) {
	events := make(mint.Events, 2)
	e := mint.Event{
		Name:  "Dusan",
		Index: 0,
	}
	events[0] = e
	e = mint.Event{
		Name:  "Cone",
		Index: 1,
	}
	events[1] = e

	exe := &Executor{}

	test := []mint.Event{
		events[0], events[1], events[1], events[0],
	}

	s := mintTestSuite.New(test, exe)
	s.Run()
}
