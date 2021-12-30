package convector

import (
	"fmt"
	"log"
)

const (
	LogError int = iota
	LogInfo
	LogDebug
)

var Logger func(level int, format string, a ...interface{})

func levellog(level int, format string, a ...interface{}) {
	if Logger != nil {
		Logger(level, format, a...)
	} else {
		log.Printf("convector [%s]: %s\n", fmtLevel(level), fmt.Sprintf(format, a...))
	}
}

func fmtLevel(level int) string {
	switch level {
	case LogError:
		return "ERROR"
	case LogInfo:
		return "INFO"
	case LogDebug:
		return "DEBUG"
	default:
		return ""
	}
}
