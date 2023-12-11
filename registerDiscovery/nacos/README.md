## GRPC Registration and Discovery Example

### Run the nacos service

Run the nacos service with the script docker-compose.yml

```yaml
version: "3"
services:
  nacos:
    image: nacos/nacos-server:v2.1.1
    container_name: nacos-standalone
    restart: always
    environment:
      - PREFER_HOST_MODE=hostname
      - MODE=standalone
    volumes:
      - ./data:/home/nacos/data
      - ./standalone-logs/:/home/nacos/logs
      #- ./init.d/custom.properties:/home/nacos/init.d/custom.properties
    ports:
      - "8848:8848"
      - "9848:9848"
```

Running the standalone nacos service

> docker-compose up -d

<br>

### Run the grpc service

Start the grpc service and register the service address with nacos

```bash
cd server
go build

# running 2 grpc server
./server -i 192.168.3.37 -p 8282
./server -i 192.168.3.37 -p 8284
```

You can see in nacos that there are two services kv

```
{"id":"helloDemo_grpc_192.168.3.79_8282","name":"helloDemo.grpc","version":"","metadata":{"kind":"grpc","version":""},"endpoints":["grpc://192.168.3.79:8282"]}

{"id":"helloDemo_grpc_192.168.3.79_8284","name":"helloDemo.grpc","version":"","metadata":{"kind":"grpc","version":""},"endpoints":["grpc://192.168.3.79:8284"]}
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
hello foo, (from 192.168.3.79:8284)
hello foo, (from 192.168.3.79:8282)
hello foo, (from 192.168.3.79:8284)
hello foo, (from 192.168.3.79:8282)
......
```

<br>

### Test service registration and discovery

Test stopping and starting the grpc service to check that service registration and discovery is dynamically aware.

1. Stop the grpc service on port 8284

You can see that the grpc service on port 8282 is detected as offline, and only the service on port 8284 is shown as responsive on the client.

```
[resolver] update instances:  [{"id":"helloDemo_grpc_192.168.3.79_8282","name":"helloDemo.grpc","version":"","metadata":{"kind":"grpc","version":""},"endpoints":["grpc://192.168.3.79:8282"]}]
```

2. Start the gprc service on port 8282

```bash
./server -i 192.168.3.37 -p 8284
```

You can see that two services are detected as running, and it reverts back to polling for grpc services.

```
[resolver] update instances: [{"id":"helloDemo_grpc_192.168.3.79_8284","name":"helloDemo.grpc","version":"","metadata":{"kind":"grpc","version":""},"endpoints":["grpc://192.168.3.79:8284"]},{"id":"helloDemo_grpc_192.168.3.79_8282","name":"helloDemo.grpc","version":"","metadata":{"kind":"grpc","version":""},"endpoints":["grpc://192.168.3.79:8282"]}]
```
