package main

import (
	"fmt"
	"time"

	"github.com/jackramey/totp"
)

func main() {
	secret := "I75NTBBDHWDSX67T" // NOTE: this is not a real secret

	for i := 0; i < 100; i++ {
		totpCode, secondsRemaining := totp.Generate(secret, 0, 30, 6, nil)
		fmt.Printf("Current TOTP code: %d  Time remaining: %d\n", totpCode, secondsRemaining)

		time.Sleep(1 * time.Second)
	}
}
