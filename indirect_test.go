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
	obj     interface{}
	num     int
	off     int
	out     []byte
	outBody []byte
	outRef  []byte
}

var indirectTests = []indirectTest{
	{"ssss", 1, 100, []byte("1 0 R"), []byte("1 0 obj\n(ssss)\nendobj\n"),
		[]byte("0000000100 00000 n\r\n")},
	{12.34, 3, 45, []byte("3 0 R"), []byte("3 0 obj\n12.34\nendobj\n"),
		[]byte("0000000045 00000 n\r\n")},
}

func TestIndirect(t *testing.T) {
	for _, dt := range indirectTests {
		i := newIndirect(dt.obj)
		i.setNum(dt.num)
		i.setOffset(dt.off)
		if bytes.Compare(i.output(), dt.out) != 0 {
			t.Errorf("indirect: toBytes() = %v, want %v",
				i.output(), dt.out)
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
		i := newIndirect(nil)
		i.set(dt.obj)
		i.setNum(dt.num)
		i.setOffset(dt.off)
		if bytes.Compare(i.output(), dt.out) != 0 {
			t.Errorf("indirect: toBytes() = %v, want %v",
				i.output(), dt.out)
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
