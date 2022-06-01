package discovery

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

const (
	schema = "etcd"
)

// Resolver for grpc client
type Resolver struct {
	schema      string
	EtcdAddrs   []string
	DialTimeout int

	closeCh      chan struct{}
	watchCh      clientv3.WatchChan
	cli          *clientv3.Client
	keyPrifix    string
	srvAddrsList []resolver.Address

	cc     resolver.ClientConn
	logger *zap.Logger
}

// NewResolver create a new resolver.Builder base on etcd
func NewResolver(etcdAddrs []string, opts ...Option) *Resolver {
	o := defaultOptions()
	o.apply(opts...)

	return &Resolver{
		schema:      schema,
		EtcdAddrs:   etcdAddrs,
		DialTimeout: o.etcdDialTimeoutSeconds,
		logger:      o.logger,
	}
}

// Scheme returns the scheme supported by this resolver.
func (r *Resolver) Scheme() string {
	return r.schema
}

// Build creates a new resolver.Resolver for the given target
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc

	r.keyPrifix = buildPrefix(server{Name: target.Endpoint, Version: target.Authority})
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// ResolveNow resolver.Resolver interface
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
	r.closeCh <- struct{}{}
}

// start
func (r *Resolver) start() (chan<- struct{}, error) {
	var err error
	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	resolver.Register(r)

	r.closeCh = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch()

	return r.closeCh, nil
}

// watch update events
func (r *Resolver) watch() {
	ticker := time.NewTicker(time.Minute)
	r.watchCh = r.cli.Watch(context.Background(), r.keyPrifix, clientv3.WithPrefix())

	for {
		select {
		case <-r.closeCh:
			return
		case res, ok := <-r.watchCh:
			if ok {
				r.update(res.Events)
			}
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", zap.Error(err))
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) {
	for _, ev := range events {
		var info server
		var err error

		switch ev.Type {
		case mvccpb.PUT:
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				continue
			}
			addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
			if !exist(r.srvAddrsList, addr) {
				r.srvAddrsList = append(r.srvAddrsList, addr)
				err := r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
				if err != nil {
					r.logger.Error("UpdateState error", zap.Error(err), zap.Any("addr", addr))
				}
			}

		case mvccpb.DELETE:
			info, err = splitPath(string(ev.Kv.Key))
			if err != nil {
				r.logger.Error("splitPath error", zap.Error(err), zap.String("key", string(ev.Kv.Key)))
				continue
			}
			addr := resolver.Address{Addr: info.Addr}
			if s, ok := remove(r.srvAddrsList, addr); ok {
				r.srvAddrsList = s
				err := r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
				if err != nil {
					r.logger.Error("UpdateState error", zap.Error(err), zap.Any("addr", addr))
				}
			}
		}
	}
}

// sync Simultaneous access to all address information
func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.cli.Get(ctx, r.keyPrifix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	r.srvAddrsList = []resolver.Address{}
	for _, v := range res.Kvs {
		info, err := ParseValue(v.Value)
		if err != nil {
			r.logger.Error("ParseValue error", zap.Error(err), zap.String("value", string(v.Value)))
			continue
		}
		addr := resolver.Address{Addr: info.Addr, Metadata: info.Weight}
		r.srvAddrsList = append(r.srvAddrsList, addr)
	}

	err = r.cc.UpdateState(resolver.State{Addresses: r.srvAddrsList})
	return err
}
