package planner

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Query struct {
	Host       string
	Start, End time.Time
}

func StreamInput(filePath string, out chan<- Query) {

	var (
		scanner *bufio.Scanner
		file    *os.File
		err     error
	)

	if filePath != "" {
		if file, err = os.Open(filePath); err != nil {
			// TODO handle error gracefully
			log.WithError(err).Fatalf("could not open file with path:%s", filePath)
			return
		}
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	defer func() {
		if file != nil {
			file.Close()
		}
		close(out)
	}()

	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
		if count == 1 {
			continue
		}
		line := scanner.Text()
		split := strings.Split(line, ",")
		host := split[0]
		start, err := toTime(split[1])
		if err != nil {
			log.WithError(err).Error("invalid start time")
			continue
		}

		end, err := toTime(split[2])
		if err != nil {
			log.WithError(err).Error("invalid end time")
			continue
		}

		out <- Query{Host: host, Start: *start, End: *end}

	}
}

func toTime(timeStr string) (*time.Time, error) {
	layout := "2006-01-02 15:04:05"
	ts, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil, err
	}
	return &ts, nil
}
