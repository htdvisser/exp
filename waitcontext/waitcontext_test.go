package waitcontext_test

import (
	"context"
	"fmt"
	"time"

	"htdvisser.dev/exp/waitcontext"
)

func Example() {
	// Create a context that waits after cancel:
	ctx, cancelAndWait := waitcontext.New(context.Background())

	// Start a goroutine:
	waitcontext.Go(ctx, func() {
		<-ctx.Done()
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Really done")
	})

	// You could extend the context elsewhere in your program:
	ctx, _ = waitcontext.New(ctx)

	// Start another goroutine:
	waitcontext.Go(ctx, func() {
		<-ctx.Done()
		time.Sleep(time.Millisecond)
		fmt.Println("Really done")
	})

	// Cancel the context and wait until all goroutines started in the context (or derived contexts) are done:
	cancelAndWait()

	// Output:
	// Really done
	// Really done
}
