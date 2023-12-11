## GRPC Registration and Discovery Example

### Run the consul service

Run the consul service with the script docker-compose.yml

```yaml
version: '3.7'

services:

  consul-server:
    image: hashicorp/consul:1.12.0
    container_name: consul-server
    restart: always
    volumes:
      - $PWD/server.json:/consul/config/server.json
      - $PWD/data/server:/consul/data
    networks:
      - consul
    ports:
      - "8500:8500"
      - "8600:8600/tcp"
      - "8600:8600/udp"
    command: "agent"

  consul-client:
    image: hashicorp/consul:1.12.0
    container_name: consul-client
    restart: always
    volumes:
      - $PWD/client.json:/consul/config/client.json
      - $PWD/data/client:/consul/data
    networks:
      - consul
    command: "agent"

networks:
  consul:
    driver: bridge
```

Running the standalone consul service

> docker-compose up -d

<br>

### Run the grpc service

Start the grpc service and register the service address with consul

```bash
cd server
go build

# running 2 grpc server
# If the consul service is deployed on another machine, the startup service configuration
# must be filled in with the specified real ip, and 127.0.0.1 cannot be used.
./server -i 192.168.3.79 -p 8282
./server -i 192.168.3.79 -p 8284
```

You can see in consul that there are two services kv.

```
{"id":"helloDemo_grpc_192.168.3.79_8282","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://192.168.3.79:8282"]}

{"id":"helloDemo_grpc_192.168.3.79_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://192.168.3.79:8284"]}
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
[resolver] update instances: [{"id":"helloDemo_grpc_192.168.3.79_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://192.168.3.79:8284"]}]
```

2. Start the gprc service on port 8282

```bash
./server -p 8282
```

You can see that two services are detected as running, and it reverts back to polling for grpc services.

```
[resolver] update instances: [{"id":"helloDemo_grpc_192.168.3.79_8282","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://192.168.3.79:8282"]},{"id":"helloDemo_grpc_192.168.3.79_8284","name":"helloDemo","version":"","metadata":null,"endpoints":["grpc://192.168.3.79:8284"]}]
```
