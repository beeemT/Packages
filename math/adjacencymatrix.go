package math

import "fmt"

type Matrix interface {
	MinPlusExp(int)
	ConvertToMatrixF64() *Matrix
	ConvertToMatrixI() *Matrix
}

//Matrix is a matrix of type float64
type MatrixF64 struct {
	matrix [][]float64
}

type MatrixI struct {
	matrix [][]int
}

/*
func (m *MatrixF64) MinPlusExp(n int) *Matrix {

}
*/

func (m *MatrixI) MinPlusExp(n int) {
	temp := &(*m)
	fmt.Println(&temp)
	fmt.Println(&m)
	for c := 0; c < n; c++ {
		for i, line := range m.matrix {
			for j := range line {
				val := 0
				for k := 0; k < len(m.matrix); k++ {
					if k == 0 {
						val = k
						continue
					}
					dn := m.matrix[i][k] + temp.matrix[k][j]
					if dn < val {
						val = dn
					}
				}
				m.matrix[i][j] = val
			}

		}
	}
}

/*func (m *MatrixF64) ConvertToMatrixF64() *Matrix {

}

func (m *MatrixI) ConvertToMatrixF64() *Matrix {

}

func (m *MatrixF64) ConvertToMatrixI() *Matrix {

}

func (m *MatrixI) ConvertToMatrixI() *Matrix {

}

func makeEmptyMatrixF64(len int) *MatrixF64 {

}*/

func makeEmptyMatrixI(length int) [][]int {
	matrix := make([][]int, length)
	for i := range matrix {
		matrix[i] = make([]int, length)
	}
	return matrix
}

func main() {
	matrix := MatrixI{matrix: makeEmptyMatrixI(4)}
	matrix.matrix = [][]int{[]int{0, 1, 111, 111},
		[]int{1, 0, 2, 5},
		[]int{111, 2, 0, 2},
		[]int{111, 5, 2, 0},
	}
	matrix.MinPlusExp(2)
	fmt.Println(matrix)
}
