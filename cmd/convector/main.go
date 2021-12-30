package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shawntoffel/convector"
	"github.com/shawntoffel/convector/internal"
)

var (
	flagVersion = false
	flagDebug   = false

	interval  = time.Minute * 10
	stdev     = time.Minute * 5
	timeout   = time.Second * 10
	endpoints = ""
)

func init() {
	flag.DurationVar(&interval, "i", interval, "base interval between requests")
	flag.DurationVar(&stdev, "j", stdev, "standard deviation for a normal distribution jitter between requests")
	flag.DurationVar(&timeout, "t", timeout, "http client timeout for requests")
	flag.StringVar(&endpoints, "e", endpoints, "comma-delimited list of endpoints")
	flag.BoolVar(&flagVersion, "v", false, "print convector version")
	flag.BoolVar(&flagDebug, "d", false, "enable debug logging")

	flag.Parse()
}

func main() {
	if flagVersion {
		fmt.Println(internal.Version)
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	if flagDebug {
		convector.LogLevel = convector.LogDebug
	}

	convector := convector.NewWithHttpClient(
		strings.Split(endpoints, ","),
		&http.Client{Timeout: timeout},
	)

	convector.Start(interval, stdev)

	select {
	case sig := <-interrupt:
		log.Println("received signal:", sig)
		convector.Stop()
		return
	}
}
