package syncgroup

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

type SyncGroup struct {
	wg sync.WaitGroup

	finishedChan chan string
	errorChan    chan error
}

func New() *SyncGroup {
	g := &SyncGroup{
		wg:           sync.WaitGroup{},
		finishedChan: make(chan string),
		errorChan:    make(chan error),
	}

	go g.listenToErrors()

	return g
}

func (g *SyncGroup) listenToErrors() {
	var accumulatedErrorMsg string

	for {
		err, ok := <-g.errorChan
		if ok {
			accumulatedErrorMsg += err.Error() + ";"
		} else {
			break
		}
	}

	g.finishedChan <- accumulatedErrorMsg
	close(g.finishedChan)

	return
}

func (g *SyncGroup) Go(f func() error) {
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

	errMsg := <-g.finishedChan

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
}
