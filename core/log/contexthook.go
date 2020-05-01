//Author xc, Created on 2020-04-09 19:00
//{COPYRIGHTS}
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

//set debug setting in a local variable
var debugIps []string

//can debug a category or not
func CanDebug(ip string) bool {
	result := false
	for _, value := range debugIps {
		if value == ip {
			result = true
			break
		}
	}
	return result
}

func SetDebugIp() {

}

func init() {
	//todo: set it from another place
	debugIps = []string{"127.0.0.1", "::1"}
}
