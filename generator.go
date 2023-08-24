package totp

// Generator is a type that simplifies generation of TOTP codes by allowing configuration to be saved
// to the generator object.
type Generator struct {
	t0            int64
	x             int64
	d             uint32
	currentTimeFn func() int64
}

type Option func(*Generator)

// WithT0 Option sets t0 to the specified value. Default value: 0
// t0 is the Unix time to start counting time steps
func WithT0(t0 int64) Option {
	return func(g *Generator) {
		g.t0 = t0
	}
}

// WithX Option sets x to the specified value. Value must be greater than 0. Default value: 30
// x represents the time step size in seconds
func WithX(x int64) Option {
	if x < 0 {
		x = 30
	}
	return func(g *Generator) {
		g.x = x
	}
}

// WithD Option sets d to the specified value. Default value: 6
// d represents the number of digits the OTP should be
func WithD(d uint32) Option {
	return func(g *Generator) {
		g.d = d
	}
}

// WithCurrentTimeFn sets the currentTimeFn to the specified function. Defaults to a function that returns time.Now().UTC().Unix()
func WithCurrentTimeFn(fn func() int64) Option {
	return func(g *Generator) {
		g.currentTimeFn = fn
	}
}

func NewGenerator(opts ...Option) *Generator {
	const (
		defaultT0 = 0
		defaultX  = 30
		defaultD  = 6
	)

	g := &Generator{
		t0:            defaultT0,
		x:             defaultX,
		d:             defaultD,
		currentTimeFn: func() int64 { return clk.Now().UTC().Unix() },
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func (g Generator) Generate(secret string) (uint32, uint64, error) {
	return Generate(secret, g.t0, g.x, g.d, g.currentTimeFn)
}
