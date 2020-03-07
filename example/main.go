package main

import (
	"fmt"

	"github.com/wiless/channel"
)

func main() {

	env := channel.NewSimpleEnv()

	env.Setup(.700, 10) // fc=700MHz, Bw=10Mhz
	env.Create(3, 1, 1)

	iidgen := channel.NewGeneratorIID()
	iidgen.Reset(1234)
	env.AttachGeneratorIID(*iidgen, true)

	for idx, link := range env.Links {

		for t := 0.0; t < 5; t++ {
			{
				fmt.Printf("\nLink (%d) t=%f : %v [%v]", idx, t, link.NextSample(), link.State())
			}
		}
	}
	fmt.Println()
	for idx, link := range env.Links {

		for t := 0.0; t < 5; t++ {
			{
				fmt.Printf("\nLink (%d) t=%f : %v [%v]", idx, t, link.NextSample(), link.State())
			}
		}
	}

	fmt.Println()
}
