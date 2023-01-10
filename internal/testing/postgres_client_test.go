//go:build integration
// +build integration

package testing

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

func theGivenCpuUsage(table *godog.Table) error {

	cpuUsages, err := godogTableToRetention(table)
	if err != nil {
		return err
	}

	for _, usage := range cpuUsages {
		err := clientLocal.Create(context.Background(), usage.ts, usage.host, usage.usage)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return nil
}

func theAggregateShouldBe(host, start, end string, table *godog.Table) error {

	startT, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return fmt.Errorf("could not parse the date %+v", err)
	}

	endT, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return fmt.Errorf("could not parse the date %+v", err)
	}

	rows, err := client.GetStatsPerHostPerMin(context.Background(), host, startT, endT)
	if err != nil {
		log.Fatalln(err)
	}

	if len(rows)+1 != len(table.Rows) {
		return fmt.Errorf("actual and expected rows not equal")
	}

	for i, row := range table.Rows {
		var (
			ts                 time.Time
			minUsage, maxUsage float64
			err                error
		)
		if i == 0 {
			continue
		}

		actualRow := rows[i-1]

		if ts, err = time.Parse(time.RFC3339, row.Cells[0].Value); err != nil {
			return fmt.Errorf("could not parse ts: %+v", err)
		}

		if ts != actualRow.Ts.UTC() {
			return fmt.Errorf("ts not equal expected:%s, actual:%s", ts, actualRow.Ts.UTC())
		}

		if minUsage, err = strconv.ParseFloat(row.Cells[1].Value, 64); err != nil {
			return fmt.Errorf("could not parse days: %+v", err)
		}

		if float32(minUsage) != actualRow.MinUsage {
			return fmt.Errorf("MinUsage not equal: %+v", err)
		}

		if maxUsage, err = strconv.ParseFloat(row.Cells[2].Value, 64); err != nil {
			return fmt.Errorf("could not parse days: %+v", err)
		}

		if float32(maxUsage) != actualRow.MaxUsage {
			return fmt.Errorf("maxUsage not equal: %+v", err)
		}
	}

	return nil
}

type cpuUsage struct {
	ts    time.Time
	host  string
	usage float64
}

func godogTableToRetention(table *godog.Table) ([]cpuUsage, error) {

	var cpuUsages []cpuUsage
	for i, row := range table.Rows {
		var (
			ts    time.Time
			host  string
			usage float64
			err   error
		)
		if i == 0 {
			continue
		}

		host = row.Cells[0].Value

		if ts, err = time.Parse(time.RFC3339, row.Cells[1].Value); err != nil {
			return nil, fmt.Errorf("could not parse ts: %+v", err)
		}

		if usage, err = strconv.ParseFloat(row.Cells[2].Value, 64); err != nil {
			return nil, fmt.Errorf("could not parse days: %+v", err)
		}
		cpuUsages = append(cpuUsages, cpuUsage{ts: ts, host: host, usage: usage})
	}

	return cpuUsages, nil
}
