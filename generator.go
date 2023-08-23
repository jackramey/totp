package totp

import "time"

type Generator struct {
	T0            int64
	X             int64
	D             uint32
	CurrentTimeFn func() int64
}

type Option func(*Generator)

func WithT0(t0 int64) Option {
	return func(g *Generator) {
		g.T0 = t0
	}
}

func WithX(x int64) Option {
	return func(g *Generator) {
		g.X = x
	}
}

func WithD(d uint32) Option {
	return func(g *Generator) {
		g.D = d
	}
}

func WithCurrentTimeFn(fn func() int64) Option {
	return func(g *Generator) {
		g.CurrentTimeFn = fn
	}
}

func NewGenerator(opts ...Option) *Generator {
	const (
		defaultT0 = 0
		defaultX  = 30
		defaultD  = 6
	)

	g := &Generator{
		T0:            defaultT0,
		X:             defaultX,
		D:             defaultD,
		CurrentTimeFn: func() int64 { return time.Now().UTC().Unix() },
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g Generator) Generate(secret string) (uint32, uint64) {
	return Generate(secret, g.T0, g.X, g.D, g.CurrentTimeFn)
}
