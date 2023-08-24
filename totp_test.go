package totp

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name              string
		secret            string
		t0                int64
		x                 int64
		d                 uint32
		currentTimeFn     func() int64
		wantCode          uint32
		wantTimeRemaining uint64
		wantErr           bool
	}{
		{
			name:              "Generate a code with 6 digits",
			secret:            "ITSASECRETSHHHHH",
			t0:                0,
			x:                 30,
			d:                 6,
			currentTimeFn:     func() int64 { return 59 },
			wantCode:          467194,
			wantTimeRemaining: 1,
		},
		{
			name:              "Generate a code with 4 digits",
			secret:            "ITSASECRETSHHHHH",
			t0:                0,
			x:                 30,
			d:                 4,
			currentTimeFn:     func() int64 { return 59 },
			wantCode:          7194,
			wantTimeRemaining: 1,
		},
		{
			name:              "Lowercase secrets shouldn't affect the code",
			secret:            "itsasecretshhhhh",
			t0:                0,
			x:                 30,
			d:                 6,
			currentTimeFn:     func() int64 { return 59 },
			wantCode:          467194,
			wantTimeRemaining: 1,
		},
		{
			name:              "Next step should produce a different code and more time remaining",
			secret:            "ITSASECRETSHHHHH",
			t0:                0,
			x:                 30,
			d:                 6,
			currentTimeFn:     func() int64 { return 60 },
			wantCode:          858003,
			wantTimeRemaining: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotTimeRemaining, err := Generate(tt.secret, tt.t0, tt.x, tt.d, tt.currentTimeFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("Generate() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
			if gotTimeRemaining != tt.wantTimeRemaining {
				t.Errorf("Generate() gotTimeRemaining = %v, want %v", gotTimeRemaining, tt.wantTimeRemaining)
			}
		})
	}
}

func TestTFn_int64(t *testing.T) {
	type testCase[T interface{ int64 | uint64 }] struct {
		name          string
		t0            int64
		x             int64
		currentTimeFn func() int64
		wantT         int64
		wantR         int64
	}
	tests := []testCase[int64]{
		{
			t0:            0,
			x:             30,
			currentTimeFn: func() int64 { return 59 },
			wantT:         1,
			wantR:         1,
		},
		{
			t0:            0,
			x:             30,
			currentTimeFn: func() int64 { return 60 },
			wantT:         2,
			wantR:         30,
		},
	}
	for _, tt := range tests {
		t.Run("int64 tests", func(t *testing.T) {
			tFunc := tFn[int64](tt.t0, tt.x, tt.currentTimeFn)
			gotT, gotR := tFunc()
			if gotT != tt.wantT {
				t.Errorf("t = %d, want %d", gotT, tt.wantT)
			}
			if gotR != tt.wantR {
				t.Errorf("r = %d, want %d", gotR, tt.wantR)
			}
		})
	}
}

func TestTFn_uint64(t *testing.T) {
	type testCase[T interface{ int64 | uint64 }] struct {
		name          string
		t0            int64
		x             int64
		currentTimeFn func() int64
		wantT         uint64
		wantR         uint64
	}
	tests := []testCase[uint64]{
		{
			t0:            0,
			x:             30,
			currentTimeFn: func() int64 { return 59 },
			wantT:         1,
			wantR:         1,
		},
		{
			t0:            0,
			x:             30,
			currentTimeFn: func() int64 { return 60 },
			wantT:         2,
			wantR:         30,
		},
	}
	for _, tt := range tests {
		t.Run("uint64 tests", func(t *testing.T) {
			tFunc := tFn[uint64](tt.t0, tt.x, tt.currentTimeFn)
			gotT, gotR := tFunc()
			if gotT != tt.wantT {
				t.Errorf("t = %d, want %d", gotT, tt.wantT)
			}
			if gotR != tt.wantR {
				t.Errorf("r = %d, want %d", gotR, tt.wantR)
			}
		})
	}
}
