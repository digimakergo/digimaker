//Author xc, Created on 2020-04-09 10:00
//{COPYRIGHTS}
package log

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type ContextInfo struct {
	DebugID   string
	IP        string
	RequestID string
	URI       string
	UserID    int
	Timers    map[string]TimePoint
}

type TimePoint struct {
	Start int64
	End   int64
}

//system info
func Info(message interface{}) {
	log.Info(message)
}

func Warning(message interface{}, label string, ctx ...context.Context) {
	if len(ctx) == 1 {
		fields := GetContextFields(ctx[0])
		log.WithFields(fields).Warning(message, label)
	} else {
		log.Warning(message)
	}
}

//Write error
func Error(message interface{}, label string, ctx ...context.Context) {
	caller := getCallerInfo(runtime.Caller(1))
	if len(ctx) == 1 {
		fields := GetContextFields(ctx[0])
		fields["caller"] = caller
		log.WithFields(fields).Error(message, label)
	} else {
		fields := log.Fields{}
		fields["caller"] = caller
		log.WithFields(fields).Error(message, label)
	}
}

func Fatal(message interface{}) {
	log.Fatal(message)
}

//Output debug info with on category.
func Debug(message interface{}, category string, ctx ...context.Context) {
	caller := getCallerInfo(runtime.Caller(1))
	if len(ctx) == 1 {
		fields := GetContextFields(ctx[0])
		fields["caller"] = caller
		fields["category"] = category
		log.WithFields(fields).Debug(message)
	} else {
		fields := getFields(category)
		fields["caller"] = caller
		log.WithFields(fields).Debug(message)
	}
}

func getFields(category string) log.Fields {
	fields := log.Fields{}
	fields["category"] = category
	return fields
}

func GetContextFields(ctx context.Context) log.Fields {
	info := GetContextInfo(ctx)
	fields := log.Fields{}
	fields["ip"] = info.IP
	fields["request_id"] = info.RequestID
	fields["user_id"] = info.UserID
	fields["debug_id"] = info.DebugID
	fields["uri"] = info.URI
	return fields
}

type logKey struct{}

// init a context log
func InitContext(ctx context.Context, info *ContextInfo) context.Context {
	newContext := context.WithValue(ctx, logKey{}, info)
	return newContext
}

func GetContextInfo(ctx context.Context) *ContextInfo {
	return ctx.Value(logKey{}).(*ContextInfo)
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
		log.WithFields(fields).Debug(strconv.Itoa(duration) + "ms")
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
	//todo: log it to file based on parameters
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.0000"})
}
