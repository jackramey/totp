package main

import (
	"fmt"
	"time"

	"github.com/jackramey/totp"
)

func main() {
	secret := "I75NTBBDHWDSX67T" // NOTE: this is not a real secret

	g := totp.NewGenerator()

	for i := 0; i < 100; i++ {
		totpCode, secondsRemaining := g.Generate(secret)
		fmt.Printf("Current TOTP code: %d  Time remaining: %d\n", totpCode, secondsRemaining)

		time.Sleep(1 * time.Second)
	}
}
