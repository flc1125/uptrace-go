module github.com/uptrace/uptrace-go/extra/otelgorm

go 1.17

replace github.com/uptrace/uptrace-go/extra/otelsql => ../otelsql

require (
	github.com/uptrace/uptrace-go/extra/otelsql v1.0.4
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	gorm.io/gorm v1.21.16
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	go.opentelemetry.io/otel/internal/metric v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v0.24.0 // indirect
)
