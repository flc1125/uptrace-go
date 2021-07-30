package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/uptrace/uptrace-go/example/grpc/api"
	"github.com/uptrace/uptrace-go/uptrace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	target := os.Getenv("GRPC_TARGET")
	if target == "" {
		target = ":9999"
	}

	log.Println("connecting to", target)

	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() { _ = conn.Close() }()

	client := api.NewHelloServiceClient(conn)
	if err := sayHello(client); err != nil {
		log.Fatal(err)
		return
	}
}

func sayHello(client api.HelloServiceClient) error {
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client",
		"user-id", "test-user",
	))

	resp, err := client.SayHello(ctx, &api.HelloRequest{Greeting: "World"})
	if err != nil {
		return err
	}
	log.Println("reply:", resp.Reply)

	return nil
}
