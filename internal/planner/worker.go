package planner

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"timescale/internal/db"
	"timescale/internal/logger"
	"timescale/internal/stats"
)

var log = logger.GetLogger()

type worker struct {
	db          *db.PGClient
	stats       *stats.Stats
	queriesChan chan Query
	done        bool
	id          string
}

func (w *worker) run() {
	ctx := context.Background()
	for query := range w.queriesChan {
		log.WithFields(logrus.Fields{"worker_id": w.id, "host": query.Host}).Debug("Worker querying")
		now := time.Now()
		_, err := w.db.GetStatsPerHostPerMin(ctx, query.Host, query.Start, query.End)
		if err != nil {
			panic(err)
		}
		w.stats.Report(time.Since(now).Milliseconds())
	}
	log.WithField("worker_id", w.id).Debug("Worker done")
	w.done = true
}

func (w worker) don() bool {
	return w.done
}

func (w worker) push(q Query) {
	w.queriesChan <- q
}
func (w worker) close() {
	close(w.queriesChan)
}

func newWorker(db *db.PGClient, stats *stats.Stats) *worker {
	worker := &worker{
		db:          db,
		queriesChan: make(chan Query, 10),
		id:          uuid.New().String(),
		stats:       stats,
	}

	log.WithField("worker_id", worker.id).Debug("Worker created")
	go worker.run()
	return worker
}
