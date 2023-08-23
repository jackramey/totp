package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

func Generate(secret string, t0, x int64, d uint32, currentTimeFn func() int64) (uint32, uint64) {
	secretBytes, err := decodeSecret(secret)
	if err != nil {
		panic(err)
	}

	tFunc := tFn[uint64](t0, x, currentTimeFn)
	timeBytes := make([]byte, 8)
	t, r := tFunc()
	binary.BigEndian.PutUint64(timeBytes, t)

	return hotp(secretBytes, timeBytes, d), r
}

// tFn takes a time function as an input, and returns a function that provides the T and R values where T is the number
// of time steps between the initial counter time T0 and the current time and R is the number of seconds remaining until
// the next code is produced.
// t0 is the Unix time to start counting time steps (default value is 0, i.e., the Unix epoch) and is also a system parameter
// x represents the time step in seconds (default value X = 30 seconds) and is a system parameter
// If currentTimeFn is nil it will default to time.Now().UTC().Unix
func tFn[T int64 | uint64](t0, x int64, currentTimeFn func() int64) func() (T, T) {
	return func() (T, T) {
		if currentTimeFn == nil {
			currentTimeFn = time.Now().UTC().Unix
		}
		if x == 0 {
			x = 30
		}

		t := T((currentTimeFn() - t0) / x)
		r := T(x - ((currentTimeFn() - t0) % x))
		return t, r
	}
}

// decodeSecret decodes the base32 input secret string and returns the binary value.
// In order to ensure consistency the decoder will trim all whitespace from the secret string and cast all characters
// to their uppercase values before decoding.
func decodeSecret(secret string) ([]byte, error) {
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	secret = strings.TrimSpace(strings.ToUpper(secret))
	secretBytes, err := decoder.DecodeString(secret)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error decoding secret value: %s", err.Error()))
	}

	return secretBytes, nil
}

// hotp calculates an HMAC-Based One-Time Passcode based on [RFC-4226](https://datatracker.ietf.org/doc/html/rfc4226)
func hotp(k, c []byte, d uint32) uint32 {
	h := hash(k, c)
	return dt(h, d)
}

// hash performs a sha1 hash on c with the secret k
func hash(k, c []byte) []byte {
	h := hmac.New(sha1.New, k)
	h.Write(c)
	return h.Sum(nil)
}

// dt performs the dynamic truncation defined in [RFC-4226 Section 5.3](https://datatracker.ietf.org/doc/html/rfc4226#section-5.3)
// and perfoms a modulo on the result value with 10^{d} to return a uint32 value with d number of digits
func dt(b []byte, d uint32) uint32 {
	offset := b[len(b)-1] & 0x0F // Offset bits are the 4 least significant bits of the hash
	// AND with 0x7FFFFFFF clears out the most significant bit to avoid confusion around signed vs unsigned modulo computations
	sNum := binary.BigEndian.Uint32(b[offset:]) & 0x7FFFFFFF // AND with 0x7FFFFFF
	return sNum % uint32(math.Pow(10, float64(d)))           // truncate to the defined number of digits of the value
}
