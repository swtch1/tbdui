package test

import "encoding/json"

func PrettyJSONSample() string {
	s := `
	{
		"foo:": "bar",
		"bar": [
			{"baz": "zaz"},
			{"bam": "zam"}
		],
		"zim": "zam"
	}`
	var x interface{}
	if err := json.Unmarshal([]byte(s), &x); err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
