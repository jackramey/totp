package main

import (
	"fmt"
	"time"

	"github.com/jackramey/totp"
)

func main() {
	secret := "I75NTBBDHWDSX67T"

	g := totp.NewGenerator(totp.WithD(8))

	for i := 0; i < 100; i++ {
		totpCode, secondsRemaining := totp.Generate(secret, 0, 30, 6, nil)
		fmt.Printf("Current TOTP code: %d  Time remaining: %d\n", totpCode, secondsRemaining)

		totpCode, secondsRemaining = g.Generate(secret)
		fmt.Printf("Current TOTP code: %d  Time remaining: %d\n", totpCode, secondsRemaining)

		time.Sleep(1 * time.Second)
	}
}
