package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

const (
	defaultPort = "8080"
)

func main() {

	funcframework.RegisterHTTPFunction("/", hello)
	port := getPort()

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

func getPort() string {

	port := defaultPort

	if envPort := os.Getenv("PORT"); port != "" {
		port = envPort
	}

	return port
}
