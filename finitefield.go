package main

import "math/big"

// finite field implementation

// Sort of useless at the moment, not really sure what to do with this
type FiniteFielder interface {
	add(a, b big.Int) big.Int
	mul(a, b big.Int) big.Int
	addInverse(n big.Int) big.Int
	mulInverse(n big.Int) big.Int
}

type FiniteField struct {
	prime *big.Int
}

func newFiniteField(prime *big.Int) FiniteField {
	f := new(FiniteField)
	f.prime = new(big.Int)
	f.prime.Set(prime)
	return *f
}

// Verifies that numbers exist in the finite field
func (f FiniteField) verify(nums ...*big.Int) bool {
	for i := range nums {
		if nums[i].Cmp(f.prime) == 1 {
			return false
		}
	}
	return true
}

// Adds two numbers in a finite field. Returns -1 if they cannot be added.
func (f FiniteField) add(a, b *big.Int) *big.Int {
	if !f.verify(a, b) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	ans = ans.Mod(ans.Add(a, b), f.prime)
	return ans
}

// Multiplies two numbers in a finite field. Returns -1 if they cannot be multiplied.
func (f FiniteField) mul(a, b *big.Int) *big.Int {
	if !f.verify(a, b) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	ans = ans.Mod(ans.Mul(a, b), f.prime)
	return ans
}

// Returns the additive inverse for a number
func (f FiniteField) addInverse(n *big.Int) *big.Int {
	if !f.verify(n) {
		return big.NewInt(-1)
	}
	ans := new(big.Int)
	return ans.Sub(f.prime, n)

}

func (f FiniteField) mulInverse(n *big.Int) *big.Int {
	if !f.verify(n) {
		return big.NewInt(-1)
	}

	//n^(p-2) mod p
	ans := new(big.Int)
	return ans.Exp(n, ans.Sub(f.prime, big.NewInt(2)), f.prime)
}
