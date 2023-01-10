package planner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrAssignWorkers(t *testing.T) {

	planner := New(nil, nil, 1)

	t.Run("adding two hosts", func(t *testing.T) {
		worker1 := planner.createOrAssignWorker("host1")
		assert.Equal(t, 1, len(planner.workers))
		hosts1, found1 := planner.workers[worker1]
		assert.Equal(t, true, found1)
		assert.Equal(t, hosts{"host1"}, hosts1)

		// should be the same as worker 1 wince max concurrency is 1
		worker2 := planner.createOrAssignWorker("host2")
		assert.Equal(t, worker1, worker2)

		//
		actualHosts, found := planner.workers[worker1]
		assert.Equal(t, true, found)
		assert.Equal(t, hosts{"host1", "host2"}, actualHosts)

	})

}
