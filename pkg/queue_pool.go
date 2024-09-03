package github.com/palage4a/go-playground/pkg/queue_pool

import (
    "sync"
)

type Publisher interface {}
type PublishBatcher interface {}
type Subscriber interface {}

type queuePool struct {
	Connections map[string]Publisher
	Replicasets []string
	Cursor      int
	Mutex       *sync.Mutex
}

func New(queue string, connections map[string]Publisher) {
    p = &queuePool{
        Connections: make(map[string]Publisher, len(connections)),
        Replicasets: make([]string, len(connections)),
        Cursor:      0,
        Mutex:       &sync.Mutex{},
    }

    i := 0
    for shard, addrs := range connections {
        p.Connections[shard] = publisher
        p.Replicasets[i] = shard
        i++
    }
}
