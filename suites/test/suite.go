package test

import "github.com/DusanDjordjic/mint"

type TestSuite struct {
	test     []mint.Event
	executor mint.Executor
}

func New(test []mint.Event, executor mint.Executor) TestSuite {
	return TestSuite{
		test:     test,
		executor: executor,
	}
}

func (suite *TestSuite) Run() {
	suite.executor.Init()
	for _, event := range suite.test {
		suite.executor.Execute(event)
	}
	suite.executor.Deinit()
}
