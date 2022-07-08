package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	random "math/rand"
	"testing"
)

const MaxSize = 100

func randomFFM(prime *big.Int) *FFM {
	finiteField := NewGFP(prime)
	i := random.Intn(MaxSize)
	i += 1
	ffm := NewFFM(finiteField, i, i)

	for l := 0; l < ffm.numRows; l++ {
		for c := 0; c < ffm.numCols; c++ {

			element, _ := rand.Int(rand.Reader, prime)

			finiteField := NewGFP(prime)
			ffm.matrix[l][c], _ = finiteField.NewElement(element)

		}
	}
	return ffm
}

func TestFiniteFieldMatrix(t *testing.T) {

	// random.Seed(time.Now().Unix())
	random.Seed(0)
	
	prime, _ := rand.Prime(rand.Reader, 3)

	for i := 0; i < NumOfTests; i++ {
		ffm := randomFFM(prime)
		
		// Trying to do impossible things with addition and multiplication
		err := ffm.Add(NewFFM(ffm.finiteField, ffm.numRows+1, ffm.numCols+3))
		if err == nil {
			t.Errorf("Addition function is able to take matricies of different sizes")
		}
		err = ffm.Mul(NewFFM(ffm.finiteField, ffm.numRows+10, ffm.numCols+2))
		if err == nil {
			t.Errorf("Multiplication function takes in matricies that cannot be multiplied together")
		}
		k := ffm.Copy()
		ffm.Mul(ffm.MulIdentity())
		if !ffm.Equals(k) {
			t.Errorf("Multiplication identity not working")
		}
		ffm.Add(ffm.AddIdentity())
		if !ffm.Equals(k) {
			t.Errorf("Additive identity not working")
		}
		
		inv, err := ffm.MulInverse()
		if err == nil {
			ffm.Mul(inv)
			if !ffm.Equals(ffm.MulIdentity()) {
				t.Errorf("Multiplicative inverse not working")
			}
		}
		ffm = randomFFM(prime)
		ffm.Add(ffm.AddInverse())
		if !ffm.Equals(ffm.AddIdentity()) {
			t.Errorf("Additive inverse not working")
		}
		
	}
	fmt.Println("Done")
	// trying to invert a non invertable matrix
	m := NewFFMButEasier(7, [][]int{
		{1, 0},
		{1, 0},
	})
	m.MulInverse()
}
