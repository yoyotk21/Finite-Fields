package main

import (
	"fmt"
	"math/big"
)

func main() {
	f := NewFFM(newGFP(big.NewInt(7)), 2, 2)
	arr := [][]*big.Int{
		{big.NewInt(6), big.NewInt(4)},
		{big.NewInt(1), big.NewInt(1)}}
	f.set(arr)
	// g := NewFFM(newGFP(big.NewInt(30)), 2, 2)
	// arr2 := [][]*big.Int {
	// 	{big.NewInt(4), big.NewInt(3)},
	// 	{big.NewInt(2), big.NewInt(1)}}
	// g.set(arr2)

	fmt.Println(FFMTest(f))
}
