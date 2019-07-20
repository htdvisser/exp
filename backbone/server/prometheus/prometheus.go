// Package prometheus can be used to export opencensus metrics to prometheus.
package prometheus

import (
	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats/view"
	"htdvisser.dev/exp/backbone/server"
)

type options struct {
	namespace string
	registry  *prometheus.Registry
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

// WithNamespace sets the namespace of the Prometheus exporter.
func WithNamespace(namespace string) Option {
	return option(func(opts *options) {
		opts.namespace = namespace
	})
}

// WithRegistry sets the registry to be used by the Prometheus exporter.
func WithRegistry(registry *prometheus.Registry) Option {
	return option(func(opts *options) {
		opts.registry = registry
	})
}

var defaultRegistry = prometheus.NewRegistry()

func init() {
	defaultRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	defaultRegistry.MustRegister(prometheus.NewGoCollector())
}

// Register registers a Prometheus exporter to the server.
func Register(s *server.Server, opts ...Option) error {
	options := &options{
		registry: defaultRegistry,
	}
	options.apply(opts...)
	pe, err := ocprom.NewExporter(ocprom.Options{
		Namespace: options.namespace,
		Registry:  options.registry,
	})
	if err != nil {
		return err
	}
	view.RegisterExporter(pe)
	s.InternalHTTP.Handle("/metrics", pe)
	return nil
}
