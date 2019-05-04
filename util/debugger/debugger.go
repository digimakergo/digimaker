//Author xc, Created on 2019-05-04 14:26
//{COPYRIGHTS}
package debugger

import (
	"context"
	"time"
)

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
	Identifier string
	StartPoint int64
	EndPoint   int64
	Duration   int
}

type Debugger struct {
	List   []DebugMessage
	Timers map[string]*DebugTimer
}

//Add debug into debug struct
func (d *Debugger) Add(debugType string, message string, category string) {
	d.List = append(d.List, DebugMessage{category, debugType, message})
	if debugType == "error" || debugType == "warning" {
		//todo: log into log file
	}
}

//Get Debugger instance from context
func GetDebugger(ctx context.Context) *Debugger {
	return ctx.Value(debugKey).(*Debugger)
}

//Add debug into context
func AddDebug(ctx context.Context, message string, category string) {
	debugger := GetDebugger(ctx)
	debugger.Add("debug", message, category)
}

func AddWarning(ctx context.Context, message string, category string) {
	debugger := GetDebugger(ctx)
	debugger.Add("warning", message, category)
}

func AddError(ctx context.Context, message string, category string) {
	debugger := GetDebugger(ctx)
	debugger.Add("error", message, category)
}

func StartTiming(ctx context.Context, category string, identifier string) {
	debugger := GetDebugger(ctx)
	now := time.Now().UnixNano()
	debugger.Timers[category] = &DebugTimer{Identifier: identifier, StartPoint: now}
}

func EndTiming(ctx context.Context, category string, identifier string) {
	debugger := GetDebugger(ctx)
	now := time.Now().UnixNano()
	timer := debugger.Timers[category]
	timer.EndPoint = now
	timer.Duration = int((now - timer.StartPoint) / 1000000)
}

//Initialize Debug from context.
//This should be invoked before context starts(eg. before request/run command line ).
func Init(ctx context.Context) context.Context {
	return context.WithValue(ctx, debugKey, &Debugger{Timers: map[string]*DebugTimer{}})
}
