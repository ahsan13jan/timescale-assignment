package planner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParTime(t *testing.T) {
	ts, err := toTime("2017-01-01 08:59:22")
	assert.NoError(t, err)
	expectedTime := time.Date(2017, time.Month(1), 1, 8, 59, 22, 0, time.UTC)
	assert.Equal(t, expectedTime, *ts)
}

func TestStream(t *testing.T) {

	out := make(chan Query, 1)
	StreamInput("../../testdata/query_params_unit.csv", out)
	assert.Equal(t, 1, len(out))

	s, err := toTime("2017-01-01 08:59:22")
	assert.NoError(t, err)

	e, err := toTime("2017-01-01 09:59:22")
	assert.NoError(t, err)
	expectedQuery := Query{
		Host:  "host_000008",
		Start: *s,
		End:   *e,
	}
	// 2017-01-01 08:59:22,2017-01-01 09:59:22
	for q := range out {
		assert.Equal(t, expectedQuery, q)
	}
}
