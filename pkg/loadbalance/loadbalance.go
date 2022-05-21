package loadbalance

import (
	"fmt"
	"sync"

	"google.golang.org/grpc/resolver"
)

var mux = &sync.Mutex{}

// Register 注册地址到resolver map
func Register(schemeStr string, serviceNameStr string, address []string) string {
	mux.Lock()
	defer mux.Unlock()

	resolver.Register(&ResolverBuilder{
		SchemeVal:   schemeStr,
		ServiceName: serviceNameStr,
		Addrs:       address,
	})

	return fmt.Sprintf("%s:///%s", schemeStr, serviceNameStr)
}

type ResolverBuilder struct {
	SchemeVal   string // SchemeVal作为唯一标致，重复会被覆盖addrs
	ServiceName string
	Addrs       []string
}

func (r *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	blr := &blResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			r.ServiceName: r.Addrs,
		},
	}
	blr.start()
	return blr, nil
}

func (r *ResolverBuilder) Scheme() string {
	return r.SchemeVal
}

type blResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *blResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*blResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*blResolver) Close() {}
