package jobgroup

import (
	"context"
	"sync"
	"time"
)

type Group struct {
	ctx context.Context
	wg  sync.WaitGroup
}

func New(ctx context.Context) (grp *Group) { return &Group{ctx: ctx} }

func (grp *Group) Go(job func(context.Context)) {
	grp.wg.Add(1)
	go func() {
		defer grp.wg.Done()

		job(grp.ctx)
	}()
}

func (grp *Group) Add(delta int) (done DoneFunc, ctx context.Context) {
	grp.wg.Add(delta)
	return grp.wg.Done, grp.ctx
}

func (grp *Group) Wait(timeout ...time.Duration) error {
	if len(timeout) == 0 || timeout[0] < 0 {
		grp.wg.Wait()
		return nil
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		grp.wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(timeout[0]):
		return ErrTimeout
	}

	return nil
}

// implement context.Context

func (grp *Group) Deadline() (deadline time.Time, ok bool) { return grp.ctx.Deadline() }
func (grp *Group) Done() <-chan struct{}                   { return grp.ctx.Done() }
func (grp *Group) Err() error                              { return grp.ctx.Err() }
func (grp *Group) Value(key interface{}) interface{}       { return grp.ctx.Value(key) }

//

type DoneFunc func()

const (
	ErrTimeout = sentinelError("timed-out")
)

type sentinelError string

func (v sentinelError) Error() string { return string(v) }
