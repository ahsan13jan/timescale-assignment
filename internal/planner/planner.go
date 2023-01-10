package planner

import (
	"math/rand"
	"time"

	"timescale/internal/db"
	"timescale/internal/stats"
)

type hosts []string
type Planner struct {
	client     *db.PGClient
	stats      *stats.Stats
	workers    map[*worker]hosts
	maxWorkers int
}

func New(client *db.PGClient, stats *stats.Stats, maxWorkers int) *Planner {
	return &Planner{
		workers:    make(map[*worker]hosts, maxWorkers),
		client:     client,
		stats:      stats,
		maxWorkers: maxWorkers,
	}
}

// Add some unit test, some functionality is convered by bdd tests
func (p Planner) Plan(input chan Query) {
	for query := range input {
		w := p.workerExists(query.Host)
		if w == nil {
			// limits numbers of workers
			w = p.createOrAssignWorker(query.Host)
		}
		w.push(query)
	}

	// Closing all query channels of workers
	for w := range p.workers {
		w.close()
	}
}

func (p Planner) workerExists(host string) *worker {
	for w, hosts := range p.workers {
		for _, h := range hosts {
			if h == host {
				return w
			}
		}
	}
	return nil
}

// Ensuring Max concurrent workers
func (p Planner) createOrAssignWorker(host string) *worker {
	if len(p.workers) < p.maxWorkers {
		w := newWorker(p.client, p.stats)
		p.workers[w] = []string{host}
		return w
	} else {
		// number of concurrent workers maxed out, randomly assign host to worker
		randomIndex := randomNum(p.maxWorkers-1, 0)
		index := 0
		for w, h := range p.workers {
			if index == randomIndex {
				h = append(h, host)
				p.workers[w] = h
				return w
			}
			index++
		}
	}
	return nil
}

// TODO add unit tests
func randomNum(max, min int) int {
	rand.Seed(time.Now().UnixNano())
	if max == min {
		return max
	}
	return rand.Intn(max-min) + min
}

func (p Planner) WorkersDone() chan bool {
	status := make(chan bool, 1)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		for range ticker.C {
			totalWorkers := len(p.workers)
			doneWorkers := 0
			for w := range p.workers {
				if w.don() {
					doneWorkers++
				}
			}
			if totalWorkers == doneWorkers {
				status <- true
				close(status)
				break
			}
		}
	}()

	return status
}
