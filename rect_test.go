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
	"rand"
	"testing"
)

func TestRect(t *testing.T) {
	const n = 10
	// testing with float64
	for i := 0; i < n; i++ {
		llx, lly := rand.NormFloat64(), rand.NormFloat64()
		urx, ury := rand.NormFloat64(), rand.NormFloat64()
		r := newRect(llx, lly, urx, ury)
		a := []float64{llx, lly, urx, ury}
		if bytes.Compare(r.output(), output(a)) != 0 {
			t.Errorf("newRect(%f, %f, %f, %f) doesn't work",
				llx, lly, urx, ury)
		}
	}
	// testing with int
	for i := 0; i < n; i++ {
		// random function that returns both positive and negative
		rnd := func() int {
			sign := 1
			if rand.Intn(2) == 0 {
				sign = -1
			}
			return sign * rand.Int()
		}

		llx, lly, urx, ury := rnd(), rnd(), rnd(), rnd()
		r := newRectInt(llx, lly, urx, ury)
		a := []int{llx, lly, urx, ury}
		if bytes.Compare(r.output(), output(a)) != 0 {
			t.Errorf("newRect(%d, %d, %d, %d) doesn't work",
				llx, lly, urx, ury)
		}
	}
}
