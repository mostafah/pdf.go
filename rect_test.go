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

type rectTest struct {
	llx, lly, urx, ury float64
	out []byte
}

var rectTests = []rectTest{
	rectTest{0, 10, 20, 30, []byte("[ 0 10 20 30 ]")},
	rectTest{-100, -200, 300, 400, []byte("[ -100 -200 300 400 ]")},
}

func TestRectNew(t *testing.T) {
	for _, rt := range rectTests {
		r := newRect(rt.llx, rt.lly, rt.urx, rt.ury)
		if bytes.Compare(r.pObject().toBytes(), rt.out) != 0 {
			t.Errorf("rect: after newRect(%v, %v, %v, %v), toBytes() ="+
				" %v, want %v", rt.llx, rt.lly, rt.urx, rt.ury,
				r.pObject().toBytes(), rt.out)
		}
	}
}