package main

import (
	"errors"
	"fmt"
)

type FFM struct {
	numRows     int
	numCols     int
	matrix      [][]Element
	finiteField FiniteField
}

func PrintMatrix(m [][]Element) {
	for i := range m {
		for j := range m[i] {
			fmt.Print(m[i][j].Val(), " ")
		}
		fmt.Println()
	}
}

func NewFFM(finiteField FiniteField, length, width int) *FFM {
	f := new(FFM)
	f.finiteField = finiteField
	f.numRows = length
	f.numCols = width
	f.matrix = make([][]Element, length)
	for i := range f.matrix {
		f.matrix[i] = make([]Element, width)
	}
	return f
}

func NewFFMButEasier(prime int64, arr [][]int) *FFM {
	f := NewFFM(NewGFPButEasier(prime), len(arr), len(arr[0]))
	f.SetButEasier(arr)
	return f
}

func (f FFM) Print() {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			if f.matrix[i][j] == nil {
				fmt.Print("<nil>", " ")
			} else {
				fmt.Print(f.matrix[i][j].Val(), " ")
			}
		}
		fmt.Println()
	}
}

func (f FFM) SameSize(matrix *FFM) bool {
	if matrix.numRows == 0 && f.numRows == 0  {
		return true
	}
	if len(matrix.matrix) == f.numRows && len(matrix.matrix[0]) == f.numCols {
		return true
	}
	return false
}

func (f FFM) Set(arr [][]Element) error {
	if len(arr) != f.numRows || len(arr[0]) != f.numCols {
		return errors.New("matrix not of correct proportions")
	}
	for i := range arr {
		for j := range arr[i] {
			if !f.finiteField.Contains(arr[i][j]) {
				return errors.New("matrix is not compatible with Finite Field")
			}
			f.matrix[i][j] = arr[i][j]
		}
	}
	return nil
}

func (f FFM) SetButEasier(arr [][]int) error {
	arrButWayMoreComplicated := make([][]Element, len(arr))
	for i := range arrButWayMoreComplicated {
		arrButWayMoreComplicated[i] = make([]Element, len(arr[0]))
		for j := range arrButWayMoreComplicated[i] {
			num := f.finiteField.(*GFP).NewElementButEasier(arr[i][j])
			arrButWayMoreComplicated[i][j] = num
		}
	}
	return f.Set(arrButWayMoreComplicated)
}

func (f FFM) Add(m *FFM) error {
	if !f.SameSize(m) {
		return fmt.Errorf("matricies are not of the same size")
	}
	for i := range f.matrix {
		for j := range f.matrix[i] {
			f.matrix[i][j] = f.matrix[i][j].Add(m.matrix[i][j])
		}
	}
	return nil
}

// returns new slice, however the values themselves still point to the same elements
func (f FFM) Col(c int) []Element {
	ans := make([]Element, f.numRows)
	for i := range f.matrix {
		ans[i] = f.matrix[i][c]
	}
	return ans
}

func (f FFM) Row(c int) []Element {
	return f.matrix[c]
}

// dots 2 vectors
func (f FFM) Dot(v1, v2 []Element) Element {
	ans := f.finiteField.AddIdentity()
	for i := range v1 {
		ans = ans.Add(v1[i].Mul(v2[i]))
	}
	return ans
}

func (f FFM) Mul(m *FFM) error {
	if f.numCols != len(m.matrix) {
		return fmt.Errorf("the matrix attempted to multiply here is not of the same size, unfortunatly")
	}
	ans := make([][]Element, f.numRows)
	for i := range f.matrix {
		row := make([]Element, f.numCols)
		for j := range row {
			row[j] = f.Dot(f.matrix[i], m.Col(j))
		}
		ans[i] = row
	}
	if len(ans) != 0 {
		f.Set(ans)
	}
	return nil
}

func (f FFM) AddInverse() *FFM {
	ans := f.Copy()
	for i := range ans.matrix {
		for j := range ans.matrix[i] {
			ans.matrix[i][j] = ans.matrix[i][j].AddInverse()
		}
	}
	return ans
}

func (f FFM) MulIdentity() *FFM {
	ans := NewFFM(f.finiteField, f.numRows, f.numCols)
	for i := range ans.matrix {
		for j := range ans.matrix[i] {
			ans.matrix[i][j] = f.finiteField.AddIdentity()
		}
		if len(ans.matrix[i]) > i {
			ans.matrix[i][i] = f.finiteField.MulIdentity()
		}
	}
	return ans
}

func (f FFM) AddIdentity() *FFM {
	ffm := NewFFM(f.finiteField, f.numRows, f.numCols)
	for i := range f.matrix {
		for j := range f.matrix[i] {
			ffm.matrix[i][j] = ffm.finiteField.AddIdentity()
		}
	}
	return ffm
}

func (f FFM) MulInverse() (*FFM, error) {
	identity := f.MulIdentity()
	var arr [][]Element

	arr = make([][]Element, len(f.matrix))
	for i := range f.matrix {
		arr[i] = make([]Element, f.numCols)
		for j := range arr[i] {
			// b := new(big.Int)
			// b.Set(f.matrix[i][j])
			// arr[i][j] = b
			arr[i][j] = f.matrix[i][j]
		}
	}
	for col := range arr[0] {
		if len(arr[col]) > col && arr[col][col].Equals(f.finiteField.AddIdentity()) {
			for k := col + 1; k < len(arr); k++ {
				if arr[col][k].Equals(f.finiteField.AddIdentity()) {
					tmp := arr[col]
					tmpIdentity := identity.matrix[col]
					arr[col] = arr[k]
					arr[k] = tmp
					identity.matrix[col] = identity.matrix[k]
					identity.matrix[k] = tmpIdentity
					break
				}
			}
			if arr[col][col].Equals(f.finiteField.AddIdentity()) {
				return nil, errors.New("matrix does not have an inverse")
			}
		}
		for row := range arr {
			if row != col {
				c := arr[row][col].Mul(arr[col][col].MulInverse())
				for j := range arr[row] {
					arr[row][j] = arr[row][j].Add(arr[col][j].Mul(c).AddInverse())
					identity.matrix[row][j] = identity.matrix[row][j].Add(identity.matrix[col][j].Mul(c).AddInverse())
				}
			}
		}
	}
	for i := range arr {
		v := arr[i][i].MulInverse()
		
		for j := range arr[i] {
			arr[i][j] = arr[i][j].Mul(v)
			identity.matrix[i][j] = identity.matrix[i][j].Mul(v)
		}
	}
	return identity, nil
}

func (f FFM) Scale(s Element) {
	for i := range f.matrix {
		for j := range f.matrix[i] {
			f.matrix[i][j] = s.Mul(f.matrix[i][j])
		}
	}
}

func (f FFM) Copy() *FFM {
	ans := NewFFM(f.finiteField, f.numRows, f.numCols)
	for i := range f.matrix {
		for j := range f.matrix[i] {
			ans.matrix[i][j] = f.matrix[i][j]
		}
	}
	return ans
}

func (f FFM) Equals(m *FFM) bool {
	if !f.SameSize(m) {
		return false
	}
	for i := range f.matrix {
		for j := range f.matrix[i] {
			if !f.matrix[i][j].Equals(m.matrix[i][j]) {
				return false
			}
		}
	}
	return true
}
