package syncgroup

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

type SyncGroup struct {
	wg sync.WaitGroup

	finishedChan chan []error
	errorChan    chan error

	listeningStarted bool
}

type GroupError struct {
	Errs []error
}

func (e GroupError) Error() string {
	var accumulated string
	for _, err := range e.Errs {
		accumulated += err.Error() + ";"
	}

	return accumulated
}

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

func (g *SyncGroup) Go(f func() error) {
	if !g.listeningStarted {
		go g.listenToErrors()
		g.listeningStarted = true
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		defer func() {
			if msg := recover(); msg != nil {
				switch msg.(type) {
				case error:
					g.errorChan <- errors.Wrap(msg.(error), "recovered from panic")
				case fmt.Stringer:
					g.errorChan <- errors.Errorf("recovered from panic:%s", msg.(fmt.Stringer).String())
				default:
					g.errorChan <- errors.Errorf("recovered from panic:%v", msg)
				}
			}
		}()

		err := f()

		if err != nil {
			g.errorChan <- err
		}
	}()
}

func (g *SyncGroup) Wait() error {
	g.wg.Wait()
	close(g.errorChan)

	errs := <-g.finishedChan

	if len(errs) == 0 {
		return nil
	}

	return GroupError{Errs: errs}
}
