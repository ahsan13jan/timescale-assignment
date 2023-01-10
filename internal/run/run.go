package run

import (
	"context"

	"golang.org/x/sync/errgroup"

	"timescale/internal/config"
	"timescale/internal/db"
	"timescale/internal/logger"
	plan "timescale/internal/planner"
	"timescale/internal/stats"
)

var log = logger.GetLogger()

func Execute(fileName string, maxWorkers int, c config.Config) stats.StatsOutput {
	if fileName != "" {
		c.CsvPath = fileName
	}

	if maxWorkers != -1 {
		c.MaxConcurrentWorkers = maxWorkers
	}

	ctx := context.Background()

	client, err := db.New(c)
	if err != nil {
		log.Fatalf("failed to initialise db client:%v", err)
	}
	defer client.Close(ctx)

	stats := stats.New()

	errG, _ := errgroup.WithContext(ctx)

	errG.Go(func() error {
		stats.Listen()
		return nil
	})

	planner := plan.New(client, stats, c.MaxConcurrentWorkers)
	queryChan := make(chan plan.Query)

	errG.Go(func() error {
		plan.StreamInput(c.CsvPath, queryChan)
		return nil
	})

	errG.Go(func() error {
		planner.Plan(queryChan) // plan queries
		<-planner.WorkersDone() // wait for all the workers to be done
		stats.Close()           // close stats
		return nil
	})

	if err := errG.Wait(); err != nil {
		log.WithError(err).Error("errgroup error")
	}

	log.Debugf("Workers done")
	return stats.Metrics()
}
