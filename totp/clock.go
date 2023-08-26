package totp

import "time"

var Clock clock = stdClock{}

type clock interface {
	Now() time.Time
}

type stdClock struct{}

func (s stdClock) Now() time.Time {
	return time.Now()
}
