package ui

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"calcMatx/internal/matrix"
)

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="ru">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width,initial-scale=1" />
<title>Калькулятор матриц и векторов</title>
<style>
body{font-family:Arial,sans-serif;max-width:1100px;margin:16px auto;padding:0 12px}
fieldset{border:1px solid #bcd;padding:12px;margin-bottom:12px}
textarea,input,select,button{font-size:14px}
textarea{width:100%;min-height:100px}
.row{display:grid;grid-template-columns:1fr 1fr;gap:12px}
.hidden{display:none}
pre{background:#f7f7f7;border:1px solid #ddd;padding:10px;white-space:pre-wrap}
.note{font-size:12px;color:#555}
</style>
<script>
function toggle(){
  const mode=document.getElementById('mode').value;
  for (const id of ['matrixBlock','slauBlock','vectorBlock']) document.getElementById(id).classList.add('hidden');
  document.getElementById(mode+'Block').classList.remove('hidden');
}
window.addEventListener('DOMContentLoaded', toggle);
</script>
</head>
<body>
<h1>Калькулятор матриц, СЛАУ и векторов</h1>
<p class="note">Вводите числа через пробел. Для матриц — по строкам, каждая строка на новой строке.</p>

<form method="POST">
<fieldset>
  <legend>Режим</legend>
  <select id="mode" name="mode" onchange="toggle()">
    <option value="matrix" {{if eq .Mode "matrix"}}selected{{end}}>Операции с матрицами</option>
    <option value="slau" {{if eq .Mode "slau"}}selected{{end}}>Решение СЛАУ</option>
    <option value="vector" {{if eq .Mode "vector"}}selected{{end}}>Операции с векторами</option>
  </select>
</fieldset>

<div id="matrixBlock">
<fieldset>
<legend>Матрицы</legend>
<div class="row">
<div><label>Матрица A</label><textarea name="A">{{.A}}</textarea></div>
<div><label>Матрица B</label><textarea name="B">{{.B}}</textarea></div>
</div>
<p>Операция:
<select name="matrix_op">
  <option value="add">A + B</option>
  <option value="sub">A - B</option>
  <option value="mul">A × B</option>
  <option value="scalar">A × число</option>
  <option value="transpose">Транспонировать A</option>
  <option value="det">Определитель A</option>
  <option value="rank">Ранг A</option>
</select>
</p>
<label>Число:</label> <input name="num" value="{{.Num}}" />
</fieldset>
</div>

<div id="slauBlock" class="hidden">
<fieldset>
<legend>СЛАУ A·x=b</legend>
<div class="row">
<div><label>Матрица A</label><textarea name="slau_A">{{.SlauA}}</textarea></div>
<div><label>Вектор b</label><textarea name="slau_b">{{.SlauB}}</textarea></div>
</div>
<p>Метод:
<select name="slau_op">
  <option value="gauss">Гаусса</option>
  <option value="cramer">Крамера</option>
  <option value="matrix">Матричный (через A⁻¹)</option>
</select>
</p>
</fieldset>
</div>

<div id="vectorBlock" class="hidden">
<fieldset>
<legend>Векторы</legend>
<div class="row">
<div><label>Вектор a</label><textarea name="va">{{.VA}}</textarea></div>
<div><label>Вектор b</label><textarea name="vb">{{.VB}}</textarea></div>
</div>
<label>Вектор c (для смешанного)</label><textarea name="vc">{{.VC}}</textarea>
<p>Операция:
<select name="vector_op">
  <option value="add">a + b</option>
  <option value="sub">a - b</option>
  <option value="scale">a × число</option>
  <option value="dot">a · b</option>
  <option value="cross">a × b (векторное)</option>
  <option value="triple">(a, b, c) смешанное</option>
</select>
</p>
<label>Число:</label> <input name="vnum" value="{{.VNum}}" />
</fieldset>
</div>

<button type="submit">Вычислить</button>
</form>

<h3>Промежуточные действия</h3>
<pre>{{.Steps}}</pre>
<h3>Результат</h3>
<pre>{{.Result}}</pre>
</body>
</html>`))

type Data struct {
	Mode, A, B, Num  string
	SlauA, SlauB     string
	VA, VB, VC, VNum string
	Result, Steps    string
}

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Data{
			Mode:  "matrix",
			A:     "1 2\n3 4",
			B:     "5 6\n7 8",
			Num:   "2",
			SlauA: "2 1\n5 7",
			SlauB: "11 13",
			VA:    "1 2 3",
			VB:    "4 5 6",
			VC:    "7 8 9",
			VNum:  "3",
		}
		if r.Method == "POST" {
			_ = r.ParseForm()
			data.Mode = r.FormValue("mode")
			data.A, data.B, data.Num = r.FormValue("A"), r.FormValue("B"), r.FormValue("num")
			data.SlauA, data.SlauB = r.FormValue("slau_A"), r.FormValue("slau_b")
			data.VA, data.VB, data.VC, data.VNum = r.FormValue("va"), r.FormValue("vb"), r.FormValue("vc"), r.FormValue("vnum")

			start := time.Now()
			result, steps, err := execute(data, r)
			elapsed := time.Since(start)
			if err != nil {
				data.Result = "Ошибка: " + err.Error()
			} else {
				data.Result = result
			}
			if elapsed > 10*time.Second {
				data.Result += "\nПредупреждение: вычисление заняло более 10 секунд."
			}
			if steps == "" {
				data.Steps = "Промежуточные действия не требуются для этой операции."
			} else {
				data.Steps = steps
			}
		}
		_ = page.Execute(w, data)
	})

	fmt.Println("server started: http://localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}

func execute(d Data, r *http.Request) (string, string, error) {
	switch d.Mode {
	case "matrix":
		A, err := matrix.Parse(d.A, "A")
		if err != nil {
			return "", "", err
		}
		op := r.FormValue("matrix_op")
		switch op {
		case "add", "sub", "mul":
			B, err := matrix.Parse(d.B, "B")
			if err != nil {
				return "", "", err
			}
			if op == "add" {
				m, e := matrix.Add(A, B)
				return m.String(), "Складываем матрицы A и B поэлементно.", e
			}
			if op == "sub" {
				m, e := matrix.Sub(A, B)
				return m.String(), "Вычитаем матрицу B из A поэлементно.", e
			}
			m, e := matrix.Mul(A, B)
			return m.String(), "Перемножаем матрицы по правилу строка×столбец.", e
		case "scalar":
			k, err := strconv.ParseFloat(strings.TrimSpace(d.Num), 64)
			if err != nil {
				return "", "", fmt.Errorf("неверное число")
			}
			m, e := matrix.Scalar(A, k)
			return m.String(), fmt.Sprintf("Умножаем каждый элемент A на %.8g.", k), e
		case "transpose":
			m, e := matrix.Transpose(A)
			return m.String(), "Меняем строки и столбцы местами.", e
		case "det":
			v, e := matrix.Determinant(A)
			return fmt.Sprintf("det(A)=%.8g", v), "Находим определитель методом исключения Гаусса.", e
		case "rank":
			v, e := matrix.Rank(A)
			return fmt.Sprintf("rank(A)=%d", v), "Приводим матрицу к ступенчатому виду и считаем ведущие строки.", e
		default:
			return "", "", fmt.Errorf("неизвестная операция")
		}
	case "slau":
		A, err := matrix.Parse(d.SlauA, "A")
		if err != nil {
			return "", "", err
		}
		b, err := matrix.ParseVector(d.SlauB, "b")
		if err != nil {
			return "", "", err
		}
		op := r.FormValue("slau_op")
		var x []float64
		var steps []string
		switch op {
		case "gauss":
			x, steps, err = matrix.SolveGauss(A, b)
		case "cramer":
			x, steps, err = matrix.SolveCramer(A, b)
		case "matrix":
			x, steps, err = matrix.SolveMatrixMethod(A, b)
		default:
			return "", "", fmt.Errorf("неизвестный метод")
		}
		return "x = " + matrix.FormatVector(x), strings.Join(steps, "\n"), err
	case "vector":
		a, err := matrix.ParseVector(d.VA, "a")
		if err != nil {
			return "", "", err
		}
		op := r.FormValue("vector_op")
		switch op {
		case "add", "sub", "dot", "cross":
			b, err := matrix.ParseVector(d.VB, "b")
			if err != nil {
				return "", "", err
			}
			if op == "add" {
				v, e := matrix.VectorAdd(a, b)
				return matrix.FormatVector(v), "Поэлементно складываем a и b.", e
			}
			if op == "sub" {
				v, e := matrix.VectorSub(a, b)
				return matrix.FormatVector(v), "Поэлементно вычитаем b из a.", e
			}
			if op == "dot" {
				v, e := matrix.Dot(a, b)
				return fmt.Sprintf("a·b=%.8g", v), "Умножаем соответствующие координаты и суммируем.", e
			}
			v, e := matrix.Cross(a, b)
			return matrix.FormatVector(v), "Вычисляем векторное произведение для 3D-векторов.", e
		case "triple":
			b, err := matrix.ParseVector(d.VB, "b")
			if err != nil {
				return "", "", err
			}
			c, err := matrix.ParseVector(d.VC, "c")
			if err != nil {
				return "", "", err
			}
			v, e := matrix.Triple(a, b, c)
			return fmt.Sprintf("(a,b,c)=%.8g", v), "Сначала a×b, затем скалярно умножаем на c.", e
		case "scale":
			k, err := strconv.ParseFloat(strings.TrimSpace(d.VNum), 64)
			if err != nil {
				return "", "", fmt.Errorf("неверное число")
			}
			v, e := matrix.VectorScale(a, k)
			return matrix.FormatVector(v), fmt.Sprintf("Умножаем координаты вектора a на %.8g.", k), e
		default:
			return "", "", fmt.Errorf("неизвестная операция")
		}
	default:
		return "", "", fmt.Errorf("неизвестный режим")
	}
}
