package main

import (
	"fmt"
	"math/big"
)

type FFM struct {
	numRows     int
	numCols     int
	matrix      [][]*big.Int
	finiteField FiniteField
}

func NewFFM(finiteField FiniteField, length, width int) *FFM {
	f := new(FFM)
	f.finiteField = finiteField
	f.numRows = length
	f.numCols = width
	f.matrix = make([][]*big.Int, length)
	for i := range f.matrix {
		f.matrix[i] = make([]*big.Int, width)
		for j := range f.matrix[i] {
			f.matrix[i][j] = big.NewInt(0)
		}
	}
	return f
}

func (f FFM) print() {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			fmt.Print(f.matrix[i][j], " ")
		}
		fmt.Println()
	}
}

func (f FFM) sameSize(matrix *FFM) bool {
	if len(matrix.matrix) == f.numRows && len(matrix.matrix[0]) == f.numCols {
		return true
	}
	return false
}

func (f FFM) set(arr [][]*big.Int) error {
	if len(arr) != f.numRows || len(arr[0]) != f.numCols {
		return fmt.Errorf("matrix not of correct proportions")
	}
	for i := range arr {
		copy(f.matrix[i], arr[i])
	}
	return nil
}

func (f FFM) add(m *FFM) error {
	if !f.sameSize(m) {
		return fmt.Errorf("matricies are not of the same size")
	}
	for i := range f.matrix {
		for j := range f.matrix[i] {
			f.matrix[i][j] = f.finiteField.add(f.matrix[i][j], m.matrix[i][j])
		}
	}
	return nil
}

// returns new slice, however the values themselves still point to the same big Ints
func (f FFM) col(c int) []*big.Int {
	ans := make([]*big.Int, f.numRows)
	for i := range f.matrix {
		ans[i] = f.matrix[i][c]
	}
	return ans
}

func (f FFM) row(c int) []*big.Int {
	return f.matrix[c]
}

// dots 2 vectors
func (f FFM) dot(v1, v2 []*big.Int) *big.Int {
	ans := new(big.Int)
	ans.SetInt64(0)
	for i := range v1 {
		ans = f.finiteField.add(ans, f.finiteField.mul(v1[i], v2[i]))
	}
	return ans
}

func (f FFM) mul(m *FFM) error {
	if f.numCols != len(m.matrix) {
		return fmt.Errorf("the matrix attempted to add here is not of the same size, unfortunatly")
	}
	ans := make([][]*big.Int, f.numRows)
	for i := range f.matrix {
		row := make([]*big.Int, f.numCols)
		for j := range row {
			row[j] = f.dot(f.matrix[i], m.col(j))
		}
		ans[i] = row
	}
	f.set(ans)
	return nil
}

func (f FFM) addInverse() *FFM {
	ans := f.copy()
	for i := range ans.matrix {
		for j := range ans.matrix[i] {
			ans.matrix[i][j] = f.finiteField.addInverse(ans.matrix[i][j])
		}
	}
	return ans
}

func (f FFM) mulIdentity() *FFM {
	ans := make([][]*big.Int, len(f.matrix))
	for i := range ans {
		row := make([]*big.Int, f.numCols)
		for j := range row {
			row[j] = big.NewInt(0)
		}
		row[i] = big.NewInt(1)
		ans[i] = row
	}
	h := NewFFM(f.finiteField, f.numRows, f.numCols)
	h.set(ans)
	return h
}

func (f FFM) addIdentity() *FFM {
	return NewFFM(f.finiteField, f.numRows, f.numCols)
}

func (f FFM) mulInverse() *FFM {
	identity := f.mulIdentity()
	arr := make([][]*big.Int, len(f.matrix))
	for i := range f.matrix {
		arr[i] = make([]*big.Int, f.numCols)
		for j := range arr[i] {
			b := new(big.Int)
			b.Set(f.matrix[i][j])
			arr[i][j] = b
		}
	}
	var behind int
	for col := range arr[0] {
		for row := range arr {
			if row != col {
				behind = col
				c := f.finiteField.mul(arr[row][col], f.finiteField.mulInverse(arr[behind][col]))
				for j := range arr[row] {
					arr[row][j] = f.finiteField.add(arr[row][j], f.finiteField.addInverse(f.finiteField.mul(arr[behind][j], c)))
					identity.matrix[row][j] = f.finiteField.add(identity.matrix[row][j], f.finiteField.addInverse(f.finiteField.mul(identity.matrix[behind][j], c)))
				}
			}
		}
	}
	for i := range arr {
		v := f.finiteField.mulInverse(arr[i][i])
		for j := range arr[i] {
			arr[i][j] = f.finiteField.mul(arr[i][j], v)
			identity.matrix[i][j] = f.finiteField.mul(identity.matrix[i][j], v)
		}
	}
	return identity
}

func (f FFM) scale(s *big.Int) {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			f.matrix[i][j] = f.finiteField.mul(s, f.matrix[i][j])
		}
	}
}

func (f FFM) copy() *FFM{
	ans := NewFFM(f.finiteField, f.numRows, f.numCols)
	for i := range f.matrix {
		for j := range f.matrix[i] {
			ans.matrix[i][j] = f.matrix[i][j]
		}
	}
	return ans
}

func (f FFM) equals(m *FFM) bool {
	if !f.sameSize(m) {
		return false
	}
	for i := range f.matrix {
		for j := range f.matrix[i] {
			if f.matrix[i][j].Cmp(m.matrix[i][j]) != 0 {
				return false
			}
		}
	}
	return true
}