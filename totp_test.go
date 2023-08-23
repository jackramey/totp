package totp

import (
	"testing"
)

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
