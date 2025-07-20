package test

import "github.com/DusanDjordjic/mint"

type TestSuite struct {
	test     []mint.Event
	executor mint.Executor
	tracker  mint.Tracker
}

func New(test []mint.Event, executor mint.Executor, tracker mint.Tracker) TestSuite {
	return TestSuite{
		test:     test,
		executor: executor,
		tracker:  tracker,
	}
}

func (suite *TestSuite) Run() {
	suite.executor.Init()

	// cold start.. prvom treba 14 mikro sekundi svakom sledecem 1-2 max
	// TODO: istraziti cold start
	suite.executor.Execute(suite.test[0], &suite.tracker)
	suite.tracker.Reset()

	for _, event := range suite.test {
		suite.tracker.EventStart(event.Name)
		suite.executor.Execute(event, &suite.tracker)
		suite.tracker.EventStop()
	}
	suite.executor.Deinit()
	suite.tracker.PrintReportTree()
}
