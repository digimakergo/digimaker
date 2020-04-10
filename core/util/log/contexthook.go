package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type ContextHook struct {
}

func (hook *ContextHook) Fire(entry *log.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}
	//todo: based on settings(eg. debug by ip/user), output context log information.
	f, err := os.OpenFile("request-debug.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	f.Write(line)
	defer f.Close()
	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *ContextHook) Levels() []log.Level {
	return []log.Level{
		log.DebugLevel,
		log.InfoLevel,
		log.ErrorLevel,
		log.InfoLevel,
	}
}
