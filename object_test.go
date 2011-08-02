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

type pBooleanTest struct {
	in  bool
	out []byte
}

var pBooleanTests = []pBooleanTest{
	pBooleanTest{false, []byte("false")},
	pBooleanTest{true, []byte("true")},
}

func TestPBooleanNew(t *testing.T) {
	for _, bt := range pBooleanTests {
		b := newPBoolean(bt.in)
		if bytes.Compare(b.toBytes(), bt.out) != 0 {
			t.Errorf("pBoolean: after newPBoolean(%v), toBytes() ="+
				" %v, want %v", bt.in, b.toBytes(), bt.out)
		}
	}
}

func TestPBooleanSet(t *testing.T) {
	b := newPBoolean(false)
	for _, bt := range pBooleanTests {
		b.set(bt.in)
		if bytes.Compare(b.toBytes(), bt.out) != 0 {
			t.Errorf("pBoolean: after set(%v), toBytes() = %v, want %v",
				bt.in, b.toBytes(), bt.out)
		}
	}
}

type pNumberTest struct {
	in  float64
	out []byte
}

type pNumberIntTest struct {
	in  int
	out []byte
}

var pNumberTests = []pNumberTest{
	pNumberTest{0.0, []byte("0")},
	pNumberTest{0.23, []byte("0.23")},
	pNumberTest{10, []byte("10")},
	pNumberTest{-2.4, []byte("-2.4")},
	pNumberTest{-1.0, []byte("-1")},
}

var pNumberIntTests = []pNumberIntTest{
	pNumberIntTest{0, []byte("0")},
	pNumberIntTest{2, []byte("2")},
	pNumberIntTest{-10, []byte("-10")},
}

