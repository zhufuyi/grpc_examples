## GRPC example of swagger document generation from protobuf

### Using

**(1) use file server**

start server `go run main.go`

open a browser and type in the address `http://127.0.0.1:8080/swagger/index.html`.

<br>

**(2) use docker**

> docker run -p 8080:8080 -v $PWD/openapiv2/:/openapiv2 -e SWAGGER_JSON=/openapiv2/hello.swagger.json swaggerapi/swagger-ui:v4.14.0

- -v share json file to container/openapiv2 directory
- -e SWAGGER_JSON specify json file

open a browser and type in the address `http://127.0.0.1:8080`.
