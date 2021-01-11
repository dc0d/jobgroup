package jobgroup_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dc0d/jobgroup"
	"github.com/stretchr/testify/assert"

	_ "github.com/dc0d/jobgroup"
	_ "github.com/stretchr/testify/assert"
)

func TestGroup_default_context_is_nil(t *testing.T) {
	const delay = time.Millisecond * 300
	var defaultContext context.Context

	startedAt := time.Now()

	var grp jobgroup.Group

	grp.Go(func(ctx context.Context) {
		time.Sleep(delay)
		defaultContext = ctx
	})

	_ = grp.Wait()

	elapsed := time.Since(startedAt)
	assert.GreaterOrEqual(t, int64(elapsed), int64(delay))
	assert.Nil(t, defaultContext)
}

func TestGroup_jobs_are_notified_of_cancellation_through_passed_context(t *testing.T) {
	const delay = time.Millisecond * 300

	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()

	grp := jobgroup.New(ctx)

	startedAt := time.Now()
	var elapsed time.Duration
	grp.Go(func(ctx context.Context) {
		<-ctx.Done()
		elapsed = time.Since(startedAt)
	})

	_ = grp.Wait()

	assert.GreaterOrEqual(t, int64(elapsed), int64(delay))
}

func TestGroup_implements_context_interface(t *testing.T) {
	const delay = time.Millisecond * 300

	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()

	grp := jobgroup.New(ctx)
	var gctx context.Context = grp

	startedAt := time.Now()
	var elapsed time.Duration
	grp.Go(func(ctx context.Context) {
		<-gctx.Done()
		elapsed = time.Since(startedAt)
	})

	now := time.Now()
	deadline, ok := gctx.Deadline()
	assert.True(t, deadline.Equal(now) || deadline.After(now))
	assert.True(t, ok)

	_ = grp.Wait()

	assert.True(t, gctx.Err() == context.DeadlineExceeded)
	assert.Equal(t, nil, gctx.Value(""))
	assert.GreaterOrEqual(t, int64(elapsed), int64(delay))
}

func TestGroup_wait_with_timeout(t *testing.T) {
	t.Run(`timed out`, func(t *testing.T) {
		const delay = time.Millisecond * 50

		ctx, cancel := context.WithTimeout(context.Background(), delay)
		defer cancel()

		grp := jobgroup.New(ctx)

		grp.Go(func(context.Context) { select {} })

		err := grp.Wait(time.Millisecond * 300)

		assert.Equal(t, jobgroup.ErrTimeout, err)
		assert.Equal(t, "timed-out", err.Error())
	})

	t.Run(`finished properly`, func(t *testing.T) {
		const delay = time.Millisecond * 50

		ctx, cancel := context.WithTimeout(context.Background(), delay)
		defer cancel()

		grp := jobgroup.New(ctx)

		grp.Go(func(context.Context) {})

		err := grp.Wait(time.Millisecond * 300)

		assert.NoError(t, err)
	})
}

func ExampleGroup() {
	const n = 1000
	var counter int64
	grp := jobgroup.New(context.Background())

	for i := 0; i < n; i++ {
		grp.Go(func(context.Context) {
			atomic.AddInt64(&counter, 1)
		})
	}

	_ = grp.Wait()

	fmt.Println(counter)

	// Output:
	// 1000
}

func ExampleGroup_done() {
	const n = 1000

	var (
		counter int64
		cancel  context.CancelFunc
		grp     *jobgroup.Group
	)

	{
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		grp = jobgroup.New(ctx)
	}

	for i := 0; i < n; i++ {
		go func(done func(), ctx context.Context) {
			defer done()
			<-ctx.Done()
			atomic.AddInt64(&counter, 1)
		}(grp.Add(1))
	}

	go func() { cancel() }()

	_ = grp.Wait()

	fmt.Println(counter)

	// Output:
	// 1000
}

func ExampleGroup_done_ctx() {
	var (
		cancel context.CancelFunc
		grp    *jobgroup.Group
	)

	{
		var ctx context.Context
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		grp = jobgroup.New(ctx)
	}

	go func(done func(), ctx context.Context) {
		defer done()
		<-ctx.Done()
		fmt.Println("job done")
	}(grp.Add(1))

	go func() { cancel() }()

	_ = grp.Wait()

	// Output:
	// job done
}
