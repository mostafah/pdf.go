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

type indirectTest struct {
	obj     pObject
	num     int
	off     int
	out     []byte
	outBody []byte
	outRef  []byte
}

var indirectTests = []indirectTest{
	indirectTest{
		obj("ssss"),
		1,
		100,
		[]byte("1 0 R"),
		[]byte("1 0 obj\n(ssss)\nendobj\n"),
		[]byte("0000000100 00000 n\r\n")},
	indirectTest{
		obj(12.34),
		3,
		45,
		[]byte("3 0 R"),
		[]byte("3 0 obj\n12.34\nendobj\n"),
		[]byte("0000000045 00000 n\r\n")},
}

func TestIndirect(t *testing.T) {
	for _, dt := range indirectTests {
		i := newIndirect(dt.obj)
		i.setNum(dt.num)
		i.setOffset(dt.off)
		if bytes.Compare(i.toBytes(), dt.out) != 0 {
			t.Errorf("indirect: toBytes() = %v, want %v",
				i.toBytes(), dt.out)
		}
		if bytes.Compare(i.body(), dt.outBody) != 0 {
			t.Errorf("indirect: body() = %v, want %v",
				i.body(), dt.outBody)
		}
		if bytes.Compare(i.ref(), dt.outRef) != 0 {
			t.Errorf("indirect: ref() = %v, want %v",
				i.ref(), dt.outRef)
		}
	}
}

// TODO merge these two and other similar pairs.
func TestIndirectSet(t *testing.T) {
	for _, dt := range indirectTests {
		i := newIndirect(newPNull())
		i.set(dt.obj)
		i.setNum(dt.num)
		i.setOffset(dt.off)
		if bytes.Compare(i.toBytes(), dt.out) != 0 {
			t.Errorf("indirect: toBytes() = %v, want %v",
				i.toBytes(), dt.out)
		}
		if bytes.Compare(i.body(), dt.outBody) != 0 {
			t.Errorf("indirect: body() = %v, want %v",
				i.body(), dt.outBody)
		}
		if bytes.Compare(i.ref(), dt.outRef) != 0 {
			t.Errorf("indirect: ref() = %v, want %v",
				i.ref(), dt.outRef)
		}
	}
}

type ot struct {
	in  interface{}
	out []byte
}

func TestObj(t *testing.T) {
	// TODO add test with Persian text for string
	// TODO test escape characters in string
	// TODO test space and other special characters for names in dictionaries
	// TODO test empty stream
	ots := []ot{
		// simple types: null, boolean, numbers, strings
		ot{nil, []byte("null")},
		ot{false, []byte("false")},
		ot{true, []byte("true")},
		ot{-10, []byte("-10")},
		ot{-2.4, []byte("-2.4")},
		ot{-1.0, []byte("-1")},
		ot{0.0, []byte("0")},
		ot{0.23, []byte("0.23")},
		ot{2, []byte("2")},
		ot{10, []byte("10")},
		ot{"", []byte("()")},
		ot{"hello", []byte("(hello)")},
		// arrays
		ot{[]int{}, []byte("[ ]")},
		ot{[]float64{1.1}, []byte("[ 1.1 ]")},
		ot{[]string{"a", "b"}, []byte("[ (a) (b) ]")},
		// dictionaries
		ot{map[string]int{}, []byte("<<\n>>")},
		ot{map[string]string{"k": "v"}, []byte("<<\n/k (v)\n>>")},
		ot{map[string]int{"A": 1}, []byte("<<\n/A 1\n>>")},
		// arrays of arras and dictionaries
		ot{[]interface{}{
			[]int{1, 2},
			map[string]int{"a": 0},
		},
			[]byte("[ [ 1 2 ] <<\n/a 0\n>> ]"),
		},
		ot{map[string]interface{}{
			"a": []string{"p", "q"},
		},
			[]byte("<<\n/a [ (p) (q) ]\n>>"),
		},
		ot{map[string]interface{}{
			"d": map[string]int{"r": -1},
		},
			[]byte("<<\n/d <<\n/r -1\n>>\n>>"),
		},
		// streams
		ot{[]byte("ssss"),
			[]byte("<<\n/Length 4\n>>\nstream\nssss\nendstream")},
	}

	for _, pt := range ots {
		o := obj(pt.in)
		if bytes.Compare(o.toBytes(), pt.out) != 0 {
			t.Errorf("obj: after obj(%v), toBytes() = %q, "+
				"want %q", pt.in, o.toBytes(), pt.out)
		}
	}
}
