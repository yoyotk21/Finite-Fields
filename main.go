package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	r := big.NewInt(3)
	fmt.Println(r)
	b := big.NewInt(17)
	f := newFiniteField(b)
	b.Set(big.NewInt(301))
	fmt.Println(f.add(big.NewInt(12), big.NewInt(12)))

	for i := 0; i < 1000; i++ {
		p, _ := rand.Prime(rand.Reader, 128)
		if !Test(newFiniteField(p)) {
			fmt.Println("Error!")
		}
	}
}
