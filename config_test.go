package mint

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	type TestCase struct {
		Input         io.Reader
		Expected      config
		ExpectedError bool
		ExpectedPanic bool
	}

	cases := map[string]TestCase{
		"nil reader": {
			Input:         nil,
			ExpectedPanic: true,
		},
		"one event": {
			Input: bytes.NewBufferString(
				`
{
	"events": ["1"],
	"eventsTable": {
		"probabilities": [100]
	}
}
`),
			ExpectedPanic: false,
			ExpectedError: false,
			Expected: config{
				Events: []Event{"1"},
				EventsTable: eventsTable{
					Probabilities: []uint16{100},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				err := recover()
				if c.ExpectedPanic {
					require.NotNil(t, err)
				} else {
					require.Nil(t, err)
				}
			}()

			config, err := ParseConfig(c.Input)
			if c.ExpectedPanic {
				t.Fatalf("expected panic but one didn't happen, %v %s", config, err)
				return
			}

			if c.ExpectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, c.Expected, config)
		})

	}
}
