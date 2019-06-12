//Author xc, Created on 2019-05-04 14:26
//{COPYRIGHTS}
//This is a debug based on context.
package debug

import (
	"context"
	"dm/dm/util"
	"errors"
	"time"
)

//This is debug with context

//Key type to main unique
type DebugKey struct{}

//key for context
var debugKey DebugKey = DebugKey{}

type DebugMessage struct {
	Category string
	Type     string
	Message  string
}

type DebugTimer struct {
	Tag        string
	StartPoint int64
	EndPoint   int64
	Duration   int
}

//Debugger Stores all debug information
type Debugger struct {
	List   []DebugMessage
	Timers map[string]*DebugTimer
}

//Add debug into debug struct
func (d *Debugger) Add(debugType string, message string, category string) {
	d.List = append(d.List, DebugMessage{category, debugType, message})
	if debugType == "error" {
		util.Error("[" + category + "]" + message)
	}
}

//Get Debugger instance from context
func GetDebugger(ctx context.Context) *Debugger {
	return ctx.Value(debugKey).(*Debugger)
}

//Add debug into context
func Debug(ctx context.Context, message string, category string) {
	debugger := GetDebugger(ctx)
	debugger.Add("debug", message, category)
}

func Error(ctx context.Context, message string, category string) {
	debugger := GetDebugger(ctx)
	debugger.Add("error", message, category)
}

//start timing
//eg. integration with crm where io: 10ms, rest-request: 200ms, logic handle: 20ms )
//so crm-io, crm-rest,crm-logic are the categories and crm-module is the tag
//They can be the same, or tag can be empty("").
func StartTiming(ctx context.Context, category string, tag string) {
	debugger := GetDebugger(ctx)
	now := time.Now().UnixNano()
	debugger.Timers[category] = &DebugTimer{Tag: tag, StartPoint: now}
}

func EndTiming(ctx context.Context, category string, tag string) {
	debugger := GetDebugger(ctx)
	now := time.Now().UnixNano()
	timer := debugger.Timers[category]
	timer.EndPoint = now
	timer.Duration = int((now - timer.StartPoint) / 1000000)
}

func GetDuration(ctx context.Context, category string) (int, error) {
	debugger := GetDebugger(ctx)
	if result, ok := debugger.Timers[category]; ok {
		return result.Duration, nil
	} else {
		return 0, errors.New("category " + category + "doesn't exist.")
	}
}

//Get duration based on tag.
//Return total and details
func GetDurationByTag(ctx context.Context, tagName string) (int, map[string]int) {
	debugger := GetDebugger(ctx)
	total := 0
	detail := map[string]int{}
	for category, timer := range debugger.Timers {
		if timer.Tag == tagName {
			detail[category] = timer.Duration
			total += timer.Duration
		}
	}
	return total, detail
}

//Initialize Debug from context.
//This should be invoked before context starts(eg. before request/run command line ).
func Init(ctx context.Context) context.Context {
	return context.WithValue(ctx, debugKey, &Debugger{Timers: map[string]*DebugTimer{}})
}
