module github.com/zhufuyi/grpc_examples

go 1.19

require (
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/gin-gonic/gin v1.8.1
	github.com/gogo/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.1
	github.com/prometheus/client_golang v1.13.0
	github.com/reugn/equalizer v0.0.0-20210216135016-a959c509d7ad
	github.com/swaggo/files v0.0.0-20220728132757-551d4a08d97a
	github.com/swaggo/gin-swagger v1.5.2
	github.com/swaggo/swag v1.8.4
	github.com/zhufuyi/pkg v1.3.5
	go.etcd.io/etcd/client/v3 v3.5.4
	go.opentelemetry.io/otel v1.9.0
	go.opentelemetry.io/otel/trace v1.9.0
	go.uber.org/zap v1.21.0
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cactus/go-statsd-client/statsd v0.0.0-00010101000000-000000000000 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gunerhuseyin/goprometheus v0.0.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.4 // indirect
	go.opentelemetry.io/contrib v1.9.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.34.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.9.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.9.0 // indirect
	go.opentelemetry.io/otel/sdk v1.9.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.100.2
	cloud.google.com/go/compute => cloud.google.com/go/compute v1.5.0
	github.com/cactus/go-statsd-client/statsd => github.com/cactus/go-statsd-client/statsd v0.0.0-20200423205355-cb0885a1018c
	gopkg.in/yaml.v2 => gopkg.in/yaml.v3 v3.0.0
)
