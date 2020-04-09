package log

import (
	"context"
	"errors"
	"time"
)

type Timer struct {
	Points []int64
}

type TimerCategory map[string]*Timer

//start timing
func StartTiming(ctx context.Context, category string) {
	timer := GetTimer(ctx)
	now := time.Now().UnixNano()
	timer[category] = &Timer{Points: []int64{now}}
}

func TickTiming(ctx context.Context, category string) {
	timers := GetTimer(ctx)
	now := time.Now().UnixNano()
	timer := timers[category]
	timer.Points = append(timer.Points, now)
	timers[category] = timer
}

func GetDuration(ctx context.Context, category string) (int, error) {
	timers := GetTimer(ctx)
	if result, ok := timers[category]; ok {
		points := result.Points
		duration := int((points[len(points)-1] - points[0]) / 1000000)
		return duration, nil
	} else {
		return 0, errors.New("category " + category + "doesn't exist.")
	}
}
