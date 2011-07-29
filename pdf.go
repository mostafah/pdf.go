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

// Package pdf aims to be a pretty low-level library for generating PDF files.
package pdf

import (
	"io"
	"fmt"
	"os"
)

// type Document holds all the objects of a PDF document.
type Document struct {
	objs []indirect // list of main PDF objects to be put in the 'body'
	w    io.Writer  // the output for writing the PDF file
	off  int        // keeps track of number of bytes already written to PDF file
	xOff int        // offset of corss reference table of PDF file in bytes

	// Variables of type *indirect that come below will be pointers to elements
	// of objs slice.
	cat   *indirect // PDF catalog
	ptree *indirect // page tree

	pg  *page       // current page
	pgs []*indirect // list of pages as pointers to elements of objs
}

// New initializes a new Document objects and returns a pointer to it. The
// returned Document is ready for adding PDF objects like page, text, etc.
// Save() method should be called when work with this Document is done.
func New(w io.Writer) *Document {
	// initiate the docuemnt
	d := new(Document)
	d.w = w
	d.objs = make([]indirect, 0, 10)
	d.pgs = make([]*indirect, 0, 1)
	// Add catalog and page tree as null first, so that they can be reffered to
	// by others. At the end the object in this indirect will be replaced by
	// pDict objects containing real catalog and page tree.
	d.cat = d.add(newPNull())
	d.ptree = d.add(newPNull())

	return d
}

// NewPage appends a new empty page to the document with the given size.
func (d *Document) NewPage(w, h int) {
	d.savePage() // save the current one before starting anew
	d.pg = newPage(w, h, d.ptree)
}

// savePage adds the current page to obj.
func (d *Document) savePage() {
	if d.pg == nil {
		return
	}
	i := d.add(d.pg.pObject())
	d.pgs = append(d.pgs, i)
}

// add makes o an indirect object and appends it to objects of d. It returns a
// pointer to the indirect object.
func (d *Document) add(o pObject) (i *indirect) {
	i = newIndirect(o)
	i.setNum(len(d.objs))
	d.objs = append(d.objs, *i)
	return i
}

// Save writes the PDF file into the writer of d.
func (d *Document) Save() {
	d.updatePageTree()
	d.saveCatalog()
	d.write()
}

// savePageTree makes up page tree dictionary.
func (d *Document) updatePageTree() {
	d.savePage() // save the last page first

	tree := newPDictType("Pages")
	tree.put("Count", newPNumberInt(len(d.pgs)))
	// kids is an array of indirect references to pages.
	kids := newPArray()
	for _, p := range d.pgs {
		kids.add(p)
	}
	tree.put("Kids", kids)
	d.ptree.set(tree)
}

func (d *Document) saveCatalog() {
	if d.ptree == nil {
		return
	}
	cat := newPDictType("Catalog")
	cat.put("Pages", d.ptree)
	d.cat.set(cat)
}

// write saves the PDF document d to d.w.
func (d *Document) write() (n int, err os.Error) {
	if d.w == nil {
		return 0, error("writer is nil; cannot write the document")
	}

	d.off = 0

	// TODO see if it error-handling can be done in another way
	if err = d.writeHeader(); err != nil {
		n = d.off
		return
	}

	if err = d.writeBody(); err != nil {
		n = d.off
		return
	}

	if err = d.writeRefs(); err != nil {
		n = d.off
		return
	}

	if err = d.writeTrailer(); err != nil {
		n = d.off
		return
	}

	n = d.off
	return
}

// writeHeader writes the PDF header to d.w.
func (d *Document) writeHeader() (err os.Error) {
	// Four non-ASCII charcters as a comment after header line are
	// recommended by PDF Reference for PDF files containing binary data.
	// This helps other applications treat the file as binary.
	b := []byte("%PDF-1.7\n%سلام\n")
	n, err := d.w.Write(b)
	d.off += n
	return
}

// writeBody prints PDF objects to d.w.
func (d *Document) writeBody() (err os.Error) {
	// Writing the objects to body and saving their offsets at the same time.
	for _, o := range d.objs {
		o.setOffset(d.off)
		n, err := d.w.Write(o.body())
		d.off += n
		if err != nil {
			return
		}
	}
	return nil
}

// writeRefs prints the cross-reference table for the objects.
func (d *Document) writeRefs() (err os.Error) {
	d.xOff = d.off

	// Print the beginning 'xref' and number of objects
	n, err := fmt.Fprintf(d.w, "xref\n%d %d\n", 0, len(d.objs)+1)
	d.off += n
	if err != nil {
		return
	}

	// Print the first line in xref
	n, err = d.w.Write([]byte("0000000000 65535 f\r\n"))
	d.off += n
	if err != nil {
		return
	}

	// write references of the objects
	for _, obj := range d.objs {
		n, err := d.w.Write(obj.ref())
		d.off += n
		if err != nil {
			return
		}
	}
	return
}

// writeTrailer finishes of the PDF document.
func (d *Document) writeTrailer() (err os.Error) {
	// 'trailer' title
	n, err := d.w.Write([]byte("trailer\n"))
	d.off += n
	if err != nil {
		return
	}

	// dictionary referring to the catalog as root
	dic := newPDict()
	dic.put("Size", newPNumberInt(len(d.objs)+1))
	//	dic.put("Root", d.objects[d.catalog])
	b := dic.toBytes()
	n, err = d.w.Write(b)
	d.off += n
	if err != nil {
		return
	}

	// ending the document
	n, err = d.w.Write([]byte("%%EOF\n"))
	d.off += n
	if err != nil {
		return
	}
	return
}

// error is a convenient function for generating errors in the this package.
func error(s string) os.Error {
	return os.NewError("PDF error:" + s)
}
