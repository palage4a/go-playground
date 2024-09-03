package github.com/palage4a/go-playground/pkg/core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	// "github.com/tarantool/message-queue-ee/internal/config"
	// "github.com/tarantool/message-queue-ee/internal/dto"
	// "github.com/tarantool/message-queue-ee/pkg/health"
	// "github.com/tarantool/message-queue-ee/pkg/log"


	"github.com/tarantool/go-iproto"
    "github.com/palage4a/go-playground/pkg/queue_pool
)

type Core struct {
	pool  *queuePool
	pools map[string]*queuePool
}

func New(ctx context.Context, clientCfg *config.TarantoolClientConfig) (*DB, error) {
	pools := make(map[string]*queuePool, len(clientCfg.Queues))

	var p *queuePool
	connections := clientCfg.Connections

	if connections != nil {
	}

	for queueName := range clientCfg.Queues {
	}

	return &DB{
		pool:  p,
		pools: pools,
	}, nil
}

func (r *DB) NextPool(queue string) (*pool.ConnectionPool, error) {
	pool, ok := r.pools[queue]
	if !ok {
		pool = r.pool
		if pool == nil {
			return nil, fmt.Errorf("db: queue `%s` does not exists in publisher config", queue)
		}
	}

	if len(pool.Replicasets) == 1 {
		return pool.Connections[pool.Replicasets[pool.Cursor]], nil
	}

	pool.Mutex.Lock()
	defer pool.Mutex.Unlock()

	replicaset := pool.Replicasets[pool.Cursor]
	if len(pool.Connections)-1 == pool.Cursor {
		pool.Cursor = 0
	} else {
		pool.Cursor++
	}

	return pool.Connections[replicaset], nil
}

func (r *DB) Each(cb func(shard string, connection *pool.ConnectionPool) error) error {
	if r.pool != nil {
		for shard, connection := range r.pool.Connections {
			err := cb(shard, connection)
			if err != nil {
				return err
			}
		}
	}

	for queueName := range r.pools {
		for shard, connection := range r.pools[queueName].Connections {
			err := cb(shard, connection)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *DB) EachWithQueue(queue string, cb func(shard string, connection *pool.ConnectionPool) error) error {
	pool, ok := r.pools[queue]
	if !ok {
		pool = r.pool
		if pool == nil {
			return fmt.Errorf("db: queue `%s` does not exists in publisher config", queue)
		}
	}

	for shard, connection := range pool.Connections {
		err := cb(shard, connection)
		if err != nil {
			return err
		}
	}

	return nil
}

func Subscribe(
	ctx context.Context,
	conn *pool.ConnectionPool,
	timeout time.Duration,
	queue string,
	rk *string,
	cursor *uint64,
	sk *string,
	sks []string,
) ([]*dto.PersistedMessage, error) {
	res, err := CallAny[[]*dto.PersistedMessage](
		ctx,
		conn,
		"queue.subscribe",
		timeout/time.Millisecond,
		queue,
		rk,
		cursor,
		sk,
		sks,
	)
	if err != nil {
		return nil, err
	}

	return *res, nil
}

func (r *DB) Publish(ctx context.Context, queue string, msg dto.PublishMessage) (uint64, error) {
	conn, err := r.nextPool(queue)
	if err != nil {
		return 0, err
	}

	res, err := CallRW[uint64](ctx, conn, "queue.publish", msg)
	if err != nil {
		return 0, err
	}

	return *res, nil
}

func (r *DB) PublishBatch(
	ctx context.Context,
	queue string,
	shardingKey *string,
	messages dto.InBatchMessages,
) ([]uint64, error) {
	conn, err := r.nextPool(queue)
	if err != nil {
		return []uint64{}, err
	}

	res, err := CallRW[[]uint64](ctx, conn, "queue.publish_batch", queue, shardingKey, messages)
	if err != nil {
		return []uint64{}, err
	}

	return *res, nil
}

func LivenessProbe(db *DB) health.Check {
	return func() error {
		return db.Each(func(_ string, connection *pool.ConnectionPool) error {
			_, err := connection.Do(
				tarantool.NewCallRequest("box.info").Args([]any{}),
				pool.ANY,
			).Get()

			return err
		})
	}
}
