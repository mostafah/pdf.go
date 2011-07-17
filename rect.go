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

// rect holds a rectangle and can product a PDF array for it. It's a common
// data structure in PDF.
type rect struct {
	// x and y for lower-left and upper-right
	llx, lly, urx, ury float64
}

func newRect(llx, lly, urx, ury float64) *rect {
	return &rect{llx, lly, urx, ury}
}

func (r *rect) pObject() (a *pArray) {
	a = newPArray()
	a.add(newPNumber(r.llx))
	a.add(newPNumber(r.lly))
	a.add(newPNumber(r.urx))
	a.add(newPNumber(r.ury))
	return
}