package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DusanDjordjic/mint"
	mintTestSuite "github.com/DusanDjordjic/mint/suites/test"
)

type Executor struct{}

func (e *Executor) Init()   {}
func (e *Executor) Deinit() {}

func (e *Executor) Execute(event mint.Event, trackerEventHandle *mint.TrackerEventHandle) {
	trackerEventHandle.Step("processing", func() {
		trackerEventHandle.Step("println", func() {
			time.Sleep(time.Millisecond * 1)
		})
		trackerEventHandle.Step("something", func() {
			time.Sleep(time.Millisecond * 2)
		})
	})

	step := trackerEventHandle.StepStart("manual sleep")
	time.Sleep(time.Millisecond * 2)
	step.Stop()
}

func TestMintTrackerIntegration(t *testing.T) {
	events := mint.Events{
		{Name: "Dusan", Index: 0},
		{Name: "Cone", Index: 1},
	}

	executor := &Executor{}
	testEvents := []mint.Event{
		events[0], events[1],
		events[0], events[1],
	}

	suite := mintTestSuite.New(testEvents, executor)
	suite.Run()

	report := suite.GetReport()

	t.Run("TreePrintDoesNotPanic", func(t *testing.T) {
		report.ToTree()
	})

	t.Run("JSONExportIsValid", func(t *testing.T) {
		data, err := report.ToJSON()
		if err != nil {
			t.Fatalf("ToJSON failed: %v", err)
		}

		var decoded map[string]any
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal of ToJSON output failed: %v", err)
		}

		if _, ok := decoded["events"]; !ok {
			t.Error("JSON output missing 'events' key")
		}
	})
}
