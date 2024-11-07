package syncgroup

import (
	"testing"

	"github.com/kurt212/syncgroup/internal/testutil"
)

type MyError struct {
	a string
}

func (m MyError) Error() string {
	return m.a
}

func TestListenTo(t *testing.T) {
	t.Parallel()

	syncgrp := New()

	syncgrp.startListening()

	syncgrp.errorChan <- MyError{"err1"}
	syncgrp.errorChan <- MyError{"err2"}
	syncgrp.errorChan <- MyError{"err3"}

	close(syncgrp.errorChan)

	res := <-syncgrp.finishedChan

	expected := []error{
		MyError{"err1"},
		MyError{"err2"},
		MyError{"err3"},
	}

	testutil.EqualSlices(t, expected, res)
}
