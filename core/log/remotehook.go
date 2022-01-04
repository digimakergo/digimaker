//Author xc, Created on 2020-04-09 19:00
//{COPYRIGHTS}
package log

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type RemoteHook struct {
}

func (hook *RemoteHook) Fire(entry *log.Entry) error {
	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	dmapp := os.Getenv("dmapp")
	if dmapp == "" {
		dmapp = "."
	}

	//todo: based on settings(eg. debug by ip/user), output context log information.
	logFolder := dmapp + "/var/debug"
	logFile := logFolder + "/request-debug.log"
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		os.Mkdir(logFolder, 0775)
	}

	f, err := os.OpenFile(logFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	f.Write(line)
	defer f.Close()
	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *RemoteHook) Levels() []log.Level {
	return log.AllLevels
}
