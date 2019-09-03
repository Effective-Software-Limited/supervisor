// Package supervisor implements a common pattern for network programs: multiple
// goroutines are listening at network sockets, and the first error should be transmitted
// to all the other routines.
//
// See also https://godoc.org/golang.org/x/sync/errgroup
package supervisor

import (
	"context"
	"sync"
)

// Supervisor is the top-level group
type Supervisor struct {
	cancel func()
	ctx    context.Context

	defers []func()

	err      error
	errfirst sync.Once
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time an agent returns a non-nil error, or
// when all agents are successfully completed.
func WithContext(ctx context.Context) (*Supervisor, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Supervisor{cancel: cancel, ctx: ctx}, ctx
}

// Agent adds a new routine to the group.
func (s *Supervisor) Agent(f func() error) {

	go func() {
		if err := f(); err != nil {
			s.errfirst.Do(func() { s.err = err })
		}
		s.cancel()
	}()
}

// Err returns a channel that will deliver a single value on the first error in the group,
// or nil if no error occurred, but there are no routines left in the group.
func (s *Supervisor) Err() <-chan error {
	errchan := make(chan error)
	go func() {
		<-s.ctx.Done()
		errchan <- s.err
	}()
	return errchan
}
