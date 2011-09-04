/*
Copyright 2011 Mostafa Hajizdeh

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pdf

import (
	"bytes"
	"testing"
)

type ot struct {
	name string
	in  interface{}
	out []byte
}

func TestOutput(t *testing.T) {
	// TODO add test with Persian text for string
	// TODO test escape characters in string
	// TODO test space and other special characters for names in dictionaries
	// TODO test empty stream
	tests := []ot{
		// simple types: null, boolean, numbers, strings
		{"null", nil, []byte("null")},
		{"false bool", false, []byte("false")},
		{"true bool", true, []byte("true")},
		{"negative ten", -10, []byte("-10")},
		{"negative float number", -2.4, []byte("-2.4")},
		{"negative one", -1.0, []byte("-1")},
		{"zero", 0.0, []byte("0")},
		{"small float", 0.23, []byte("0.23")},
		{"two", 2, []byte("2")},
		{"ten", 10, []byte("10")},
		{"empty string", "", []byte("()")},
		{"simple string", "hello", []byte("(hello)")},
		// arrays
		{"empty array", []int{}, []byte("[ ]")},
		{"array of one", []float64{1.1}, []byte("[ 1.1 ]")},
		{"array of two", []string{"a", "b"}, []byte("[ (a) (b) ]")},
		// dictionaries
		{"empty dictionary", map[string]int{}, []byte("<<\n>>")},
		{"dictionary of string",
			map[string]string{"k": "v"}, []byte("<<\n/k (v)\n>>")},
		{"dictionary of number",
			map[string]int{"A": 1}, []byte("<<\n/A 1\n>>")},
		// arrays of arras and dictionaries
		{"array including array and dictionary", []interface{}{
			[]int{1, 2},
			map[string]int{"a": 0},
		},
			[]byte("[ [ 1 2 ] <<\n/a 0\n>> ]"),
		},
		{"dictionary including array", map[string]interface{}{
			"a": []string{"p", "q"},
		},
			[]byte("<<\n/a [ (p) (q) ]\n>>"),
		},
		{"dictionary inclduing dictionary", map[string]interface{}{
			"d": map[string]int{"r": -1},
		},
			[]byte("<<\n/d <<\n/r -1\n>>\n>>"),
		},
		// streams
		{"bytes", []byte("ssss"),
			[]byte("<<\n/Length 4\n>>\nstream\nssss\nendstream")},
		{"buffer", bytes.NewBufferString("a"),
			[]byte("<<\n/Length 1\n>>\nstream\na\nendstream")},
	}

	for _, test := range tests {
		o := output(test.in)
		if bytes.Compare(o, test.out) != 0 {
			t.Errorf("%s: got\n\t%v\nexpected\n\t%v",
				test.name, o, test.out)
		}
	}
}
