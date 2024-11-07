package syncgroup_test

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kurt212/syncgroup"
	"github.com/kurt212/syncgroup/internal/testutil"
)

type MyError struct {
	a string
}

func (m MyError) Error() string {
	return m.a
}

func TestGoOK(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.Go(func() error {
		return nil
	})

	syncgrp.Go(func() error {
		return nil
	})

	syncgrp.Go(func() error {
		return nil
	})

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestTryGoOK(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.TryGo(func() error {
		return nil
	})

	syncgrp.TryGo(func() error {
		return nil
	})

	syncgrp.TryGo(func() error {
		return nil
	})

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestGoWithOneError(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.Go(func() error {
		return nil
	})

	returnMyErr := MyError{"123"}

	syncgrp.Go(func() error {
		return returnMyErr
	})

	syncgrp.Go(func() error {
		return nil
	})

	err := syncgrp.Wait()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != returnMyErr.Error() {
		t.Fatalf("expected %v, got %v", returnMyErr, err)
	}

	testutil.True(t, errors.Is(err, returnMyErr), "Result error should be found by errors.Is")
}

func TestTryGoWithOneError(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.Go(func() error {
		return nil
	})

	returnMyErr := MyError{"123"}

	syncgrp.TryGo(func() error {
		return returnMyErr
	})

	syncgrp.TryGo(func() error {
		return nil
	})

	err := syncgrp.Wait()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != returnMyErr.Error() {
		t.Fatalf("expected %v, got %v", returnMyErr, err)
	}

	testutil.True(t, errors.Is(err, returnMyErr), "Result error should be found by errors.Is")
}

//nolint:dupl
func TestGoWithTwoErrors(t *testing.T) {
	t.Parallel()

	firstErr := MyError{"123"}
	secondErr := MyError{"456"}

	syncgrp := syncgroup.New()

	syncgrp.Go(func() error {
		return nil
	})

	syncgrp.Go(func() error {
		return firstErr
	})

	syncgrp.Go(func() error {
		return secondErr
	})

	err := syncgrp.Wait()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	testutil.True(t, errors.Is(err, firstErr), "Result error should be found by errors.Is")
	testutil.True(t, errors.Is(err, secondErr), "Result error should be found by errors.Is")

	unwrappableErr, ok := err.(interface {
		Unwrap() []error
	})

	testutil.True(t, ok, "Result error should be unwrappable and implement Unwrap() []error interface")

	gotErrors := unwrappableErr.Unwrap()

	testutil.True(t, len(gotErrors) == 2, "Result error should contain 2 errors")
}

//nolint:dupl
func TestTryGoWithTwoErrors(t *testing.T) {
	t.Parallel()

	firstErr := MyError{"123"}
	secondErr := MyError{"456"}

	syncgrp := syncgroup.New()

	syncgrp.TryGo(func() error {
		return nil
	})

	syncgrp.TryGo(func() error {
		return firstErr
	})

	syncgrp.TryGo(func() error {
		return secondErr
	})

	err := syncgrp.Wait()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	testutil.True(t, errors.Is(err, firstErr), "Result error should be found by errors.Is")
	testutil.True(t, errors.Is(err, secondErr), "Result error should be found by errors.Is")

	unwrappableErr, ok := err.(interface {
		Unwrap() []error
	})

	testutil.True(t, ok, "Result error should be unwrappable and implement Unwrap() []error interface")

	gotErrors := unwrappableErr.Unwrap()

	testutil.True(t, len(gotErrors) == 2, "Result error should contain 2 errors")
}

func TestNoGoroutines(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestGoRecoversNonErrorPanic(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	panicMsg := "this is message from panic"

	syncgrp.Go(func() error {
		panic(panicMsg)
	})

	err := syncgrp.Wait()

	testutil.True(
		t,
		errors.Is(err, syncgroup.ErrPanicRecovered),
		"On panic should return special panic error",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), panicMsg),
		"Error should contain panic message",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), "goroutine"),
		"Error should contain stack trace",
	)
}

func TestTryGoRecoversNonErrorPanic(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	panicMsg := "this is message from panic"

	syncgrp.TryGo(func() error {
		panic(panicMsg)
	})

	err := syncgrp.Wait()

	testutil.True(
		t,
		errors.Is(err, syncgroup.ErrPanicRecovered),
		"On panic should return special panic error",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), panicMsg),
		"Error should contain panic message",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), "goroutine"),
		"Error should contain stack trace",
	)
}

