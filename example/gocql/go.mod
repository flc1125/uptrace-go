module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210817081954-bc256bbb90de
	github.com/golang/snappy v0.0.4 // indirect
	github.com/uptrace/uptrace-go v1.0.2
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.24.0
	go.opentelemetry.io/otel v1.0.0
	go.opentelemetry.io/otel/trace v1.0.0
)
