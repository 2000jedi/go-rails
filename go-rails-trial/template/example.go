package template

func genHTML_example(content string, xs []map[string]interface{}, A string, B string) (_gen string) {
	_gen = ``
	_gen += `<body>    `
	_gen += content
	_gen += `    `
	_gen += `    `
	for _, x := range xs {
		x__id := x[`id`].(map[string]interface{})
		x__prim := x[`prim`].(map[string]interface{})
		x__id__first := x__id[`first`].(string)
		x__prim__first := x__prim[`first`].(string)
		x__id__second := x__id[`second`].(string)
		_gen += `        `
		_gen += x__id__first
		_gen += `        `
		_gen += x__prim__first
		_gen += `        `
		_gen += x__id__second
		_gen += `    `
	}
	_gen += `    `
	if A == "" {
		_gen += ` STH `
	}
	_gen += `</body>`

	return
}

func Construct_example(m map[string]interface{}) string {
	A := m[`A`].(string)
	content := m[`content`].(string)
	B := m[`B`].(string)
	xs := m[`xs`].([]map[string]interface{})

	return genHTML_example(content, xs, A, B)
}
