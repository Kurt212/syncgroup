package syncgroup

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type MyErr struct {
	a string
}

func (m MyErr) Error() string {
	return m.a
}

func TestListenTo(t *testing.T) {
	as := assert.New(t)

	sg := &SyncGroup{
		wg:           sync.WaitGroup{},
		finishedChan: make(chan []error),
		errorChan:    make(chan error),
	}

	go sg.listenToErrors()

	expected := []error{MyErr{"err1"}, MyErr{"err2"}, MyErr{"err3"}}

	sg.errorChan <- MyErr{"err1"}
	sg.errorChan <- MyErr{"err2"}
	sg.errorChan <- MyErr{"err3"}

	close(sg.errorChan)

	res := <-sg.finishedChan

	as.Equal(expected, res)
}

func TestSyncOK(t *testing.T) {
	as := assert.New(t)

	sg := New()

	sg.Go(func() error {
		return nil
	})

	sg.Go(func() error {
		return nil
	})

	sg.Go(func() error {
		return nil
	})

	err := sg.Wait()
	as.Nil(err)
}

func TestSyncBad1(t *testing.T) {
	as := assert.New(t)

	sg := New()

	sg.Go(func() error {
		return nil
	})

	sg.Go(func() error {
		return MyErr{"123"}
	})

	sg.Go(func() error {
		return nil
	})

	err := sg.Wait()

	as.NotNil(err)
	as.Equal("123;", err.Error())
}

func TestSyncBad2(t *testing.T) {
	as := assert.New(t)

	sg := New()

	sg.Go(func() error {
		return nil
	})

	sg.Go(func() error {
		return MyErr{"123"}
	})

	sg.Go(func() error {
		return MyErr{"456"}
	})

	err := sg.Wait()

	as.NotNil(err)
	as.Contains([]string{"123;456;", "456;123;"}, err.Error())
}
