package github.com/palage4a/go-playground/pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/tarantool/go-iproto"

	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/uuid"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

type Pool struct {
    p tarantool.Doer
}

func New(ctx context.Context, user string, pass string, addrs []string) (*Pool, error) {
	instances := make([]pool.Instance, len(addrs))
	for i_addr, addr := range addrs {
		instances[i_addr] = pool.Instance{
			Name: addr,
			Dialer: tarantool.NetDialer{
				Address:  addr,
				User:     user,
				Password: pass,
			},
			Opts: tarantool.Opts{},
		}
	}

	p, err := retry.DoWithData(func() (*pool.ConnectionPool, error) {
		return pool.Connect(ctx, instances)
	},
		retry.Attempts(5),
		retry.Delay(3*time.Second),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Infof("failed to connect after %d attempts: %s", n+1, err)
		}),
	)
	if err != nil {
		return nil, err
	}

    return &Pool{p:p}, nil
}

func (p *Pool) Publish(ctx context.Context, queue string, msg dto.PublishMessage) (uint64, error) {
    return callRW[uint64](ctx, p.p, "queue.publish", msg)
}

func (p *Pool) PublishBatch(
	ctx context.Context,
	queue string,
	shardingKey *string,
	messages dto.InBatchMessages,
) ([]uint64, error) {
    return callRW[[]uint64](ctx, p.p, "queue.publish_batch", queue, shardingKey, messages)
}

func (p *Pool) Subscribe(
	ctx context.Context,
	timeout time.Duration,
	queue string,
	rk *string,
	cursor *uint64,
	sk *string,
	sks []string,
) ([]*dto.PersistedMessage, error) {
	return callAny[[]*dto.PersistedMessage](
		ctx,
		p.p,
		"queue.subscribe",
		timeout/time.Millisecond,
		queue,
		rk,
		cursor,
		sk,
		sks,
	)
}

var ErrDeduplication = errors.New("duplicate found")

func call[R any](ctx context.Context, p tarantool.Doer , userMode pool.Mode, name string, args ...any) (*R, error) {
	fut := p.Do(
		tarantool.NewCallRequest(name).Args(args),
		userMode,
	)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-fut.WaitChan():
	}

	var resp Resp[*R]
	if err := fut.GetTyped(&resp); err != nil {
		if tntErr, ok := err.(tarantool.Error); ok && tntErr.Code == iproto.ER_TUPLE_FOUND {
			tntErr.ExtendedInfo = nil // NOTE: remove stack trace from error message
			err = fmt.Errorf("%w: %w", ErrDeduplication, tntErr)
		}

		return nil, fmt.Errorf("tarantool: %w when %s(%+v)", err, name, args)
	}

	if resp.Value == nil {
		err := fmt.Errorf("tarantool: empty result of %s(%+v)", name, args)
		log.Error(err)

		return nil, err
	}

	return resp.Value, nil
}

func callRW[R any](ctx context.Context, p *pool.ConnectionPool, name string, args ...any) (*R, error) {
	return call[R](ctx, p, pool.RW, name, args...)
}

func callAny[R any](ctx context.Context, p *pool.ConnectionPool, name string, args ...any) (*R, error) {
	return call[R](ctx, p, pool.ANY, name, args...)
}

