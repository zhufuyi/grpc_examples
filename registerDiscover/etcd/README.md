## GRPC Registration and Discovery Example

### Using

Script to run the etcd service docker-compose.yml

```yaml
version: "3"

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.4
    container_name: etcd-standalone
    restart: always
    environment:
      - ETCDCTL_API=3
    command:
      - etcd
      - --name=etcd-single
      - --data-dir=/etcd-data
      - --listen-client-urls=http://0.0.0.0:2379
      - --advertise-client-urls=http://etcd:2379
      - --listen-peer-urls=http://0.0.0.0:2380
      - --initial-advertise-peer-urls=http://etcd:2380
      #- --initial-cluster=etcd-single=http://0.0.0.0:2380
    ports:
      - 2379:2379
      - 2380:2380
    volumes:
      - $PWD/etcd-data:/etcd-data
```

Running the standalone etcd service

> docker-compose up -d

<br>

 Start the rpc service and register the service address with etcd

> cd server && go run main.go

<br>

Start the rpc client and discover the rpc service port address by service name

> cd client && go run main.go
