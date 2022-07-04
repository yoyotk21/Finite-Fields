package main

import "math/big"

// finite field implementation

type FiniteField interface {
	add(a, b *big.Int) *big.Int
	mul(a, b *big.Int) *big.Int
	addInverse(n *big.Int) *big.Int
	mulInverse(n *big.Int) *big.Int
}

type GFP struct {
	prime *big.Int
}

func newGFP(prime *big.Int) *GFP {
	f := new(GFP)
	f.prime = new(big.Int)
	f.prime.Set(prime)
	return f
}

// Verifies that numbers exist in the finite field
func (f GFP) verify(nums ...*big.Int) bool {
	for i := range nums {
		if nums[i].Cmp(f.prime) == 1 {
			return false
		}
	}
	return true
}

// Adds two numbers in a finite field. Returns -1 if they cannot be added.
func (f GFP) add(a, b *big.Int) *big.Int {
	if !f.verify(a, b) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	ans = ans.Mod(ans.Add(a, b), f.prime)
	return ans
}

// Multiplies two numbers in a finite field. Returns -1 if they cannot be multiplied.
func (f GFP) mul(a, b *big.Int) *big.Int {
	if !f.verify(a, b) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	ans = ans.Mod(ans.Mul(a, b), f.prime)
	return ans
}

// Returns the additive inverse for a number
func (f GFP) addInverse(n *big.Int) *big.Int {
	if !f.verify(n) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	return ans.Sub(f.prime, n)

}

func (f GFP) mulInverse(n *big.Int) *big.Int {
	if !f.verify(n) {
		return big.NewInt(-1)
	}

	//n^(p-2) mod p
	ans := new(big.Int)
	return ans.Exp(n, ans.Sub(f.prime, big.NewInt(2)), f.prime)
}
