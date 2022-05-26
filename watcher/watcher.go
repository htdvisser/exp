// Package watcher provides a simple interface for watching for changes to a value.
package watcher

import (
	"context"
	"sync"

	"github.com/zyedidia/generic/list"
)

// Notifier[X] is an interface for getting notified about value changes.
type Notifier[X any] interface {
	Notify(X)
}

// Func is a Notifier[X] that calls a function when the value changes.
type Func[X any] func(X)

// Notify implements the Notifier[X] interface.
func (w Func[X]) Notify(x X) { w(x) }

// Channel[X] is a Notifier[X] that sends value changes to a channel.
type Channel[X any] chan X

// Notify implements the Notifier[X] interface.
// If the channel is blocked (full), changed values are dropped.
func (c Channel[X]) Notify(x X) {
	select {
	case c <- x:
	default:
	}
}

// Value[X] is a value that can be watched for changes.
type Value[X any] struct {
	mu           sync.Mutex
	currentValue X
	equals       func(X, X) bool
	watchers     list.List[Notifier[X]]
}

// NewComparableValue creates a new Value[X] for types that are comparable with the given initial value.
func NewComparableValue[X comparable](initialValue X) *Value[X] {
	return NewValue(initialValue, func(a X, b X) bool { return a == b })
}

// NewValue creates a new Value[X] with the given initial value.
func NewValue[X any](initialValue X, equals func(X, X) bool) *Value[X] {
	return &Value[X]{
		equals:       equals,
		currentValue: initialValue,
	}
}

// Set sets the value of the Value[X] and notifies all watchers if the new value is different than the old one.
// Set does not return until all watchers have been notified.
func (v *Value[X]) Set(newValue X) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.equals(newValue, v.currentValue) {
		v.currentValue = newValue
		v.watchers.Front.Each(func(watcher Notifier[X]) {
			watcher.Notify(newValue)
		})
	}
}

// Watch adds a new watcher to the Value[X] and notifies the watcher with the current value.
// It returns a function that can be used to remove the watcher.
func (v *Value[X]) Watch(watcher Notifier[X]) (unwatch func()) {
	v.mu.Lock()
	defer v.mu.Unlock()
	node := &list.Node[Notifier[X]]{Value: watcher}
	v.watchers.PushBackNode(node)
	watcher.Notify(v.currentValue)
	return func() {
		v.mu.Lock()
		defer v.mu.Unlock()
		v.watchers.Remove(node)
	}
}

// WatchFunc adds a new watcher func to the Value[X] and calls the given function with the current value.
func (v *Value[X]) WatchFunc(f func(X)) (unwatch func()) {
	return v.Watch(Func[X](f))
}

// WaitForChange waits for the value to change to a value different than sourceValue.
// If the value is not changed before the context is done, an error is returned.
func (v *Value[X]) WaitForChange(ctx context.Context, sourceValue X) (X, error) {
	v.mu.Lock()
	currentValue := v.currentValue
	if !v.equals(currentValue, sourceValue) {
		v.mu.Unlock()
		return currentValue, nil
	}
	watcher := make(Channel[X], 1)
	node := &list.Node[Notifier[X]]{Value: watcher}
	v.watchers.PushBackNode(node)
	v.mu.Unlock()
	defer func() {
		v.mu.Lock()
		defer v.mu.Unlock()
		v.watchers.Remove(node)
	}()
	select {
	case <-ctx.Done():
		return sourceValue, ctx.Err()
	case newValue := <-watcher:
		return newValue, nil
	}
}
