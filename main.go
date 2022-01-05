package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cgauge/aws-secretsmanager-caching-extension/extensions_api"
	"github.com/cgauge/aws-secretsmanager-caching-extension/secrets"
	"github.com/cgauge/aws-secretsmanager-caching-extension/server"
)

var (
	extensionClient = extensions_api.NewClient(os.Getenv("AWS_LAMBDA_RUNTIME_API"))
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		s := <-sigs
		cancel()
		println(secrets.PrintPrefix, "Received", s)
		println(secrets.PrintPrefix, "Exiting")
	}()

	_, err := extensionClient.Register(ctx, secrets.ExtensionName)
	if err != nil {
		panic(err)
	}

	go server.StartHTTPServer("8015")

	processEvents(ctx)
}

func processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			println(secrets.PrintPrefix, "Waiting for event...")
			res, err := extensionClient.NextEvent(ctx)
			if err != nil {
				println(secrets.PrintPrefix, "Error:", err)
				println(secrets.PrintPrefix, "Exiting")
				return
			}

			if res.EventType == extensions_api.Shutdown {
				println(secrets.PrintPrefix, "Received SHUTDOWN event")
				println(secrets.PrintPrefix, "Exiting")
				return
			}
		}
	}
}
