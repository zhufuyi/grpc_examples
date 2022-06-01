package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo server
	srvTTL  int64 // unit is second
	cli     *clientv3.Client
	logger  *zap.Logger
}

// RegisterRPCAddr Register the rpc service address to etcd
func RegisterRPCAddr(serverName string, grpcAddr string, etcdAddrs []string, opts ...Option) *Register {
	o := defaultOptions()
	o.apply(opts...)

	if strings.Split(grpcAddr, ":")[0] == "" {
		panic("invalid grpc address " + grpcAddr)
	}

	//etcdRegister := newRegister(etcdAddrs, o.etcdDialTimeoutSeconds, o.logger)
	srvInfo := server{
		Name:    serverName,
		Addr:    grpcAddr,
		Version: o.version,
		Weight:  o.weight,
	}

	etcdRegister := &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: o.etcdDialTimeoutSeconds,
		logger:      o.logger,
		closeCh:     make(chan struct{}),
		srvInfo:     srvInfo,
		srvTTL:      o.ttlSeconds,
	}

	if err := etcdRegister.registerServer(); err != nil {
		panic("server register failed: " + err.Error())
	}

	return etcdRegister
}

// Register a service
func (r *Register) registerServer() error {
	var err error

	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return err
	}

	if err = r.register(); err != nil {
		return err
	}

	go r.keepAlive()

	return nil
}

// Stop stop register
func (r *Register) Stop() {
	defer func() {
		if err := recover(); err != nil {
			r.logger.Warn("repeated stopped")
		}
	}()

	r.closeCh <- struct{}{}
	time.Sleep(time.Millisecond * 100) // Allow a little time before exiting the service to give the goroutine a chance to unregister node
	close(r.closeCh)
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}

	_, err = r.cli.Put(context.Background(), buildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), buildRegPath(r.srvInfo))
	return err
}

// keepAlive
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}

		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}

		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			} else {
				r.logger.Info("unregister server ok", zap.String("server", buildRegPath(r.srvInfo)))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		}
	}
}

// UpdateHandler return http handler
func (r *Register) UpdateHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(context.Background(), buildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	})
}

func (r *Register) GetServerInfo() (server, error) {
	resp, err := r.cli.Get(context.Background(), buildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}
