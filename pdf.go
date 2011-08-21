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
	"log"
	"runtime"
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

	con *pStream
}

// New initializes a new PDF document, ready to be filled by new pages, graphics,
// text, etc. The PDF file will be written to the given Writer, w. Save function
// should be called when document is ready.
func New(w io.Writer) *Document {
	// initiate the docuemnt
	d := new(Document)
	d.w = w
	d.objs = make([]indirect, 0, 10)
	d.pgs = make([]*indirect, 0, 1)

	// Add catalog and page tree as null first, so that they can be reffered to
	// by others. At the end the object in this indirect will be replaced by
	// pDict objects containing real catalog and page tree.
	d.cat = d.add(nil)   // to be later updated by saveCatalog function
	d.ptree = d.add(nil) // to be later updated by updatePageTree

	d.con = newPStream([]byte{})

	return d
}

// Save finalizes the document and writes the rest of the PDF file to the Writer,
// returning the total number of bytes written to it.
func (d *Document) Save() (n int, err os.Error) {
	// Error-handling
	defer func() {
		if r := recover(); r != nil {
			n = d.off
			switch e := r.(type) {
			case string:
				err = error(e)
			case os.Error:
				err = e
			default:
				panic(r)
			}
		}
	}()

	// Save the pages and catalog.
	d.updatePageTree()
	d.saveCatalog()

	// Write the document to d.w. Write functions increase d.off as they go.
	d.off = 0
	if d.w == nil {
		panic("writer is nil; cannot write the document")
	}
	d.writeHeader()
	d.writeBody()
	d.writeRefs()
	d.writeTrailer()
	return d.off, nil
}

// writeHeader writes the PDF header to d.w.
func (d *Document) writeHeader() {
	// Four non-ASCII charcters as a comment after header line are
	// recommended by PDF Reference for PDF files containing binary data.
	// This helps other applications treat the file as binary.
	b := []byte("%PDF-1.7\n%سلام\n")
	n, err := d.w.Write(b)
	d.off += n
	check(err)
}

// writeBody prints PDF objects to d.w.
func (d *Document) writeBody() {
	// Writing the objects to body and saving their offsets at the same time.
	for i, _ := range d.objs {
		d.objs[i].setOffset(d.off)
		n, err := d.w.Write(d.objs[i].body())
		d.off += n
		check(err)
	}
}

// writeRefs prints the cross-reference table for the objects.
func (d *Document) writeRefs() {
	d.xOff = d.off

	// Print the beginning 'xref' and number of objects
	n, err := fmt.Fprintf(d.w, "xref\n%d %d\n", 0, len(d.objs)+1)
	d.off += n
	check(err)

	// Print the first line in xref
	n, err = d.w.Write([]byte("0000000000 65535 f\r\n"))
	d.off += n
	check(err)

	// write references of the objects
	for _, o := range d.objs {
		n, err := d.w.Write(o.ref())
		d.off += n
		check(err)
	}
}

// writeTrailer finishes of the PDF document.
func (d *Document) writeTrailer() {
	// 'trailer' title
	n, err := d.w.Write([]byte("trailer\n"))
	d.off += n
	check(err)

	// dictionary referring to the catalog as root
	dic := newPDict()
	dic.put("Size", len(d.objs)+1)
	dic.put("Root", d.cat)
	b := dic.toBytes()
	n, err = d.w.Write(b)
	d.off += n
	check(err)

	// writing xref offset
	n, err = d.w.Write([]byte(fmt.Sprintf("startxref\n%d\n", d.xOff)))

	// ending the document
	n, err = d.w.Write([]byte("%%EOF\n"))
	d.off += n
	check(err)
}

// add makes o an indirect object and appends it to objects of d. It returns a
// pointer to the indirect object.
func (d *Document) add(o interface{}) (i *indirect) {
	i = newIndirect(obj(o))
	i.setNum(len(d.objs) + 1)
	d.objs = append(d.objs, *i)
	return &d.objs[len(d.objs)-1]
}

// savePageTree makes up page tree dictionary.
func (d *Document) updatePageTree() {
	d.savePage() // save the last page first

	tree := newPDictType("Pages")
	tree.put("Count", len(d.pgs))
	tree.put("Kids", d.pgs)
	d.ptree.set(tree)
}

// saveCatalog saves catalog!
func (d *Document) saveCatalog() {
	if d.ptree == nil {
		return
	}
	cat := newPDictType("Catalog")
	cat.put("Pages", d.ptree)
	d.cat.set(cat)
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
	d.pg.addContent(d.add(d.con))
	i := d.add(d.pg)
	d.pgs = append(d.pgs, i)
}

// addc writes string to the current content stream. Functions that work
// with content, like Line and Stroke, use this to add content.
func (d *Document) addc(s string) {
	d.con.append([]byte(s + "\n"))
}

// Line draws a single line from (x0, y0) to (x1, y1).
func (d *Document) Line(x0, y0, x1, y1 int) {
	d.addc(fmt.Sprint(x0, y0, " m\n", x1, y1, " l"))
}


// Stroke paints the current path with stroke.
func (d *Document) Stroke() {
	d.addc("S")
}

// error is a convenient function for generating errors in the this package.
func error(s string) os.Error {
	return os.NewError("PDF error:" + s)
}

// check panics if err is not nil
func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

// here prints the current line number and file name every time it's called.
func here() {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d", file, line)
}
