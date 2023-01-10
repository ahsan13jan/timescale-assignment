package stats

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats_Listen(t *testing.T) {
	stats := New()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		stats.Listen()
		wg.Done()
	}()

	stats.Report(10)
	stats.Report(20)
	stats.Report(40)
	stats.Close()

	wg.Wait()

	expectedStats := StatsOutput{
		NumOfQueries:        3,
		TotalProcessingTime: 70,
		MinProcessingTime:   10,
		MaxProcessingTime:   40,
		AvgProcessingTime:   23.33,
		MedProcessingTime:   20,
	}
	actualStats := stats.Metrics()

	assert.Equal(t, expectedStats, actualStats)
}

func TestMedian(t *testing.T) {
	assert.Equal(t, float64(10), calculateMedian([]int64{10, 20, 5}))
	assert.Equal(t, float64(15), calculateMedian([]int64{10, 20, 5, 60}))
}
