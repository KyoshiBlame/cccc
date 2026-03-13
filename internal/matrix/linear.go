package matrix

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParseVector(text, name string) ([]float64, error) {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) == 0 {
		return nil, fmt.Errorf("%s: пустой вектор", name)
	}
	out := make([]float64, len(parts))
	for i, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: неверный формат числа", name)
		}
		out[i] = v
	}
	return out, nil
}

func SolveMatrixMethod(A Matrix, b []float64) ([]float64, []string, error) {
	if err := A.Validate("A"); err != nil {
		return nil, nil, err
	}
	if !A.isSquare() {
		return nil, nil, errors.New("матричный метод возможен только для квадратной матрицы A")
	}
	n := len(A)
	if len(b) != n {
		return nil, nil, errors.New("размер вектора b должен совпадать с числом строк A")
	}
	steps := []string{"1) Находим обратную матрицу A^-1."}
	inv, err := Inverse(A)
	if err != nil {
		return nil, steps, err
	}
	steps = append(steps, "A^-1:\n"+inv.String())
	bcol := make(Matrix, n)
	for i := 0; i < n; i++ {
		bcol[i] = []float64{b[i]}
	}
	steps = append(steps, "2) Умножаем X = A^-1 * b.")
	xm, err := Mul(inv, bcol)
	if err != nil {
		return nil, steps, err
	}
	x := make([]float64, n)
	for i := range x {
		x[i] = xm[i][0]
	}
	steps = append(steps, "X:\n"+xm.String())
	return x, steps, nil
}

func SolveCramer(A Matrix, b []float64) ([]float64, []string, error) {
	if err := A.Validate("A"); err != nil {
		return nil, nil, err
	}
	if !A.isSquare() {
		return nil, nil, errors.New("метод Крамера возможен только для квадратной матрицы A")
	}
	n := len(A)
	if len(b) != n {
		return nil, nil, errors.New("размер вектора b должен совпадать с числом строк A")
	}
	steps := []string{}
	detA, err := Determinant(A)
	if err != nil {
		return nil, steps, err
	}
	steps = append(steps, fmt.Sprintf("1) det(A) = %.8g", detA))
	if math.Abs(detA) < eps {
		return nil, steps, errors.New("метод Крамера неприменим: det(A)=0")
	}
	x := make([]float64, n)
	for c := 0; c < n; c++ {
		ac := A.Clone()
		for r := 0; r < n; r++ {
			ac[r][c] = b[r]
		}
		dc, err := Determinant(ac)
		if err != nil {
			return nil, steps, err
		}
		x[c] = dc / detA
		steps = append(steps, fmt.Sprintf("%d) det(A_%d)=%.8g -> x%d=%.8g", c+2, c+1, dc, c+1, x[c]))
	}
	return x, steps, nil
}

func SolveGauss(A Matrix, b []float64) ([]float64, []string, error) {
	if err := A.Validate("A"); err != nil {
		return nil, nil, err
	}
	rows, cols := A.Shape()
	if rows != cols {
		return nil, nil, errors.New("метод Гаусса реализован для квадратной матрицы A")
	}
	if len(b) != rows {
		return nil, nil, errors.New("размер вектора b должен совпадать с числом строк A")
	}
	n := rows
	aug := make(Matrix, n)
	for i := 0; i < n; i++ {
		aug[i] = append(append([]float64{}, A[i]...), b[i])
	}
	steps := []string{"1) Прямой ход."}
	for c := 0; c < n; c++ {
		p := pivotRow(aug, c, c)
		if math.Abs(aug[p][c]) < eps {
			return nil, steps, errors.New("система не имеет единственного решения")
		}
		if p != c {
			aug[p], aug[c] = aug[c], aug[p]
			steps = append(steps, fmt.Sprintf("перестановка строк R%d <-> R%d", c+1, p+1))
		}
		for r := c + 1; r < n; r++ {
			f := aug[r][c] / aug[c][c]
			for j := c; j <= n; j++ {
				aug[r][j] -= f * aug[c][j]
			}
			steps = append(steps, fmt.Sprintf("R%d = R%d - %.6g*R%d", r+1, r+1, f, c+1))
		}
	}
	steps = append(steps, "2) Обратный ход.")
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := aug[i][n]
		for j := i + 1; j < n; j++ {
			sum -= aug[i][j] * x[j]
		}
		if math.Abs(aug[i][i]) < eps {
			return nil, steps, errors.New("деление на ноль при обратном ходе")
		}
		x[i] = sum / aug[i][i]
		steps = append(steps, fmt.Sprintf("x%d = %.8g", i+1, x[i]))
	}
	return x, steps, nil
}

func FormatVector(v []float64) string {
	parts := make([]string, len(v))
	for i := range v {
		x := math.Round(v[i]*1e8) / 1e8
		if math.Abs(x) < eps {
			x = 0
		}
		parts[i] = fmt.Sprintf("%.8g", x)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
