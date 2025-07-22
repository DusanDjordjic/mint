package test

import (
	"github.com/DusanDjordjic/mint"
	"github.com/DusanDjordjic/mint/tracker"
)

type TestSuite struct {
	test     []mint.Event
	executor mint.Executor
	tracker  *tracker.Tracker
}

func New(test []mint.Event, executor mint.Executor, t *tracker.Tracker) TestSuite {
	return TestSuite{
		test:     test,
		executor: executor,
		tracker:  t,
	}
}

func (suite *TestSuite) Run() {
	suite.executor.Init(suite.tracker)
	for _, event := range suite.test {
		entry := suite.tracker.Start(event.Name)
		suite.executor.Execute(event, suite.tracker)
		entry.Stop()
	}
	suite.executor.Deinit(suite.tracker)
}
