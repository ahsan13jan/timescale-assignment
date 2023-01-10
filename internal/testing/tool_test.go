//go:build integration
// +build integration

package testing

import (
	"fmt"
	"strconv"

	"timescale/internal/run"
)

func theCliToolExecutes(csvPath string) error {
	metrics = run.Execute(csvPath, -1, conf)
	metrics.Print()
	return nil
}

func theNumOfQuriesShouldBe(numStr string) error {
	var (
		num int64
		err error
	)
	if num, err = strconv.ParseInt(numStr, 10, 64); err != nil {
		return fmt.Errorf("could not parse days: %+v", err)
	}

	if num != metrics.NumOfQueries {
		return fmt.Errorf("num of quries stats not equal")
	}
	return nil
}
