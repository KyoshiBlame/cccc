package ui

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"calcMatx/internal/matrix"
)

var page = template.Must(template.New("page").Parse(`
<html>
<head>
<title>Matrix Calculator</title>
<style>
textarea {width:300px;height:120px}
</style>
</head>

<body>

<h2>Matrix Calculator</h2>

<form method="POST">

<p>Matrix A</p>
<textarea name="A">{{.A}}</textarea>

<p>Operation</p>
<select name="op">
<option value="det">Determinant</option>
<option value="transpose">Transpose</option>
<option value="rank">Rank</option>
<option value="trace">Trace</option>
<option value="power">Power</option>
<option value="scalar">Scalar</option>
</select>

<p>Number (power/scalar)</p>
<input name="num" value="2">

<br><br>
<button type="submit">Calculate</button>

</form>

<h3>Result</h3>
<pre>{{.Result}}</pre>

</body>
</html>
`))

type Data struct {
	A      string
	Result string
}

func Start() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := Data{
			A: "1 2 3\n0 1 4\n5 6 0",
		}

		if r.Method == "POST" {

			A, err := matrix.Parse(r.FormValue("A"), "A")
			if err != nil {
				data.Result = err.Error()
				page.Execute(w, data)
				return
			}

			op := r.FormValue("op")
			num := r.FormValue("num")

			switch op {

			case "det":

				v, err := matrix.Determinant(A)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = fmt.Sprintf("det = %v", v)
				}

			case "transpose":

				m, err := matrix.Transpose(A)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = m.String()
				}

			case "rank":

				v, err := matrix.Rank(A)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = fmt.Sprintf("rank = %d", v)
				}

			case "trace":

				v, err := matrix.Trace(A)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = fmt.Sprintf("trace = %v", v)
				}

			case "power":

				n, _ := strconv.Atoi(num)

				m, err := matrix.Power(A, n)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = m.String()
				}

			case "scalar":

				k, _ := strconv.ParseFloat(num, 64)

				m, err := matrix.Scalar(A, k)
				if err != nil {
					data.Result = err.Error()
				} else {
					data.Result = m.String()
				}
			}
		}

		page.Execute(w, data)
	})

	fmt.Println("server started: http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
