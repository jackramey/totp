package main

import (
	"fmt"
	"time"

	"github.com/jackramey/totp/totp"
)

func main() {
	secret := "FIQAUAK4DIDGHEN5FTT7Q7BW3GMUDWMM" // NOTE: this is not a real secret

	g := totp.NewGenerator()

	for i := 0; i < 100; i++ {
		totpCode, secondsRemaining, err := g.Generate(secret)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Current TOTP code: %d  Time remaining: %d\n", totpCode, secondsRemaining)

		time.Sleep(1 * time.Second)
	}
}
