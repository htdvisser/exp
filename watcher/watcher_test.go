package watcher_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"htdvisser.dev/exp/watcher"
)

func ExampleValue_WatchFunc() {
	v := watcher.NewComparableValue(0)

	v.WatchFunc(func(i int) {
		fmt.Printf("The value is now %d\n", i)
	})

	v.Set(1)
	v.Set(2)
	v.Set(3)

	// Output:
	// The value is now 0
	// The value is now 1
	// The value is now 2
	// The value is now 3
}

func TestWatch(t *testing.T) {
	v := watcher.NewComparableValue("")
	v.Set("foo")

	ch := make(watcher.Channel[string], 2)
	unwatch := v.Watch(ch)

	if len(ch) != 1 {
		t.Error("Expected channel to contain one element")
	}
	if <-ch != "foo" {
		t.Error("Expected channel to contain foo")
	}

	v.Set("bar")
	v.Set("baz")
	v.Set("qux") // This value is dropped.

	if len(ch) != 2 {
		t.Error("Expected channel to contain two elements")
	}
	if <-ch != "bar" {
		t.Error("Expected channel to contain bar")
	}
	if <-ch != "baz" {
		t.Error("Expected channel to contain bar")
	}

	v.Set("qux") // Even though we haven't read the qux change, this is not a change.

	if len(ch) > 0 {
		t.Error("Expected channel to be empty")
	}

	unwatch()

	v.Set("quux")

	if len(ch) > 0 {
		t.Error("Expected channel to be empty")
	}
}

func TestWaitForChange(t *testing.T) {
	v := watcher.NewComparableValue("")
	v.Set("foo")

	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	t.Run("already changed", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		changed, err := v.WaitForChange(ctx, "")
		if err != nil {
			t.Error("Expected no error")
		}
		if changed != "foo" {
			t.Error("Expected changed value to be foo")
		}
	})

	t.Run("changing", func(t *testing.T) {
		go func() {
			time.Sleep(100 * time.Millisecond)
			v.Set("bar")
		}()

		changed, err := v.WaitForChange(ctx, "foo")
		if err != nil {
			t.Error("Expected no error")
		}
		if changed != "bar" {
			t.Error("Expected changed value to be bar")
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()

		_, err := v.WaitForChange(ctx, "bar")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Error("Expected context to be canceled")
		}
	})
}
