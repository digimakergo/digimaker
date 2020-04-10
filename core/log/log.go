//Author xc, Created on 2020-04-09 10:00
//{COPYRIGHTS}
package log

import (
	"context"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

//system info
func Info(message string) {
	log.Info(message)
}

func Warning(message string, label string, ctx ...context.Context) {
	if len(ctx) == 1 {
		logger := GetLogger(ctx[0])
		//write warning to context log and global log
		logger.Warning(message, label)
		log.WithFields(logger.Data).Warning(message, label)
	} else {
		log.Warning(message)
	}
}

//Write error
func Error(message string, label string, ctx ...context.Context) {
	if len(ctx) == 1 {
		logger := GetLogger(ctx[0])
		//white both to context log and global log
		logger.Error(message)
		log.WithFields(logger.Data).Error(message, label)
	} else {
		log.Error(message, label)
	}
}

func Debug(message interface{}, category string, ctx ...context.Context) {
	if len(ctx) == 1 {
		logger := GetLogger(ctx[0])
		logger.Debug(message, "["+category+"]")
	} else {
		log.Debug(message, "["+category+"]")
	}
}

type loggerKey struct{}
type debugInfoKey struct{}
type timerKey struct{}

//Init before request.
func WithLogger(ctx context.Context, fields log.Fields) context.Context {
	//create a new context logger
	logger := log.New()
	logger.SetOutput(ioutil.Discard)
	logger.AddHook(&ContextHook{})
	logger.SetLevel(log.DebugLevel)
	entry := logger.WithFields(fields)
	result := context.WithValue(ctx, loggerKey{}, entry)
	result = context.WithValue(result, timerKey{}, TimerCategory{})
	return result
}

func GetLogger(ctx context.Context) *log.Entry {
	value := ctx.Value(loggerKey{})
	logger := value.(*log.Entry)
	return logger
}

func GetTimer(ctx context.Context) TimerCategory {
	timers := ctx.Value(timerKey{}).(TimerCategory)
	return timers
}

func LogTiming(ctx context.Context) {
	timer := GetTimer(ctx)
	for category, _ := range timer {
		duration, _ := GetDuration(ctx, category)
		Debug(strconv.Itoa(duration)+"ms", category, ctx)
	}
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.0000"})
}
