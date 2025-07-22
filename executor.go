package mint

import "github.com/DusanDjordjic/mint/tracker"

type Executor interface {
	Init(tracker *tracker.Tracker)
	Deinit(tracker *tracker.Tracker)
	Execute(event Event, tracker *tracker.Tracker)
}
