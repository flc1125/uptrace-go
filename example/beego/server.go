package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/astaxie/beego/otelbeego"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/label"
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("beego-tracer")
)

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	// To enable tracing on template rendering, disable autorender and
	// call otelbeego.Render manually.
	beego.BConfig.WebConfig.AutoRender = false

	beego.Router("/profiles/:username", &ProfileController{})

	mware := otelbeego.NewOTelBeegoMiddleWare("service-name")
	beego.RunWithMiddleWares(":9999", mware)
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}

type ProfileController struct {
	beego.Controller
}

func (c *ProfileController) Get() {
	ctx := c.Ctx.Request.Context()

	username := c.Ctx.Input.Param(":username")
	name, err := selectUser(ctx, username)
	if err != nil {
		c.Abort("404")
		return
	}

	c.Data["username"] = username
	c.Data["name"] = name
	c.TplName = "hello.tpl"

	// Don't forget to call render manually.
	if err := otelbeego.Render(&c.Controller); err != nil {
		c.Abort("500")
	}
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(label.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}