// Package syncgroup package that contains an implementation of an abstract
// synchronisation mechanism - synchronisation group.
// The main idea is to have an ability to run independent tasks in separate goroutines which way return errors.
// A user can wait until all goroutines finish running and collect all occurred errors.
//
// The design is similar to errgroup (https://godoc.org/golang.org/x/sync/errgroup),
// but it does not cancel the context of the goroutines if any of them returns an error.
package syncgroup

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// ErrPanicRecovered is a special error that is returned when a panic is recovered from a goroutine.
// It is used to wrap the original panic error and stack trace.
// You can use errors.Is(err, ErrPanicRecovered) to check if the error was caused by a panic.
// If the panic value was an error, you can use errors.Unwrap(err) to get the original error.
var ErrPanicRecovered = errors.New("recovered from panic")

// SyncGroup is the main class for working with syncgroups. It's a collection of goroutines that can be waited for.
// Additionally, SyncGroup collects all errors returned by goroutines,
// handles panics and provides a way to limit the number of concurrent goroutines.
// It has two main methods: Go() and Wait()
//
// Go() spawns a new goroutine, which may return an error.
// The returned error will be saved and returned by Wait() method.
//
// Wait() waits until all spawned goroutines finish and returns a wrapper for a slice of errors.
// If there was no error, Wait() would return nil,
// otherwise a non nil error, which can be unwrapped to access all errors.
type SyncGroup struct {
	wg        sync.WaitGroup
	semaphore chan semaphoreToken

	finishedChan chan []error
	errorChan    chan error

	listeningStarted        atomic.Bool
	listeningRoutineStarter *sync.Once
}

type semaphoreToken struct{}

// New is the default constructor for SyncGroup.
func New() *SyncGroup {
	grp := &SyncGroup{
		wg:                      sync.WaitGroup{},
		semaphore:               nil,
		finishedChan:            make(chan []error),
		errorChan:               make(chan error),
		listeningStarted:        atomic.Bool{},
		listeningRoutineStarter: new(sync.Once),
	}

	return grp
}

// Go spawns given function in a new goroutine.
// If group has a limit of concurrent goroutines, goroutine execution will be blocked until a slot is available.
// The returned error will be saved and returned wrapped by Wait() method.
func (g *SyncGroup) Go(fnc func() error) {
	g.startListening()

	g.wg.Add(1)

	go func() {
		defer g.done()

		// blocks until semaphore slot is acquired
		if g.semaphore != nil {
			g.semaphore <- semaphoreToken{}
		}

		err := fnc()
		if err != nil {
			g.errorChan <- err
		}
	}()
}

func (g *SyncGroup) TryGo(fnc func() error) bool {
	if g.semaphore != nil {
		select {
		case g.semaphore <- semaphoreToken{}:
		default:
			return false
		}
	}

	g.startListening()
	g.wg.Add(1)

	go func() {
		defer g.done()

		err := fnc()
		if err != nil {
			g.errorChan <- err
		}
	}()

	return true
}

// done is called in every goroutine spawned by SyncGroup in defer statement.
// Its job is to handle panics, release all resources and decrement the WaitGroup counter.
func (g *SyncGroup) done() {
	if msg := recover(); msg != nil {
		var err error

		switch val := msg.(type) {
		case error:
			err = fmt.Errorf("%w: %w\n%s", ErrPanicRecovered, val, string(debug.Stack()))
		default:
			err = fmt.Errorf("%w: %v\n%s", ErrPanicRecovered, val, string(debug.Stack()))
		}

		g.errorChan <- err
	}

	if g.semaphore != nil {
		<-g.semaphore
	}

	g.wg.Done()
}

func (g *SyncGroup) startListening() {
	g.listeningRoutineStarter.Do(func() {
		g.listeningStarted.Store(true)
		go g.listenToErrors()
	})
}

// listenToErrors is a single per group goroutine that listens to all errors and accumulates them.
func (g *SyncGroup) listenToErrors() {
	defer func() {
		close(g.finishedChan)
	}()

	var accumulatedErrors []error //nolint:prealloc // false positive
	for err := range g.errorChan {
		accumulatedErrors = append(accumulatedErrors, err)
	}

	g.finishedChan <- accumulatedErrors
}

// Wait waits until all spawned goroutines are finished and returns a wrapped error for all collected errors.
// The result is nil if none of the spawned goroutines returned an error
//
// If error is not nil, the result is guaranteed to implement `Unwrap() []errors` methods to access all errors.
// The error supports unwrapping with standard errors.Unwrap(), errors.Is() and errors.As() functions.
func (g *SyncGroup) Wait() error {
	if !g.listeningStarted.Load() {
		return nil
	}

	g.wg.Wait()
	close(g.errorChan)

	errs := <-g.finishedChan

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func (g *SyncGroup) SetLimit(limit int) {
	if g.listeningStarted.Load() {
		panic("cannot set limit after starting goroutines")
	}

	if limit <= 0 {
		g.semaphore = nil

		return
	}

	g.semaphore = make(chan semaphoreToken, limit)
}
