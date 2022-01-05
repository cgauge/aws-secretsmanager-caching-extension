package main

import (
	"github.com/cgauge/aws-secretsmanager-caching-extension/server"
)

func main() {
	server.StartHTTPServer("8015")
}
