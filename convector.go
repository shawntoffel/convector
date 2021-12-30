package convector

import (
	"net/http"
	"strings"
	"time"

	"github.com/lthibault/jitterbug/v2"
)

var LogLevel = LogInfo

type ConvectorOptions struct {
	Interval time.Duration
	Jitter   time.Duration
}

type Convector struct {
	ticker    *jitterbug.Ticker
	endpoints []string
	client    *http.Client
	opts      ConvectorOptions
	quit      chan bool
}

func New(endpoints []string, opts ConvectorOptions) Convector {
	return NewWithHttpClient(endpoints, opts, &http.Client{})
}

func NewWithHttpClient(endpoints []string, opts ConvectorOptions, client *http.Client) Convector {
	ticker := jitterbug.New(
		opts.Interval,
		&jitterbug.Norm{Stdev: opts.Jitter},
	)

	return Convector{
		ticker:    ticker,
		endpoints: endpoints,
		client:    client,
		quit:      make(chan bool),
		opts:      opts,
	}
}

func (c Convector) Start() {
	c.log(LogInfo, "starting with interval: %s, jitter: %s", c.opts.Interval, c.opts.Jitter)
	go func() {
		for {
			select {
			case _ = <-c.ticker.C:
				c.tick()
			case <-c.quit:
				return
			}
		}
	}()
}

func (c Convector) Stop() {
	c.log(LogInfo, "stopping")
	c.ticker.Stop()
	c.quit <- true
	c.log(LogInfo, "stopped")
}

func (c Convector) tick() {
	c.log(LogDebug, "starting endpoint execution for interval")

	for _, endpoint := range c.endpoints {
		e := strings.TrimSpace(endpoint)

		c.log(LogDebug, "executing endpoint: %s", e)

		_, err := c.client.Get(e)
		if err != nil {
			c.log(LogError, err.Error())
		} else {
			c.log(LogDebug, "finished executing endpoint: %s", e)
		}
	}

	c.log(LogDebug, "finished endpoint execution for inteval")
}

func (c Convector) log(level int, format string, a ...interface{}) {
	if level > LogLevel {
		return
	}

	levellog(level, format, a...)
}
