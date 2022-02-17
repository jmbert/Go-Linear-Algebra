package linAlg

import (
	"errors"
	"math"
	"strconv"
)

type Vector struct {
	Fields []float64
}

func (v Vector) Add(v2 Vector) Vector {
	var out Vector
	if len(v.Fields) >= len(v2.Fields) {
		for i := 0; i < len(v.Fields); i++ {
			if i <= len(v2.Fields) {
				out.Fields = append(out.Fields, v.Fields[i]+v2.Fields[i])
			} else {
				out.Fields = append(out.Fields, v.Fields[i])
			}
		}
	} else {
		for i := 0; i < len(v2.Fields); i++ {
			if i <= len(v.Fields) {
				out.Fields = append(out.Fields, v.Fields[i]+v2.Fields[i])
			} else {
				out.Fields = append(out.Fields, v2.Fields[i])
			}
		}
	}
	return out
}

func (v Vector) Scale(scalar float64) Vector {
	var out Vector
	for i := 0; i < len(v.Fields); i++ {
		out.Fields = append(out.Fields, v.Fields[i]*scalar)
	}

	return out
}

func (v Vector) Magnitude() float64 {
	var sum float64
	for i := 0; i < len(v.Fields); i++ {
		sum += v.Fields[i] * v.Fields[i]
	}

	return math.Sqrt(sum)
}

func (v Vector) Normalise() Vector {
	return v.Scale(1 / v.Magnitude())
}

func (v Vector) Dot(v2 Vector) float64 {
	var out float64

	if len(v.Fields) >= len(v2.Fields) {
		for i := 0; i < len(v.Fields); i++ {
			if i <= len(v2.Fields) {
				out += v.Fields[i] * v2.Fields[i]
			}
		}
	} else {
		for i := 0; i < len(v2.Fields); i++ {
			if i <= len(v.Fields) {
				out += v.Fields[i] * v2.Fields[i]
			}
		}
	}

	return out
}

func (v Vector) Print() string {
	var printable string = "["
	for i := 0; i < len(v.Fields); i++ {
		if i != len(v.Fields)-1 {
			printable += strconv.FormatFloat((v.Fields[i]), 'f', 3, 64) + ", "
		} else {
			printable += strconv.FormatFloat((v.Fields[i]), 'f', 3, 64) + "]"
		}

	}

	return printable
}

func (v Vector) String() string {
	var printable string
	for i := 0; i < len(v.Fields); i++ {
		if i != len(v.Fields)-1 {
			printable += strconv.FormatFloat((v.Fields[i]), 'f', 3, 64) + "|"
		} else {
			printable += strconv.FormatFloat((v.Fields[i]), 'f', 3, 64)
		}

	}

	return printable
}

type Matrix struct {
	Columns []Vector
}

func (m Matrix) Add(m2 Matrix) Matrix {
	var out Matrix
	if len(m.Columns) >= len(m2.Columns) {
		for i := 0; i < len(m.Columns); i++ {
			if i <= len(m2.Columns) {
				out.Columns = append(out.Columns, m.Columns[i].Add(m2.Columns[i]))
			} else {
				break
			}
		}
	} else {
		for i := 0; i < len(m2.Columns); i++ {
			if i <= len(m.Columns) {
				out.Columns = append(out.Columns, m.Columns[i].Add(m2.Columns[i]))
			} else {
				break
			}
		}
	}
	return out
}

func (m Matrix) Transform(v Vector) Vector {
	var out Vector

	for i := 0; i < len(m.Columns[0].Fields); i++ {
		out.Fields = append(out.Fields, m.Row(i).Dot(v))
	}

	return out
}

func (m Matrix) Multiply(m2 Matrix) (Matrix, error) {
	var out Matrix

	if len(m.Columns) == len(m2.Columns[0].Fields) {
		for i := 0; i < len(m.Columns); i++ {
			var column = m.Transform(m2.Columns[i])
			out.Columns = append(out.Columns, column)
		}
		return out, nil
	} else {
		return Matrix{nil}, errors.New("invalid matrix-matrix pair")
	}
}

func (m Matrix) Scale(scalar float64) Matrix {
	var scaled Matrix

	for i := 0; i < len(m.Columns); i++ {
		scaled.Columns = append(scaled.Columns, m.Columns[i].Scale(scalar))
	}

	return scaled
}

func (m Matrix) SubMatrix(row, column int) Matrix {
	var subMatrix Matrix

	for j := 0; j < len(m.Columns); j++ {
		if j != row {
			subMatrix.Columns = append(subMatrix.Columns, Vector{remove(m.Columns[j].Fields, row)})
		}
	}

	return subMatrix
}

func (m Matrix) Determinant() float64 {
	var det = 0.0
	for j := 0; j < len(m.Columns); j++ {
		var subDet = m.SubMatrix(0, j).Determinant()
		if j%2 == 0 {
			det += subDet
		} else {
			det -= subDet
		}
	}
	return det
}

// Figure out eigenstuff nxn formula
//func (m Matrix) EigenNumbers() ([]float64, []Vector) {
//	if len(m.Columns) != len(m.Columns[0].Fields) {
//		return []float64{}, []Vector{}
//	} else {
//		var eigenValues []float64
//		var eigenVectors = []Vector{Origin(len(m.Columns))}
//
//		return eigenValues, eigenVectors
//	}
//}

func (m Matrix) Transpose() Matrix {
	var transposed Matrix

	for i := 0; i < len(m.Columns); i++ {
		transposed.Columns = append(transposed.Columns, m.Row(i))
	}

	return transposed
}

func (m Matrix) Cofactor() Matrix {
	var cofactor Matrix

	for column := 0; column < len(m.Columns); column++ {
		var c Vector
		for row := 0; row < len(m.Columns[column].Fields); row++ {
			c.Fields = append(c.Fields, m.SubMatrix(row, column).Determinant())
		}
		cofactor.Columns = append(cofactor.Columns, c)
	}

	return cofactor
}

func (m Matrix) Inverse() Matrix {
	return m.Cofactor().Transpose()
}

func (m Matrix) Row(row int) Vector {
	var R Vector

	for i := 0; i < len(m.Columns); i++ {
		R.Fields = append(R.Fields, m.Columns[i].Fields[row])
	}

	return R
}

func (m Matrix) Print() string {
	var printable string = "[\n"

	for i := 0; i < len(m.Columns[0].Fields); i++ {
		printable += m.Row(i).Print() + "\n"
	}

	printable += "]"

	return printable
}

func remove(slice []float64, s int) []float64 {
	return append(slice[:s], slice[s+1:]...)
}

func Origin(dimension int) Vector {
	var origin Vector
	for i := 0; i < dimension; i++ {
		origin.Fields = append(origin.Fields, 0)
	}
	return origin
}

func Identity(dimension int) Matrix {
	var Identity Matrix

	for i := 0; i < dimension; i++ {
		var column Vector
		for j := 0; j < dimension; j++ {
			if j == i {
				column.Fields = append(column.Fields, 1)
			} else {
				column.Fields = append(column.Fields, 0)
			}
		}
		Identity.Columns = append(Identity.Columns, column)
	}
	return Identity
}
