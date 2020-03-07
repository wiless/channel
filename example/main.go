package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/wiless/channel"
)

func main() {
	// Reset the seed
	rand.Seed(time.Now().UnixNano())

	env := channel.NewSimpleEnv()
	env.Setup(0.700, 10) // fc=700MHz, Bw=10Mhz
	env.Create(3, 1, 1)

	// Example to set all links to an i.i.d generator
	env.AttachGeneratorIID()
	for idx, link := range env.Links {
		for t := 0.0; t < 5; t++ {
			{
				fmt.Printf("\nLink (%d) t=%f : %v", idx, t, link.NextSample())
			}
		}
	}

	fmt.Println()
}
