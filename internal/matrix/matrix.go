package matrix

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const eps = 1e-10

type Matrix [][]float64

func (m Matrix) Clone() Matrix {
	out := make(Matrix, len(m))
	for i := range m {
		out[i] = append([]float64(nil), m[i]...)
	}
	return out
}

func (m Matrix) Validate(name string) error {
	if len(m) == 0 {
		return fmt.Errorf("%s: пустая матрица", name)
	}
	cols := len(m[0])
	if cols == 0 {
		return fmt.Errorf("%s: нет столбцов", name)
	}
	for i, row := range m {
		if len(row) != cols {
			return fmt.Errorf("%s: строка %d имеет неверную длину", name, i+1)
		}
		for _, v := range row {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return fmt.Errorf("%s: содержит нечисловые значения", name)
			}
		}
	}
	return nil
}

func Parse(text, name string) (Matrix, error) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	rows := make(Matrix, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		row := make([]float64, len(parts))
		for i, p := range parts {
			v, err := strconv.ParseFloat(p, 64)
			if err != nil {
				return nil, fmt.Errorf("%s: неверный формат числа", name)
			}
			row[i] = v
		}
		rows = append(rows, row)
	}
	if err := rows.Validate(name); err != nil {
		return nil, err
	}
	return rows, nil
}

func (m Matrix) String() string {
	var b strings.Builder
	for i, row := range m {
		for j, v := range row {
			if j > 0 {
				b.WriteString("\t")
			}
			x := math.Round(v*1e8) / 1e8
			if math.Abs(x) < eps {
				x = 0
			}
			b.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
		}
		if i < len(m)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Matrix) Shape() (int, int) { return len(m), len(m[0]) }
func (m Matrix) isSquare() bool    { r, c := m.Shape(); return r == c }

func Identity(n int) Matrix {
	out := make(Matrix, n)
	for i := 0; i < n; i++ {
		out[i] = make([]float64, n)
		out[i][i] = 1
	}
	return out
}

func Add(a, b Matrix) (Matrix, error) {
	if err := a.Validate("A"); err != nil {
		return nil, err
	}
	if err := b.Validate("B"); err != nil {
		return nil, err
	}
	ra, ca := a.Shape()
	rb, cb := b.Shape()
	if ra != rb || ca != cb {
		return nil, errors.New("сложение: размеры матриц должны совпадать")
	}
	out := make(Matrix, ra)
	for i := range out {
		out[i] = make([]float64, ca)
		for j := 0; j < ca; j++ {
			out[i][j] = a[i][j] + b[i][j]
		}
	}
	return out, nil
}

func Sub(a, b Matrix) (Matrix, error) {
	if err := a.Validate("A"); err != nil {
		return nil, err
	}
	if err := b.Validate("B"); err != nil {
		return nil, err
	}
	ra, ca := a.Shape()
	rb, cb := b.Shape()
	if ra != rb || ca != cb {
		return nil, errors.New("вычитание: размеры матриц должны совпадать")
	}
	out := make(Matrix, ra)
	for i := range out {
		out[i] = make([]float64, ca)
		for j := 0; j < ca; j++ {
			out[i][j] = a[i][j] - b[i][j]
		}
	}
	return out, nil
}

func Mul(a, b Matrix) (Matrix, error) {
	if err := a.Validate("A"); err != nil {
		return nil, err
	}
	if err := b.Validate("B"); err != nil {
		return nil, err
	}
	ra, ca := a.Shape()
	rb, cb := b.Shape()
	if ca != rb {
		return nil, errors.New("умножение: число столбцов A должно равняться числу строк B")
	}
	out := make(Matrix, ra)
	for i := range out {
		out[i] = make([]float64, cb)
		for k := 0; k < ca; k++ {
			for j := 0; j < cb; j++ {
				out[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return out, nil
}

func Scalar(m Matrix, k float64) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	out := m.Clone()
	for i := range out {
		for j := range out[i] {
			out[i][j] *= k
		}
	}
	return out, nil
}

func Transpose(m Matrix) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	r, c := m.Shape()
	out := make(Matrix, c)
	for j := 0; j < c; j++ {
		out[j] = make([]float64, r)
		for i := 0; i < r; i++ {
			out[j][i] = m[i][j]
		}
	}
	return out, nil
}

func Trace(m Matrix) (float64, error) {
	if err := m.Validate("Матрица"); err != nil {
		return 0, err
	}
	if !m.isSquare() {
		return 0, errors.New("след определён только для квадратной матрицы")
	}
	s := 0.0
	for i := range m {
		s += m[i][i]
	}
	return s, nil
}
