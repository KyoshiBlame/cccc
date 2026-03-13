package matrix

import (
	"errors"
	"math"
)

func pivotRow(m Matrix, start, col int) int {
	p := start
	mx := math.Abs(m[start][col])
	for i := start + 1; i < len(m); i++ {
		if v := math.Abs(m[i][col]); v > mx {
			mx = v
			p = i
		}
	}
	return p
}

func Determinant(m Matrix) (float64, error) {
	if err := m.Validate("Матрица"); err != nil {
		return 0, err
	}
	if !m.isSquare() {
		return 0, errors.New("определитель определён только для квадратной матрицы")
	}
	a := m.Clone()
	n := len(a)
	sign, det := 1.0, 1.0
	for c := 0; c < n; c++ {
		p := pivotRow(a, c, c)
		if math.Abs(a[p][c]) < eps {
			return 0, nil
		}
		if p != c {
			a[p], a[c] = a[c], a[p]
			sign *= -1
		}
		pivot := a[c][c]
		det *= pivot
		for r := c + 1; r < n; r++ {
			f := a[r][c] / pivot
			for j := c; j < n; j++ {
				a[r][j] -= f * a[c][j]
			}
		}
	}
	return det * sign, nil
}

func Rank(m Matrix) (int, error) {
	if err := m.Validate("Матрица"); err != nil {
		return 0, err
	}
	a := m.Clone()
	r, c := a.Shape()
	rank, row := 0, 0
	for col := 0; col < c && row < r; col++ {
		p := row
		for i := row + 1; i < r; i++ {
			if math.Abs(a[i][col]) > math.Abs(a[p][col]) {
				p = i
			}
		}
		if math.Abs(a[p][col]) < eps {
			continue
		}
		a[p], a[row] = a[row], a[p]
		pv := a[row][col]
		for j := col; j < c; j++ {
			a[row][j] /= pv
		}
		for i := 0; i < r; i++ {
			if i == row {
				continue
			}
			f := a[i][col]
			for j := col; j < c; j++ {
				a[i][j] -= f * a[row][j]
			}
		}
		row++
		rank++
	}
	return rank, nil
}

func RowEchelon(m Matrix) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	a := m.Clone()
	r, c := a.Shape()
	lead := 0
	for i := 0; i < r && lead < c; i++ {
		k := i
		for k < r && math.Abs(a[k][lead]) < eps {
			k++
		}
		if k == r {
			lead++
			i--
			continue
		}
		a[i], a[k] = a[k], a[i]
		div := a[i][lead]
		for j := 0; j < c; j++ {
			a[i][j] /= div
		}
		for rr := i + 1; rr < r; rr++ {
			f := a[rr][lead]
			for j := 0; j < c; j++ {
				a[rr][j] -= f * a[i][j]
			}
		}
		lead++
	}
	return a, nil
}

func Triangular(m Matrix) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	if !m.isSquare() {
		return nil, errors.New("треугольный вид доступен только для квадратной матрицы")
	}
	a := m.Clone()
	n := len(a)
	for c := 0; c < n; c++ {
		p := pivotRow(a, c, c)
		if math.Abs(a[p][c]) < eps {
			continue
		}
		a[p], a[c] = a[c], a[p]
		for r := c + 1; r < n; r++ {
			f := a[r][c] / a[c][c]
			for j := c; j < n; j++ {
				a[r][j] -= f * a[c][j]
			}
		}
	}
	return a, nil
}

func Inverse(m Matrix) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	if !m.isSquare() {
		return nil, errors.New("обратная матрица существует только для квадратной матрицы")
	}
	n := len(m)
	a := make(Matrix, n)
	for i := range a {
		a[i] = append(append([]float64{}, m[i]...), Identity(n)[i]...)
	}
	for c := 0; c < n; c++ {
		p := pivotRow(a, c, c)
		if math.Abs(a[p][c]) < eps {
			return nil, errors.New("матрица вырождена, обратной не существует")
		}
		a[p], a[c] = a[c], a[p]
		pv := a[c][c]
		for j := 0; j < 2*n; j++ {
			a[c][j] /= pv
		}
		for r := 0; r < n; r++ {
			if r == c {
				continue
			}
			f := a[r][c]
			for j := 0; j < 2*n; j++ {
				a[r][j] -= f * a[c][j]
			}
		}
	}
	out := make(Matrix, n)
	for i := range out {
		out[i] = append([]float64{}, a[i][n:]...)
	}
	return out, nil
}

func Power(m Matrix, exp int) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	if !m.isSquare() {
		return nil, errors.New("возведение в степень доступно только для квадратной матрицы")
	}
	if exp == 0 {
		return Identity(len(m)), nil
	}
	if exp < 0 {
		inv, err := Inverse(m)
		if err != nil {
			return nil, err
		}
		return Power(inv, -exp)
	}
	res := Identity(len(m))
	base := m.Clone()
	for exp > 0 {
		if exp%2 == 1 {
			var err error
			res, err = Mul(res, base)
			if err != nil {
				return nil, err
			}
		}
		var err error
		base, err = Mul(base, base)
		if err != nil {
			return nil, err
		}
		exp /= 2
	}
	return res, nil
}

func LU(m Matrix) (Matrix, Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, nil, err
	}
	if !m.isSquare() {
		return nil, nil, errors.New("LU-разложение доступно только для квадратной матрицы")
	}
	n := len(m)
	L, U := Identity(n), make(Matrix, n)
	for i := range U {
		U[i] = make([]float64, n)
	}
	for i := 0; i < n; i++ {
		for k := i; k < n; k++ {
			sum := 0.0
			for j := 0; j < i; j++ {
				sum += L[i][j] * U[j][k]
			}
			U[i][k] = m[i][k] - sum
		}
		if math.Abs(U[i][i]) < eps {
			return nil, nil, errors.New("LU-разложение невозможно без перестановок")
		}
		for k := i + 1; k < n; k++ {
			sum := 0.0
			for j := 0; j < i; j++ {
				sum += L[k][j] * U[j][i]
			}
			L[k][i] = (m[k][i] - sum) / U[i][i]
		}
	}
	return L, U, nil
}

func Elementary(m Matrix, op string, i, j int, factor float64) (Matrix, error) {
	if err := m.Validate("Матрица"); err != nil {
		return nil, err
	}
	if i < 0 || i >= len(m) || j < 0 || j >= len(m) {
		return nil, errors.New("индексы строк выходят за границы")
	}
	a := m.Clone()
	switch op {
	case "swap":
		a[i], a[j] = a[j], a[i]
	case "scale":
		for c := range a[i] {
			a[i][c] *= factor
		}
	case "add":
		for c := range a[i] {
			a[i][c] += factor * a[j][c]
		}
	default:
		return nil, errors.New("неизвестное элементарное преобразование")
	}
	return a, nil
}
