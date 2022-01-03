//Author xc, Created on 2020-04-09 10:00
//{COPYRIGHTS}
package log

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

type ContextInfo struct {
	IP        string
	RequestID string
	URI       string
	UserID    int
	Timers    map[string]TimePoint
	Debug     bool
}

type TimePoint struct {
	Start int64
	End   int64
}

//system info
func Info(message interface{}) {
	logrus.Info(message)
}

func Warning(message interface{}, label string, ctx ...context.Context) {
	if len(ctx) == 1 {
		fields := GetContextFields(ctx[0])
		logrus.WithFields(fields).Warning(message, label)
	} else {
		logrus.Warning(message)
	}
}

//Write error
func Error(message interface{}, label string, ctx ...context.Context) {
	caller := getCallerInfo(runtime.Caller(1))
	if len(ctx) == 1 {
		fields := GetContextFields(ctx[0])
		fields["caller"] = caller
		fields["category"] = label
		logrus.WithFields(fields).Error(message, label)
	} else {
		fields := logrus.Fields{}
		fields["caller"] = caller
		fields["category"] = label
		logrus.WithFields(fields).Error(message, label)
	}
}

func Fatal(message interface{}) {
	logrus.Fatal(message)
}

//Output debug info with on category.
func Debug(message interface{}, category string, ctx ...context.Context) {
	caller := getCallerInfo(runtime.Caller(1))
	if len(ctx) == 1 {
		info := GetContextInfo(ctx[0])
		if info != nil && info.Debug {
			fields := GetContextFields(ctx[0])
			fields["caller"] = caller
			fields["category"] = category
			fields["type"] = "message"
			logrus.WithFields(fields).Debug(message)
		}
	} else {
		fields := getFields(category)
		fields["caller"] = caller
		fields["type"] = "message"
		logrus.WithFields(fields).Debug(message)
	}
}

func getFields(category string) logrus.Fields {
	fields := logrus.Fields{}
	fields["category"] = category
	return fields
}

func GetContextFields(ctx context.Context) logrus.Fields {
	info := GetContextInfo(ctx)
	fields := logrus.Fields{}
	if info != nil {
		fields["ip"] = info.IP
		fields["request_id"] = info.RequestID
		fields["user_id"] = info.UserID
		fields["uri"] = info.URI
	}
	return fields
}

type logKey struct{}

// init a context log
func InitContext(ctx context.Context, info *ContextInfo) context.Context {
	newContext := context.WithValue(ctx, logKey{}, info)
	return newContext
}

func GetContextInfo(ctx context.Context) *ContextInfo {
	info := ctx.Value(logKey{})
	if info != nil {
		return info.(*ContextInfo)
	} else {
		return nil
	}
}

//start timing
func StartTiming(ctx context.Context, category string) {
	info := GetContextInfo(ctx)
	now := time.Now().UnixNano()
	timer := TimePoint{Start: now}
	if info.Timers == nil {
		info.Timers = map[string]TimePoint{}
	}

	info.Timers[category] = timer
}

//End timing on a category
func EndTiming(ctx context.Context, category string) {
	info := GetContextInfo(ctx)

	now := time.Now().UnixNano()
	timer := info.Timers[category]
	timer.End = now

	info.Timers[category] = timer
}

//Log all timing, usually done in the end of request
func LogTiming(ctx context.Context) {
	info := GetContextInfo(ctx)
	for category, timer := range info.Timers {
		duration := int((timer.End - timer.Start) / 1000000)
		fields := GetContextFields(ctx)
		fields["type"] = "timing"
		fields["category"] = category
		logrus.WithFields(fields).Debug(strconv.Itoa(duration) + "ms")
	}
}

func getCallerInfo(pc uintptr, file string, line int, ok bool) string {
	if ok {
		name := runtime.FuncForPC(pc).Name()
		result := name + ": " + strconv.Itoa(line)
		return result
	} else {
		return ""
	}
}

func init() {
	environment := "prod"
	if envValue := os.Getenv("env"); envValue != "" {
		environment = envValue
	}

	logrus.SetLevel(logrus.DebugLevel)

	dmapp := os.Getenv("dmapp")
	if dmapp == "" {
		dmapp = "."
	}

	if environment == "prod" {
		logrus.SetOutput(ioutil.Discard)
		//log server output
		logrus.AddHook(&RemoteHook{})

		//default output
		//todo: create /var/log/digimaker/digimaker.log automatically. need to think about permission
		file, err := os.OpenFile(dmapp+"/digimaker.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			log.Fatal("Log folder doesn't exist")
		}

		logrus.AddHook(&writer.Hook{
			Writer: file,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.WarnLevel,
				logrus.ErrorLevel,
				logrus.FatalLevel,
			},
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.0000"})
	}

}
