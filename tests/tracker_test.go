package tests

import (
	"fmt"
	"testing"
	"time"

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

func (e *Executor) Execute(event mint.Event, tracker *mint.Tracker) {
	tracker.Step("processing", func() {
		tracker.Step("println", func() {
			fmt.Println("Executing...", event.Name, event.Index)
		})
		tracker.Step("something", func() {
			time.Sleep(time.Millisecond * 1)
		})
	})

	// ILI

	tracker.StepStart("time sleep")
	time.Sleep(time.Millisecond * 2)
	tracker.StepStart("123")
	time.Sleep(time.Millisecond * 1)
	tracker.StepStop("123")
	tracker.StepStop("time sleep")
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
	tracker := mint.NewTracker()

	test := []mint.Event{
		events[0], events[1], events[1], events[0],
		events[0], events[1], events[1], events[0],
		events[0], events[1], events[1], events[0],
		events[0], events[1], events[1], events[0],
		events[0], events[1], events[1], events[0],
		events[0], events[1], events[1], events[0],
	}

	s := mintTestSuite.New(test, exe, *tracker)
	s.Run()
}
