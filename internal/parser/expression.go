package parser

import (
	"errors"
	"strings"
	"unicode"

	"CalcMatxGo/CalcMatxGo-main/internal/matrix"
)

var prec = map[string]int{"+": 1, "-": 1, "*": 2}

func tokenize(expr string) ([]string, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, errors.New("выражение пустое")
	}
	out := []string{}
	for i := 0; i < len(expr); {
		r := rune(expr[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		if unicode.IsLetter(r) {
			j := i
			for j < len(expr) && unicode.IsLetter(rune(expr[j])) {
				j++
			}
			out = append(out, expr[i:j])
			i = j
			continue
		}
		s := string(r)
		if s == "(" || s == ")" || s == "+" || s == "-" || s == "*" {
			out = append(out, s)
			i++
			continue
		}
		return nil, errors.New("выражение содержит недопустимые символы")
	}
	return out, nil
}

func toRPN(tokens []string) ([]string, error) {
	out, ops := []string{}, []string{}
	for _, t := range tokens {
		if unicode.IsLetter(rune(t[0])) {
			out = append(out, t)
			continue
		}
		if _, ok := prec[t]; ok {
			for len(ops) > 0 {
				top := ops[len(ops)-1]
				if p, ok := prec[top]; ok && p >= prec[t] {
					out = append(out, top)
					ops = ops[:len(ops)-1]
					continue
				}
				break
			}
			ops = append(ops, t)
			continue
		}
		if t == "(" {
			ops = append(ops, t)
			continue
		}
		if t == ")" {
			found := false
			for len(ops) > 0 {
				top := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				if top == "(" {
					found = true
					break
				}
				out = append(out, top)
			}
			if !found {
				return nil, errors.New("несогласованные скобки")
			}
		}
	}
	for i := len(ops) - 1; i >= 0; i-- {
		if ops[i] == "(" {
			return nil, errors.New("несогласованные скобки")
		}
		out = append(out, ops[i])
	}
	return out, nil
}

func Evaluate(expr string, vars map[string]matrix.Matrix) (matrix.Matrix, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}
	rpn, err := toRPN(tokens)
	if err != nil {
		return nil, err
	}
	st := []matrix.Matrix{}
	for _, t := range rpn {
		if unicode.IsLetter(rune(t[0])) {
			m, ok := vars[t]
			if !ok {
				return nil, errors.New("матрица " + t + " не задана")
			}
			st = append(st, m)
			continue
		}
		if len(st) < 2 {
			return nil, errors.New("некорректное выражение")
		}
		b, a := st[len(st)-1], st[len(st)-2]
		st = st[:len(st)-2]
		var res matrix.Matrix
		switch t {
		case "+":
			res, err = matrix.Add(a, b)
		case "-":
			res, err = matrix.Sub(a, b)
		case "*":
			res, err = matrix.Mul(a, b)
		default:
			return nil, errors.New("неизвестный оператор")
		}
		if err != nil {
			return nil, err
		}
		st = append(st, res)
	}
	if len(st) != 1 {
		return nil, errors.New("некорректное выражение")
	}
	return st[0], nil
}