func TestPNumberNew(t *testing.T) {
	for _, nt := range pNumberTests {
		n := newPNumber(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pNumber: after newPNumber(%v), toBytes() = %v,"+
				" want %v", nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestPNumberNewInt(t *testing.T) {
	for _, nt := range pNumberIntTests {
		n := newPNumberInt(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pNumber: after newPNumberInt(%v), toBytes() ="+
				" %v, want %v", nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestPNumberSet(t *testing.T) {
	n := newPNumber(0.0)
	for _, nt := range pNumberTests {
		n.set(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pNumber: after set(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestPNumberSetInt(t *testing.T) {
	n := newPNumber(0)
	for _, nt := range pNumberIntTests {
		n.setInt(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pNumber: after setInt(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

type pStringTest struct {
	in  string
	out []byte
}

// TODO add test with Persian text
// TODO escape characters
var pStringTests = []pStringTest{
	pStringTest{"", []byte("()")},
	pStringTest{"hello", []byte("(hello)")},
}

func TestPStringNew(t *testing.T) {
	for _, st := range pStringTests {
		s := newPString(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("pString: after newPString(%v), toBytes() = %v,"+
				" want %v", st.in, s.toBytes(), st.out)
		}
	}
}

func TestPStringSet(t *testing.T) {
	s := newPString("")
	for _, st := range pStringTests {
		s.set(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("pString: after set(%v), toBytes() = %v, want %v",
				st.in, s.toBytes(), st.out)
		}
	}
}

type pNameTest struct {
	in  string
	out []byte
}

// TODO add tests with space and other special characters.
var pNameTests = []pNameTest{
	pNameTest{"hello", []byte("/hello")},
}

func TestPNameNew(t *testing.T) {
	for _, nt := range pNameTests {
		n := newPName(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pName: after newPName(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

func TestPNameSet(t *testing.T) {
	n := newPName("")
	for _, nt := range pNameTests {
		n.set(nt.in)
		if bytes.Compare(n.toBytes(), nt.out) != 0 {
			t.Errorf("pName: after set(%v), toBytes() = %v, want %v",
				nt.in, n.toBytes(), nt.out)
		}
	}
}

type pArrayTest struct {
	in  []pObject
	out []byte
}

var pArrayTests = []pArrayTest{
	pArrayTest{[]pObject{}, []byte("[ ]")},
	pArrayTest{[]pObject{newPNumber(1.1)}, []byte("[ 1.1 ]")},
	pArrayTest{[]pObject{newPNumberInt(2), newPString("text")},
		[]byte("[ 2 (text) ]")},
}

func TestPArray(t *testing.T) {
	for _, at := range pArrayTests {
		a := newPArray()
		for _, o := range at.in {
			a.add(o)
		}
		if bytes.Compare(a.toBytes(), at.out) != 0 {
			t.Errorf("pArray: toBytes() = %v, want %v",
				a.toBytes(), at.out)
		}
	}
}

type pDictTest struct {
	inKey []string
	inVal []pObject
	out   []byte
}

var pDictTests = []pDictTest{
	pDictTest{
		[]string{},
		[]pObject{},
		// Output:
		// <<
		// >>
		[]byte("<<\n>>")},
	pDictTest{
		[]string{"VarOne"},
		[]pObject{newPString("value")},
		// Output:
		// <<
		// /VarOne (value)
		// >>
		[]byte("<<\n/VarOne (value)\n>>")},
	pDictTest{
		[]string{"VarOne", "VarTwo"},
		[]pObject{newPNumber(2.3), newPNumberInt(0)},
		// Output:
		// <<
		// /VarOne 2.3
		// /VarTwo 0
		// >>
		[]byte("<<\n/VarOne 2.3\n/VarTwo 0\n>>")},
}

func TestPDict(t *testing.T) {
	for _, dt := range pDictTests {
		d := newPDict()
		for i := 0; i < len(dt.inKey); i++ {
			d.put(dt.inKey[i], dt.inVal[i])
		}
		if bytes.Compare(d.toBytes(), dt.out) != 0 {
			t.Errorf("pDict: toBytes() = %v, want %v",
				d.toBytes(), dt.out)
		}
	}
}

func TestPDictDupl(t *testing.T) {
	for _, dt := range pDictTests {
		d := newPDict()
		for i := 0; i < len(dt.inKey); i++ {
			d.put(dt.inKey[i], newPNull())
		}
		for i := 0; i < len(dt.inKey); i++ {
			d.put(dt.inKey[i], dt.inVal[i])
		}
		if bytes.Compare(d.toBytes(), dt.out) != 0 {
			t.Errorf("pDict duplicate: toBytes() = %q, want %q",
				d.toBytes(), dt.out)
		}
	}
}

// TestDictMore tests dictionaries when they contain other dictionaries and arrays.
func TestPDictMore(t *testing.T) {
	d := newPDict()
	d.put("N", newPNumber(1.0))

	d2 := newPDict()
	d2.put("Key", newPString("value"))
	d.put("D", d2)

	a := newPArray()
	a.add(newPString("array"))
	d.put("A", a)

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
		t.Errorf("pDict: toBytes() = %v, want %v", d.toBytes(), out)
	}
}

var pDictTypeTests = []string{
	"hello",
	"there",
}

func TestPDictType(t *testing.T) {
	for _, dt := range pDictTypeTests {
		d1 := newPDictType(dt)
		d2 := newPDict()
		d2.put("Type", newPName(dt))
		if bytes.Compare(d1.toBytes(), d2.toBytes()) != 0 {
			t.Errorf("pDict newPDictType: toBytes() = %v, want %v",
				d1.toBytes(), d2.toBytes())
		}
	}
}

type pStreamTest struct {
	in, out []byte
}

// TODO check PDF Reference for empty streams and add a test for that
// TODO test append
var pStreamTests = []pStreamTest{
	pStreamTest{
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

func TestPStream(t *testing.T) {
	for _, st := range pStreamTests {
		s := newPStream(st.in)
		if bytes.Compare(s.toBytes(), st.out) != 0 {
			t.Errorf("pStream: toBytes() = %v, want %v",
				s.toBytes(), st.out)
		}
	}
}

func TestPNull(t *testing.T) {
	in := newPNull()
	out := []byte("null")
	if bytes.Compare(in.toBytes(), out) != 0 {
		t.Errorf("pNull: toBytes() = %v, want %v", in.toBytes(), out)
	}
}

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
		newPString("ssss"),
		1,
		100,
		[]byte("1 0 R"),
		[]byte("1 0 obj\n(ssss)\nendobj\n"),
		[]byte("0000000100 00000 n\r\n")},
	indirectTest{
		newPNumber(12.34),
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
