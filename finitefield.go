package main

import (
	"errors"
	"math/big"
)

// finite field implementation

type Element interface {
	Val() interface{}
	Equals(e Element) bool
	Add(e Element) Element
	Mul(e Element) Element
	AddInverse() Element
	MulInverse() Element
}

type FiniteField interface {
	AddIdentity() Element
	MulIdentity() Element
	// NewElement(interface{}) (Element, error)
	Contains(e Element) bool
}

type GFPElement struct {
	value *big.Int
	gfp   GFP
}

// helper only
func (g *GFPElement) new(i *big.Int) *GFPElement {
	n := new(GFPElement)
	n.gfp = g.gfp
	n.value = i
	return n
}

func (g *GFPElement) Val() interface{} {
	return g.value
}

func (g *GFPElement) Equals(e Element) bool {
	p := e.(*GFPElement)
	return g.value.Cmp(p.value) == 0
}

func (g *GFPElement) Add(e Element) Element {
	p := e.(*GFPElement) 
	ans := new(big.Int)
	ans = ans.Mod(ans.Add(g.value, p.value), g.gfp.prime)
	return g.new(ans)
}

func (g *GFPElement) Mul(e Element) Element {
	p := e.(*GFPElement)
	ans := new(big.Int)
	ans = ans.Mod(ans.Mul(g.value, p.value), g.gfp.prime)
	return g.new(ans)
}

func (g *GFPElement) AddInverse() Element {
	ans := new(big.Int)
	return g.new(ans.Sub(g.gfp.prime, g.value))
}

func (g *GFPElement) MulInverse() Element {
	ans := new(big.Int)
	return g.new(ans.Exp(g.value, ans.Sub(g.gfp.prime, big.NewInt(2)), g.gfp.prime))
}

type GFP struct {
	prime *big.Int
}

func NewGFP(prime *big.Int) *GFP {
	f := new(GFP)
	f.prime = new(big.Int)
	f.prime.Set(prime)
	return f
}

func NewGFPButEasier(prime int64) *GFP {
	return NewGFP(big.NewInt(prime))
}

// Verifies that numbers exist in the finite field
// func (f GFP) verify(nums ...GFPElement) bool {
// 	for i := range nums {
// 		if nums[i].value.Cmp(f.prime) == 1 {
// 			return false
// 		}
// 	}
// 	return true
// }

// Adds two numbers in a finite field
// func (f GFP) add(a, b GFPElement) GFPElement {
// 	ans := new(big.Int)
// 	ans = ans.Mod(ans.Add(a.value, b.value), f.prime)
// 	return NewGFPElement(ans)
// }

// Multiplies two numbers in a finite field.
// func (f GFP) mul(a, b GFPElement) GFPElement {
// 	ans := new(big.Int)
// 	ans = ans.Mod(ans.Mul(a.value, b.value), f.prime)
// 	return NewGFPElement(ans)
// }

// Returns the additive inverse for a number
// func (f GFP) addInverse(n GFPElement) GFPElement {
// 	ans := new(big.Int)
// 	return NewGFPElement(ans.Sub(f.prime, n.value))
// }

// func (f GFP) mulInverse(n GFPElement) GFPElement {
// 	if n.value.Sign() == 0 {
// 		// Need to figure out a better system for errors
// 		return NewGFPElement(big.NewInt(0))
// 	}
// 	n^(p-2) mod p
// 	ans := new(big.Int)
// 	return NewGFPElement(ans.Exp(n.value, ans.Sub(f.prime, big.NewInt(2)), f.prime))
// }

func (f GFP) NewElement(i *big.Int) (Element, error) {
	if i.Cmp(f.prime) != -1 || i.Cmp(big.NewInt(0)) == -1 {
		return nil, errors.New(i.String() + " does not belong in the Finite Field. Only numbers between 0 and " + f.prime.String() + " belong (not inclusive).")
	}
	n := new(GFPElement)
	n.value = i
	n.gfp = f
	return n, nil
}

func (f GFP) NewElementButEasier(i int) Element {
	ans, _ := f.NewElement(big.NewInt(int64(i)))
	return ans
}

func (f GFP) AddIdentity() Element {
	return f.NewElementButEasier(0)
}

func (f GFP) MulIdentity() Element {
	return f.NewElementButEasier(1)
}

func (f GFP) Contains(e Element) bool {
	return !(e.Val().(*big.Int).Cmp(f.prime) != -1 || e.Val().(*big.Int).Cmp(big.NewInt(0)) == -1)
}
