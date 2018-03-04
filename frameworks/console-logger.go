package frameworks

import "log"

type ConsoleLogger struct{}

func (d ConsoleLogger) Log(level, message string) {
	log.Printf("[%s] %s\n", level, message)
}
