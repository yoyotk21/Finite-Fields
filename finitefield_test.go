package main

import (
	"crypto/rand"
	"testing"
)

const NumOfTests = 100

func TestFiniteField(t *testing.T) {

	for i := 0; i < NumOfTests; i++ {

		prime, _ := rand.Prime(rand.Reader, 128)
		finiteField := NewGFP(prime)

		r1, _ := rand.Int(rand.Reader, prime)
		r2, _ := rand.Int(rand.Reader, prime)
		r3, _ := rand.Int(rand.Reader, prime)

		a, _ := finiteField.NewElement(r1)
		b, _ := finiteField.NewElement(r2)
		c, _ := finiteField.NewElement(r3)

		// Associativity test

		if !a.Add(b.Add(c)).Equals(a.Add(b).Add(c)) || !a.Mul(b.Mul(c)).Equals(a.Mul(b).Mul(c)) {
			t.Errorf("Finite field not associative")
		}

		// Commutativity test

		if !a.Add(b).Equals(b.Add(a)) || !a.Mul(b).Equals(b.Mul(a)) {
			t.Errorf("Finite field not commutative")
		}

		// Identities test

		if !finiteField.AddIdentity().Add(a).Equals(a) {
			t.Errorf("Additive identity not working")
		}

		if !finiteField.MulIdentity().Mul(b).Equals(b) {
			t.Errorf("Multiplicative identity not working")
		}

		// Inverse test

		if !c.Add(c.AddInverse()).Equals(finiteField.AddIdentity()) {
			t.Errorf("Additive inverse not working")
		}

		if !a.Mul(a.MulInverse()).Equals(finiteField.MulIdentity()) {
			t.Errorf("Multiplicative inverse not working")
		}

		// Distributivaty check

		if !a.Mul(b.Add(c)).Equals(a.Mul(b).Add(a.Mul(c))) {
			t.Errorf("Distributivity not working")
		}

	}

}
