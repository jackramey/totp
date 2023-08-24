package totp

import "time"

var clk clock = stdClock{}

type clock interface {
	Now() time.Time
}

type stdClock struct{}

func (s stdClock) Now() time.Time {
	return time.Now()
}
