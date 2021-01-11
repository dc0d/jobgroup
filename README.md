[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![PkgGoDev](https://pkg.go.dev/badge/dc0d/jobgroup)](https://pkg.go.dev/github.com/dc0d/jobgroup) [![Go Report Card](https://goreportcard.com/badge/github.com/dc0d/jobgroup)](https://goreportcard.com/report/github.com/dc0d/jobgroup) [![Maintainability](https://api.codeclimate.com/v1/badges/33f3205c4f3c848e065b/maintainability)](https://codeclimate.com/github/dc0d/jobgroup/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/33f3205c4f3c848e065b/test_coverage)](https://codeclimate.com/github/dc0d/jobgroup/test_coverage)


# jobgroup

To manage a group of goroutines, it's possible to use an instance of `context.Context`. And to wait for them to finish, it's possible to use an instance of `sync.WaitGroup`. The `jobgroup` provides a utility, which combines these functionalities, under a simple API.

## example usage

Assume a set of jobs need to be finished inside a time-window. For this purpose a context can be used.

```go
ctx, cancel := context.WithTimeout(context.Background(), delay)
defer cancel()
```

Then, using that context, we create a job-group.

```go
grp := jobgroup.New(ctx)
```

Each goroutine can be registered in this job-group.

```go
grp.Go(func(ctx context.Context) {
    // ...
    // cancel this job if the context is canceled/timeout
    // <-ctx.Done()
})
```

And then we wait for them to finish.

```go
grp.Wait()
```

Also, it's possible to wait for them to finish for limited duration of time.

```go
err := grp.Wait(time.Second * 5)
```

The returned `err` will be nil, if all jobs are finished. Or `jobgroup.ErrTimeout` if any job fails to stop after five seconds of waiting.
