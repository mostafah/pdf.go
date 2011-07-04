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

type booleanTest struct {
	in  bool
	out []byte
}

var booleanTests = []booleanTest{
	booleanTest{false, []byte("false")},
	booleanTest{true, []byte("true")},
}

func TestBooleanNew(t *testing.T) {
	for _, bt := range booleanTests {
		b := newBoolean(bt.in)
		if bytes.Compare(b.toBytes(), bt.out) != 0 {
			t.Errorf("boolean: after newBoolean(%v), toBytes() = %v, want %v",
				bt.in, b.toBytes(), bt.out)
		}
	}
}

func TestBooleanSet(t *testing.T) {
	b := newBoolean(false)
	for _, bt := range booleanTests {
		b.set(bt.in)
		if bytes.Compare(b.toBytes(), bt.out) != 0 {
			t.Errorf("boolean: after set(%v), toBytes() = %v, want %v",
				bt.in, b.toBytes(), bt.out)
		}
	}
}

type numberTest struct {
	in  float64
	out []byte
}

type numberIntTest struct {
	in  int
	out []byte
}

var numberTests = []numberTest{
	numberTest{0.0, []byte("0")},
	numberTest{0.23, []byte("0.23")},
	numberTest{10, []byte("10")},
	numberTest{-2.4, []byte("-2.4")},
	numberTest{-1.0, []byte("-1")},
}

var numberIntTests = []numberIntTest{
	numberIntTest{0, []byte("0")},
	numberIntTest{2, []byte("2")},
	numberIntTest{-10, []byte("-10")},
}

func TestNumberNew(t *testing.T) {
	for _, nt := range numberTests {
		n := newNumber(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("number: after newNumber(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestNumberNewInt(t *testing.T) {
	for _, nt := range numberIntTests {
		n := newNumberInt(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("number: after newNumberInt(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestNumberSet(t *testing.T) {
	n := newNumber(0.0)
	for _, nt := range numberTests {
		n.set(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("number: after set(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestNumberSetInt(t *testing.T) {
	n := newNumber(0)
	for _, nt := range numberIntTests {
		n.setInt(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("number: after setInt(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

type strTest struct {
	in  string
	out []byte
}

// TODO add test with Persian text
// TODO escape characters
var strTests = []strTest{
	strTest{"", []byte("()")},
	strTest{"hello", []byte("(hello)")},
}

func TestStrNew(t *testing.T) {
	for _, st := range strTests {
		s := newStr(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("str: after newStr(%v), toBytes() = %v, want %v",
				st.in, s.toBytes(), st.out)
		}
	}
}

func TestStrSet(t *testing.T) {
	s := newStr("")
	for _, st := range strTests {
		s.set(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("str: after set(%v), toBytes() = %v, want %v",
				st.in, s.toBytes(), st.out)
		}
	}
}

type nameTest struct {
	in  string
	out []byte
}

var nameTests = []nameTest{
	nameTest{"hello", []byte("/hello")},
}

func TestNameNew(t *testing.T) {
	for _, nt := range nameTests {
		n := newName(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("name: after newName(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestNameSet(t *testing.T) {
	n := newName("")
	for _, nt := range nameTests {
		n.set(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("name: after set(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

type arrayTest struct {
	in  []object
	out []byte
}

var arrayTests = []arrayTest{
	arrayTest{[]object{}, []byte("[ ]")},
	arrayTest{[]object{newNumber(1.1)}, []byte("[ 1.1 ]")},
	arrayTest{[]object{newNumberInt(2), newStr("text")}, []byte("[ 2 (text) ]")},
}

func TestArray(t *testing.T) {
	for _, at := range arrayTests {
		a := newArray()
		for _, o := range at.in {
			a.add(o)
		}
		if bytes.Compare(a.toBytes(), at.out) != 0 {
			t.Errorf("array: toBytes() = %v, want %v",
				a.toBytes(), at.out)
		}
	}
}

type dictTest struct {
	inKey   []*name
	inValue []object
	out     []byte
}

var dictTests = []dictTest{
	dictTest{
		[]*name{},
		[]object{},
		// Output:
		// <<
		// >>
		[]byte("<<\n>>")},
	dictTest{
		[]*name{newName("VarOne")},
		[]object{newStr("value")},
		// Output:
		// <<
		// /VarOne (value)
		// >>
		[]byte("<<\n/VarOne (value)\n>>")},
	dictTest{
		[]*name{newName("VarOne"), newName("VarTwo")},
		[]object{newNumber(2.3), newNumberInt(0)},
		// Output:
		// <<
		// /VarOne 2.3
		// /VarTwo 0
		// >>
		[]byte("<<\n/VarOne 2.3\n/VarTwo 0\n>>")},
}

func TestDict(t *testing.T) {
	for _, dt := range dictTests {
		d := newDict()
		for i := 0; i < len(dt.inKey); i++ {
			d.add(dt.inKey[i], dt.inValue[i])
		}
		if bytes.Compare(d.toBytes(), dt.out) != 0 {
			t.Errorf("dict: toBytes() = %v, want %v",
				d.toBytes(), dt.out)
		}
	}
}

// TestDictMore tests dictionaries when they contain more dictionaries and arrays.
func TestDictMore(t *testing.T) {
	d := newDict()
	d.add(newName("N"), newNumber(1.0))

	d2 := newDict()
	d2.add(newName("Key"), newStr("value"))
	d.add(newName("D"), d2)

	a := newArray()
	a.add(newStr("array"))
	d.add(newName("A"), a)

	// Output:
	// <<
	// /N 1
	// /D <<
	// /Key (value)
	// >>
	// /A [ (array) ]
	// >>
	out := []byte("<<\n/N 1\n/D <<\n/Key (value)\n>>\n/A [ (array) ]\n>>")

	if bytes.Compare(d.toBytes(), out) != 0 {
		t.Errorf("dict: toBytes() = %v, want %v", d.toBytes(), out)
	}
}

type streamTest struct {
	in, out []byte
}

// TODO check PDF Reference for empty streams and add a test for that
var streamTests = []streamTest{
	streamTest{
		[]byte("ssss"),
		// Output:
		// <<
		// /Length 4
		// >>
		// stream
		// ssss
		// endstream
		[]byte("<<\n/Length 4\n>>\nstream\nssss\nendstream")},
}

func TestStream(t *testing.T) {
	for _, st := range streamTests {
		s := newStream(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("stream: toBytes() = %v, want %v",
				s.toBytes(), st.out)
		}
	}
}

func TestNull(t *testing.T) {
	in := newNull()
	out := []byte("null")
	if bytes.Compare(in.toBytes(), out) != 0 {
		t.Errorf("null: toBytes() = %v, want %v", in.toBytes(), out)
	}
}

type indirectTest struct {
	obj     object
	num     uint32
	offset  uint64
	out     []byte
	outBody []byte
	outRef  []byte
}

var indirectTests = []indirectTest{
	indirectTest{
		newStr("ssss"),
		1,
		100,
		[]byte("1 0 R"),
		[]byte("1 0 obj\n(ssss)\nendobj\n"),
		[]byte("0000000100 00000 n\r\n")},
	indirectTest{
		newNumber(12.34),
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
		i.setOffset(dt.offset)
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
