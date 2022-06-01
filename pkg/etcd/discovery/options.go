package discovery

import (
	"go.uber.org/zap"
)

var (
	defaultTTLSeconds             int64 = 10
	defaultVersion                      = ""
	defaultWeight                 int64 = 0
	defaultLogger                       = zap.NewExample()
	defaultEtcdDialTimeoutSeconds       = 3
)

type options struct {
	// register options
	ttlSeconds int64
	version    string
	weight     int64

	// register or resolver options
	logger                 *zap.Logger
	etcdDialTimeoutSeconds int
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option represents the etcd options
type Option func(*options)

// WithTTLSeconds sets timeout for requesting the etcd service in seconds
func WithTTLSeconds(ttlSeconds int64) Option {
	return func(o *options) {
		o.ttlSeconds = ttlSeconds
	}
}

// WithLogger sets print log
func WithLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithVersion sets etcd value version
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// WithWeight sets etcd weight, use in loadbalance
func WithWeight(weight int64) Option {
	return func(o *options) {
		o.weight = weight
	}
}

// WithDialTimeout sets rpc dial timeout
func WithDialTimeout(dialTimeoutSeconds int) Option {
	return func(o *options) {
		o.etcdDialTimeoutSeconds = dialTimeoutSeconds
	}
}

func defaultOptions() *options {
	return &options{
		ttlSeconds:             defaultTTLSeconds,
		version:                defaultVersion,
		weight:                 defaultWeight,
		logger:                 defaultLogger,
		etcdDialTimeoutSeconds: defaultEtcdDialTimeoutSeconds,
	}
}
