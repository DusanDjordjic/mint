package test

import "github.com/DusanDjordjic/mint"

type TestSuite struct {
	test     []mint.Event
	executor mint.Executor
	tracker  mint.Tracker
}

func New(test []mint.Event, executor mint.Executor) TestSuite {
	return TestSuite{
		test:     test,
		executor: executor,
		tracker:  *mint.NewTracker(),
	}
}

func (suite *TestSuite) Run() {
	suite.executor.Init()

	// cold start.. prvom treba 14 mikro sekundi svakom sledecem 1-2 max
	// TODO: istraziti cold start

	for _, event := range suite.test {
		eventTracker := suite.tracker.EventStart(event.Name)
		suite.executor.Execute(event, eventTracker)
	}
	suite.executor.Deinit()
}

func (suite *TestSuite) GetReport() *mint.Report {
	return suite.tracker.GenerateReport()
}
