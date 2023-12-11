## GRPC Registration and Discovery Example

### Run the etcd service

Run the etcd service with the script docker-compose.yml

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

  etcd-ui:
    image: deltaprojects/etcdkeeper:latest
    container_name: etcd-ui
    restart: always
    ports:
      - "12379:8080"
    environment:
      # 不是etcd地址
      HOST: "0.0.0.0"
    depends_on:
      - etcd
```

Running the standalone etcd service

> docker-compose up -d

<br>

### Run the grpc service

Start the grpc service and register the service address with etcd

```bash
cd server
go build

# running 2 grpc server
./server -p 8282
./server -p 8284
```

You can see in etcd that there are two services kv

```
/microservices/helloDemo/helloDemo_grpc_127.0.0.1_8282
{"id":"helloDemo_grpc_127.0.0.1_8282","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://127.0.0.1:8282"]}

/microservices/helloDemo/helloDemo_grpc_127.0.0.1_8284
{"id":"helloDemo_grpc_127.0.0.1_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://127.0.0.1:8284"]}
```

<br>

### Run the grpc client

Start the grpc client and discover the grpc service address by service name.

```bash
cd client
go build

./client
```

You can see the polling requesting two services.

```
hello foo, (from 127.0.0.1:8284)
hello foo, (from 127.0.0.1:8282)
hello foo, (from 127.0.0.1:8284)
hello foo, (from 127.0.0.1:8282)
......
```

<br>

### Test service registration and discovery

Test stopping and starting the grpc service to check that service registration and discovery is dynamically aware.

1. Stop the grpc service on port 8282

You can see that the grpc service on port 8282 is detected as offline, and only the service on port 8284 is shown as responsive on the client.

```
[resolver] update instances: [{"id":"helloDemo_grpc_127.0.0.1_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://127.0.0.1:8284"]}]
```

2. Start the gprc service on port 8282

```bash
./server -p 8282
```

You can see that two services are detected as running, and it reverts back to polling for grpc services.

```
[resolver] update instances: [{"id":"helloDemo_grpc_127.0.0.1_8282","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://127.0.0.1:8282"]},{"id":"helloDemo_grpc_127.0.0.1_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://127.0.0.1:8284"]}]
```
