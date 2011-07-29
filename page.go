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

// This file deals with pages in PDF.

// type page holds a PDF page, its attributes and its content.
type page struct {
	box rect      // size of the page
	par *indirect // page tree for this page
}

func newPage(w, h int, par *indirect) *page {
	p := new(page)
	p.box = *newRectInt(0, 0, w, h)
	p.par = par
	return p
}

func (p *page) pObject() (d *pDict) {
	d = newPDictType("Page")
	d.put("Parent", p.par)
	d.put("MediaBox", p.box.pObject())
	// TODO Add Resources.
	// TODO Add Contets.
	return
}
