package opencensus

import (
	"fmt"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
)

func init() {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		panic(fmt.Errorf("Failed to register client views for gRPC metrics: %v", err))
	}
	if err := view.Register(ochttp.DefaultClientViews...); err != nil {
		panic(fmt.Errorf("Failed to register client views for HTTP metrics: %v", err))
	}
}
