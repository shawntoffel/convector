package convector

import (
	"net/http"
	"strings"
	"time"

	"github.com/lthibault/jitterbug/v2"
)

var LogLevel = LogInfo

type Convector struct {
	ticker    *jitterbug.Ticker
	endpoints []string
	client    *http.Client
	quit      chan bool
}

func New(endpoints []string) Convector {
	return NewWithHttpClient(endpoints, &http.Client{})
}

func NewWithHttpClient(endpoints []string, client *http.Client) Convector {
	return Convector{
		endpoints: endpoints,
		client:    client,
		quit:      make(chan bool),
	}
}

func (c Convector) Start(interval time.Duration, stdev time.Duration) {
	if c.ticker != nil {
		c.ticker.Stop()
	}

	c.log(LogInfo, "starting with interval: %s, stdev: %s", interval, stdev)

	c.ticker = jitterbug.New(
		interval,
		&jitterbug.Norm{Stdev: stdev},
	)

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
	if c.ticker != nil {
		c.ticker.Stop()
	}
	c.quit <- true
	c.log(LogInfo, "stopped")
}

func (c Convector) tick() {
	c.log(LogDebug, "starting endpoint execution for interval")

	for _, endpoint := range c.endpoints {
		e := strings.TrimSpace(endpoint)
		if e == "" {
			continue
		}

		c.log(LogDebug, "sending request to endpoint: %s", e)

		start := time.Now()
		resp, err := c.client.Get(e)
		took := time.Since(start)
		if err == nil {
			c.log(LogDebug, "endpoint '%s' responded with status code %d in %s", e, resp.StatusCode, took)
		} else {
			c.log(LogError, err.Error())
		}
	}

	c.log(LogDebug, "finished endpoint execution for interval")
}

func (c Convector) log(level int, format string, a ...interface{}) {
	if level > LogLevel {
		return
	}

	levellog(level, format, a...)
}
