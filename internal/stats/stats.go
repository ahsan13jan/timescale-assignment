package stats

import (
	"fmt"
	"math"
	"sort"
)

type StatsOutput struct {
	NumOfQueries        int64
	TotalProcessingTime int64
	MinProcessingTime   int64
	MaxProcessingTime   int64
	AvgProcessingTime   float64
	MedProcessingTime   float64
}

type Stats struct {
	numOfQueries               int64
	totalProcessingTimeInMilli int64
	minProcessingTimeInMilli   int64
	maxProcessingTimeInMilli   int64
	avgProcessingTimeInMilli   float64
	medProcessingTimeInMilli   int64
	processingTimeChan         chan int64
	processingTimes            []int64
}

func New() *Stats {
	stats := &Stats{
		totalProcessingTimeInMilli: 0,
		numOfQueries:               0,
		minProcessingTimeInMilli:   9223372036854775807,
		maxProcessingTimeInMilli:   -1,
		avgProcessingTimeInMilli:   0,
		medProcessingTimeInMilli:   0,
		processingTimeChan:         make(chan int64),
	}

	return stats
}

func (s *Stats) Metrics() StatsOutput {
	return StatsOutput{
		NumOfQueries:        s.numOfQueries,
		TotalProcessingTime: s.totalProcessingTimeInMilli,
		MinProcessingTime:   s.minProcessingTimeInMilli,
		MaxProcessingTime:   s.maxProcessingTimeInMilli,
		AvgProcessingTime:   math.Floor(s.avgProcessingTimeInMilli*100) / 100,
		MedProcessingTime:   math.Floor(calculateMedian(s.processingTimes)*100) / 100,
	}
}

// TODO profiling shows it is taking alot or memory
// Need to look into how can we improve memory
func (s *Stats) Listen() {
	for p := range s.processingTimeChan {

		s.processingTimes = append(s.processingTimes, p)
		s.numOfQueries++
		s.totalProcessingTimeInMilli += p

		if p < s.minProcessingTimeInMilli {
			s.minProcessingTimeInMilli = p
		}

		if p > s.maxProcessingTimeInMilli {
			s.maxProcessingTimeInMilli = p
		}

		s.avgProcessingTimeInMilli = float64(s.totalProcessingTimeInMilli) / float64(s.numOfQueries)
	}
}

func (s *Stats) Report(processingTime int64) {
	s.processingTimeChan <- processingTime
}

func (s *StatsOutput) Print() {
	fmt.Printf("numOfQueries :%d\n", s.NumOfQueries)
	fmt.Printf("totalProcessingTimeInMilli :%d\n", s.TotalProcessingTime)
	fmt.Printf("minProcessingTimeInMilli :%d\n", s.MinProcessingTime)
	fmt.Printf("maxProcessingTimeInMilli :%d\n", s.MaxProcessingTime)
	fmt.Printf("avgProcessingTimeInMilli :%.2f\n", s.AvgProcessingTime)
	fmt.Printf("medProcessingTimeMilli :%.2f\n", s.MedProcessingTime)
}

func (s *Stats) Close() {
	close(s.processingTimeChan)
}

// Time complexity of this algorithm in O(n)
//
//	TODO Improve time complexity with using min max heap data structure
func calculateMedian(arr []int64) float64 {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	mid := len(arr) / 2
	if len(arr)%2 == 0 {

		return float64(arr[mid-1]+arr[mid]) / float64(2)
	} else {
		return float64(arr[mid])
	}
}
