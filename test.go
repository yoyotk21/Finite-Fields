package main

import (
	"crypto/rand"
	"math/big"
)

func Test(f FiniteField) bool {	
	a, _ := rand.Int(rand.Reader, f.prime)
	b, _ := rand.Int(rand.Reader, f.prime)
	c, _ := rand.Int(rand.Reader, f.prime)
	
	// Associativity check
	if f.add(f.add(a, b),c).Cmp(f.add(a, f.add(b,c))) != 0 {
		return false
	}

	if f.mul(f.mul(a, b),c).Cmp(f.mul(a, f.mul(b,c))) != 0 {
		return false
	}

	// abelian check
	if f.add(a, b).Cmp(f.add(b, a)) != 0 {
		return false
	}

	if f.mul(a, b).Cmp(f.mul(b, a)) != 0 {
		return false
	}

	// identity check
	if f.add(a, big.NewInt(0)).Cmp(a) != 0 {
		return false
	}

	if f.mul(a, big.NewInt(1)).Cmp(a) != 0 {
		return false
	}

	// Inverse check
	if f.add(a, f.addInverse(a)).Cmp(big.NewInt(0)) != 0 {
		return false
	}

	if f.mul(a, f.mulInverse(a)).Cmp(big.NewInt(1)) != 0 {
		return false
	}

	// Distributivity check
	if f.mul(a, f.add(b, c)).Cmp(f.add(f.mul(a, b), f.mul(a, c))) != 0 {
		return false
	}

	return true
}
