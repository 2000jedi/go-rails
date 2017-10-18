package template

func genHTML_example(content string, A string, xs []map[string]interface{}) (_gen string) {
	_gen = ``
	_gen += `<body>\n    `
	_gen += content
	_gen += `\n    `
	for _, x := range xs {
		x__id := x[`id`].(map[string]interface{})
		x__prim := x[`prim`].(map[string]interface{})
		x__id__first := x__id[`first`].(string)
		x__prim__first := x__prim[`first`].(string)
		x__id__second := x__id[`second`].(string)
		_gen += `\n        `
		_gen += x__id__first
		_gen += `\n        `
		_gen += x__prim__first
		_gen += `\n        `
		_gen += x__id__second
		_gen += `\n    `
	}
	_gen += `\n    `
	if A == "" {
		_gen += ` STH `
	}
	_gen += `\n</body>\n`

	return
}

func Construct_example(m map[string]interface{}) string {
	content := m[`content`].(string)
	A := m[`A`].(string)
	xs := m[`xs`].([]map[string]interface{})

	return genHTML_example(content, A, xs)
}
