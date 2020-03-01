// This is a package that contains an implementation of an abstract
// synchronisation mechanism - synchronisation group.
// The main idea is to have an ability to run independent tasks in separate goroutines which way return errors.
// A user can wait until all goroutines finish running and collect all occurred errors.
//
// The design is similar to errgroup (https://godoc.org/golang.org/x/sync/errgroup),
// but it does not cancel the context of the goroutines if any of them returns an error.
package syncgroup

import (
	"github.com/pkg/errors"
	"strings"
	"sync"
)

// SyncGroup is the main class for working with syncgroups.
// It has two main methods: .Go() and .Wait()
//
// .Go() spawns a new goroutine, which may return an error.
// The returned error will be saved and returned by .Wait() method.
//
// .Wait() waits until all spawned goroutines finish and returns an Error struct which is a wrapper for a slice of errors.
// If there was no error, .Wait() would return nil, otherwise it would return GroupError instance.
type SyncGroup struct {
	wg sync.WaitGroup

	finishedChan chan []error
	errorChan    chan error

	listeningStarted bool
}

// GroupError is a wrapper for errors which are returned by functions called in spawned goroutines.
type GroupError struct {
	Errs []error
}

// Error concatenates all stored errors with ';' symbol and returns a resulting string.
func (e GroupError) Error() string {
	builder := strings.Builder{}

	for _, err := range e.Errs {
		if err != nil {
			builder.WriteString(err.Error())
			builder.WriteString(";")
		}
	}

	return builder.String()
}

// New is the default constructor for SyncGroup
func New() *SyncGroup {
	g := &SyncGroup{
		wg:           sync.WaitGroup{},
		finishedChan: make(chan []error),
		errorChan:    make(chan error),
	}

	return g
}

func (g *SyncGroup) listenToErrors() {
	var accumulatedErrors []error

	for {
		err, ok := <-g.errorChan
		if ok {
			accumulatedErrors = append(accumulatedErrors, err)
		} else {
			break
		}
	}

	g.finishedChan <- accumulatedErrors
	close(g.finishedChan)

	return
}

// Go spawns given function in a new goroutine.
// The returned error will be saved and returned by .Wait() method.
func (g *SyncGroup) Go(f func() error) {
	if !g.listeningStarted {
		go g.listenToErrors()
		g.listeningStarted = true
	}

	g.wg.Add(1)
	go func() {
		defer func() {
			if msg := recover(); msg != nil {
				switch msg.(type) {
				case error:
					g.errorChan <- errors.Wrap(msg.(error), "recovered from panic")
				default:
					g.errorChan <- errors.Errorf("recovered from panic:%v", msg)
				}
			}

			g.wg.Done()
		}()

		err := f()

		if err != nil {
			g.errorChan <- err
		}
	}()
}

// Wait waits until all spawned goroutines are finished and returns a wrapper struct for all collected errors.
// The result is nil if none of the spawned goroutines returned an error
//
// The result is guaranteed to be an instance of GroupError, so that you can access the stored errors directly.
// If you only need to check the absence of errors, then only check for nil value.
func (g *SyncGroup) Wait() error {
	if !g.listeningStarted {
		return nil
	}

	g.wg.Wait()
	close(g.errorChan)

	errs := <-g.finishedChan

	if len(errs) == 0 {
		return nil
	}

	return GroupError{Errs: errs}
}
