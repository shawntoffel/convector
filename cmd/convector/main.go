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
	jitter    = time.Minute * 5
	endpoints = "https://myip2.shawntoffel.com, https://myip.shawntoffel.com"
	timeout   = time.Second * 10
)

func init() {
	flag.DurationVar(&interval, "i", interval, "base interval between requests")
	flag.DurationVar(&jitter, "j", jitter, "jitter between requests")
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
		convector.ConvectorOptions{
			Interval: interval,
			Jitter:   jitter,
		},
		&http.Client{Timeout: timeout},
	)

	convector.Start()

	select {
	case sig := <-interrupt:
		log.Println("received signal:", sig)
		convector.Stop()
		return
	}
}
