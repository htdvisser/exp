// Package jaeger can be used to export opentelemetry traces to jaeger.
package jaeger

import (
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/exporter/trace/jaeger"
)

type options struct {
	endpointOption jaeger.EndpointOption
	serviceName    string
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt.apply(o)
	}
}

// Option is an option for the Prometheus exporter.
type Option interface {
	apply(*options)
}

type option func(*options)

func (f option) apply(opts *options) {
	f(opts)
}

// WithCollectorEndpoint configures the Jaeger exporter with the given collector endpoint.
func WithCollectorEndpoint(collectorEndpoint string) Option {
	return option(func(opts *options) {
		opts.endpointOption = jaeger.WithCollectorEndpoint(collectorEndpoint)
	})
}

// WithAuthenticatedCollectorEndpoint configures the Jaeger exporter with the given collector endpoint.
func WithAuthenticatedCollectorEndpoint(collectorEndpoint, username, password string) Option {
	return option(func(opts *options) {
		opts.endpointOption = jaeger.WithCollectorEndpoint(
			collectorEndpoint,
			jaeger.WithUsername(username),
			jaeger.WithPassword(password),
		)
	})
}

// WithAgentEndpoint configures the Jaeger exporter with the given agent endpoint.
func WithAgentEndpoint(agentEndpoint string) Option {
	return option(func(opts *options) {
		opts.endpointOption = jaeger.WithAgentEndpoint(agentEndpoint)
	})
}

// WithServiceName sets the service name of the Jaeger exporter.
func WithServiceName(serviceName string) Option {
	return option(func(opts *options) {
		opts.serviceName = serviceName
	})
}

// NewExporter returns a Jaeger exporter that can be attached to the server.
func NewExporter(opts ...Option) (*jaeger.Exporter, error) {
	options := &options{
		endpointOption: jaeger.WithCollectorEndpoint("http://jaeger-collector.istio-system.svc.cluster.local:14268/api/traces"),
		serviceName:    filepath.Base(os.Args[0]),
	}
	options.apply(opts...)
	return jaeger.NewExporter(
		options.endpointOption,
		jaeger.WithProcess(jaeger.Process{
			ServiceName: options.serviceName,
		}),
	)
}