func TestGoRecoversErrorPanic(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	panicErr := errors.New("this is error from panic") //nolint:err113

	syncgrp.Go(func() error {
		panic(panicErr)
	})

	err := syncgrp.Wait()

	testutil.True(
		t,
		errors.Is(err, syncgroup.ErrPanicRecovered),
		"On panic should return special panic error",
	)

	testutil.True(
		t,
		errors.Is(err, panicErr),
		"Error should wrap panic error",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), "goroutine"),
		"Error should contain stack trace",
	)
}

func TestTryGoRecoversErrorPanic(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	panicErr := errors.New("this is error from panic") //nolint:err113

	syncgrp.TryGo(func() error {
		panic(panicErr)
	})

	err := syncgrp.Wait()

	testutil.True(
		t,
		errors.Is(err, syncgroup.ErrPanicRecovered),
		"On panic should return special panic error",
	)

	testutil.True(
		t,
		errors.Is(err, panicErr),
		"Error should wrap panic error",
	)

	testutil.True(
		t,
		strings.Contains(err.Error(), "goroutine"),
		"Error should contain stack trace",
	)
}

func TestLimitGoroutinesWithGo(t *testing.T) {
	t.Parallel()

	const limit = 2

	syncgrp := syncgroup.New()
	syncgrp.SetLimit(limit)

	activeCount := atomic.Int32{}

	runnableFunc := func() error {
		active := activeCount.Add(1)
		defer activeCount.Add(-1)

		if active > limit {
			t.Errorf("expected %d active goroutines, got %d", limit, active)
		}

		time.Sleep(100 * time.Millisecond)

		return nil
	}

	const goroutinesCount = 10

	for range goroutinesCount {
		syncgrp.Go(runnableFunc)
	}

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLimitGoroutinesWithTryGo(t *testing.T) {
	t.Parallel()

	const limit = 2

	syncgrp := syncgroup.New()
	syncgrp.SetLimit(limit)

	activeCount := atomic.Int32{}

	runnableFunc := func() error {
		active := activeCount.Add(1)
		defer activeCount.Add(-1)

		if active > limit {
			t.Errorf("expected %d active goroutines, got %d", limit, active)
		}

		time.Sleep(100 * time.Millisecond)

		return nil
	}

	const goroutinesCount = 10

	for range goroutinesCount {
		syncgrp.TryGo(runnableFunc)
	}

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestTryGoReturnsValidValue(t *testing.T) {
	t.Parallel()

	const limit = 2

	syncgrp := syncgroup.New()
	syncgrp.SetLimit(limit)

	stopChan := make(chan struct{})

	runnableFunc := func() error {
		<-stopChan

		return nil
	}

	const goroutinesCount = 10

	failedToRun := 0

	for range goroutinesCount {
		ok := syncgrp.TryGo(runnableFunc)
		if !ok {
			failedToRun++
		}
	}

	testutil.Equal(t, failedToRun, goroutinesCount-limit)

	close(stopChan)

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCanNotChangeLimitAfterGo(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	stopChan := make(chan struct{})

	syncgrp.Go(func() error {
		<-stopChan

		return nil
	})

	testutil.Panics(t, func() {
		syncgrp.SetLimit(1)
	})

	close(stopChan)

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCanNotChangeLimitAfterTryGo(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	stopChan := make(chan struct{})

	syncgrp.TryGo(func() error {
		<-stopChan

		return nil
	})

	testutil.Panics(t, func() {
		syncgrp.SetLimit(1)
	})

	close(stopChan)

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUnsetLimitWorksWithGo(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.SetLimit(1)
	syncgrp.SetLimit(0)

	stopChan := make(chan struct{})

	activeGoroutines := atomic.Int32{}

	const goroutinesCount = 10

	gouroutineFunc := func() error {
		activeGoroutines.Add(1)
		defer activeGoroutines.Add(-1)

		<-stopChan

		return nil
	}

	for range goroutinesCount {
		syncgrp.Go(gouroutineFunc)
	}

	for activeGoroutines.Load() != goroutinesCount {
		time.Sleep(10 * time.Millisecond)
	}

	close(stopChan)

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUnsetLimitWorksWithTryGo(t *testing.T) {
	t.Parallel()

	syncgrp := syncgroup.New()

	syncgrp.SetLimit(1)
	syncgrp.SetLimit(0)

	stopChan := make(chan struct{})

	activeGoroutines := atomic.Int32{}

	const goroutinesCount = 10

	goroutineFunc := func() error {
		activeGoroutines.Add(1)
		defer activeGoroutines.Add(-1)

		<-stopChan

		return nil
	}

	for range goroutinesCount {
		syncgrp.TryGo(goroutineFunc)
	}

	for activeGoroutines.Load() != goroutinesCount {
		time.Sleep(10 * time.Millisecond)
	}

	close(stopChan)

	err := syncgrp.Wait()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
