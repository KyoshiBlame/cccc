package matrix

import "errors"

func VectorAdd(a, b []float64) ([]float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return nil, errors.New("векторы не должны быть пустыми")
	}
	if len(a) != len(b) {
		return nil, errors.New("для сложения длины векторов должны совпадать")
	}
	out := make([]float64, len(a))
	for i := range a {
		out[i] = a[i] + b[i]
	}
	return out, nil
}

func VectorSub(a, b []float64) ([]float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return nil, errors.New("векторы не должны быть пустыми")
	}
	if len(a) != len(b) {
		return nil, errors.New("для вычитания длины векторов должны совпадать")
	}
	out := make([]float64, len(a))
	for i := range a {
		out[i] = a[i] - b[i]
	}
	return out, nil
}

func VectorScale(v []float64, k float64) ([]float64, error) {
	if len(v) == 0 {
		return nil, errors.New("вектор не должен быть пустым")
	}
	out := append([]float64(nil), v...)
	for i := range out {
		out[i] *= k
	}
	return out, nil
}

func Dot(a, b []float64) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, errors.New("векторы не должны быть пустыми")
	}
	if len(a) != len(b) {
		return 0, errors.New("для скалярного произведения длины векторов должны совпадать")
	}
	s := 0.0
	for i := range a {
		s += a[i] * b[i]
	}
	return s, nil
}

func Cross(a, b []float64) ([]float64, error) {
	if len(a) != 3 || len(b) != 3 {
		return nil, errors.New("векторное произведение определено только для 3D-векторов")
	}
	return []float64{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}, nil
}

func Triple(a, b, c []float64) (float64, error) {
	ab, err := Cross(a, b)
	if err != nil {
		return 0, err
	}
	return Dot(ab, c)
}
