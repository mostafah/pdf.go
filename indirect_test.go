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
	num    int
	off    int
	out    []byte
	outRef []byte
}

func TestIndirect(t *testing.T) {
	tests := []indirectTest{
		{1, 100, []byte("1 0 R"), []byte("0000000100 00000 n\r\n")},
		{3, 45, []byte("3 0 R"), []byte("0000000045 00000 n\r\n")},
	}

	for _, test := range tests {
		i := &indirect{num: test.num, off: test.off}
		if bytes.Compare(i.output(), test.out) != 0 {
			t.Errorf("indirect output: got\n\t%v\nexpected\n\t%v",
				i.output(), test.out)
		}
		if bytes.Compare(i.ref(), test.outRef) != 0 {
			t.Errorf("indirect ref: got\n\t%v\nexpected\n\t%v",
				i.ref(), test.outRef)
		}
	}
}
