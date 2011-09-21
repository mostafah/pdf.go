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
	"bytes"
	"io"
	"fmt"
	"log"
	"runtime"
	"os"
)

// Document holds all the objects of a PDF document.
type Document struct {
	objs []*indirect // All the PDF indirect objects of this document
	w    io.Writer
	off  int // Number of bytes already written to w
	xOff int // Offset of corss reference table

	// The following *indirect variables are pointers to elements of objs.
	cat   *indirect     // PDF catalog
	ptree *indirect     // Page tree
	pg    *page         // Current page
	pgs   []*indirect   // List of pages
	con   *bytes.Buffer // Current content stream.
}

// New initializes a new PDF document, ready to be filled by new pages, graphics,
// text, etc.
func New(w io.Writer) (d *Document, err os.Error) {
	defer dontPanic(&err)

	if w == nil {
		panic("pdf.New function was called with a nil parameter.")
	}

	// Initiate the docuemnt.
	d = new(Document)
	d.w = w
	d.objs = make([]*indirect, 0, 10)
	d.pgs = make([]*indirect, 0, 1)
	d.cat = d.reserveIndirect()   // to be later updated by saveCatalog
	d.ptree = d.reserveIndirect() // to be later updated by updatePageTree
	d.off = 0

	// Write header of the file.
	d.writeHeader()

	return d, nil
}

// Close finalizes the document by writing the rest of the PDF file to the output.
func (d *Document) Close() (err os.Error) {
	defer dontPanic(&err)

	// Save the pages and catalog.
	d.updatePageTree()
	d.saveCatalog()

	// Write the document to d.w.
	d.writeRefs()
	d.writeTrailer()
	return nil
}

// NewPage appends a new empty page to the document with the given size.
func (d *Document) NewPage(w, h int) (err os.Error) {
	defer dontPanic(&err)

	d.savePage() // Save the current one before starting anew.
	d.pg = newPage(w, h, d.ptree)
	return nil
}

// savePage writes the current page (d.pg) to the output.
func (d *Document) savePage() {
	if d.pg == nil {
		return
	}
	// Save the current content stream and add it to the page.
	d.pg.addContent(d.indirect(d.con))
	// Current content stream was written to the output, so we don't need it
	// anymore.
	d.con = nil

	// Add the page to the list of pages.
	d.pgs = append(d.pgs, d.indirect(d.pg))
}

// savePageTree makes page tree dictionary.
func (d *Document) updatePageTree() {
	d.savePage() // Save the current page first.

	tree := map[string]interface{}{
		"Type":  "Pages",
		"Count": len(d.pgs),
		"Kids":  d.pgs,
	}
	d.outputIndirect(d.ptree, tree)
}

// saveCatalog saves catalog!
func (d *Document) saveCatalog() {
	if d.ptree == nil {
		return
	}
	cat := map[string]interface{}{
		"Type":  "Catalog",
		"Pages": d.ptree,
	}
	d.outputIndirect(d.cat, cat)
}

// addc writes string to the current content stream. Functions that work
// with content, like Line and Stroke, use this to add content.
func (d *Document) addc(s string) {
	if d.con == nil {
		d.con = bytes.NewBuffer([]byte{})
	}
	d.con.Write([]byte(s + "\n"))
}

// writeHeader writes the PDF header to the output.
func (d *Document) writeHeader() {
	// Four non-ASCII charcters as a comment after header line are
	// recommended by PDF Reference for PDF files containing binary data.
	// This helps other applications treat the file as binary. "سلام" means
	// "hello" in Persian.
	b := []byte("%PDF-1.7\n%سلام\n")
	n, err := d.w.Write(b)
	d.off += n
	check(err)
}

// writeRefs prints the cross-reference table for the objects.
func (d *Document) writeRefs() {
	d.xOff = d.off

	// Print the beginning 'xref' and number of objects.
	n, err := fmt.Fprintf(d.w, "xref\n%d %d\n", 0, len(d.objs)+1)
	d.off += n
	check(err)

	// Print the first line in xref.
	n, err = d.w.Write([]byte("0000000000 65535 f\r\n"))
	d.off += n
	check(err)

	// Write references of the objects.
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

	// Dictionary referring to the catalog as root
	dic := map[string]interface{}{
		"Size": len(d.objs) + 1,
		"Root": d.cat,
	}
	n, err = d.w.Write(output(dic))
	d.off += n
	check(err)

	// Offset of 'xref' table
	n, err = d.w.Write([]byte(fmt.Sprintf("\nstartxref\n%d\n", d.xOff)))

	// Ending the document
	n, err = d.w.Write([]byte("%%EOF\n"))
	d.off += n
	check(err)
}

// indirect turns o into a PDF object, writes it to the output, and returns a PDF
// indirect reference to it.
func (d *Document) indirect(o interface{}) (i *indirect) {
	i = d.reserveIndirect()
	d.outputIndirect(i, o)
	return
}

// reverseIndirect makes and returns a new indirect object, but doesn't save it. The
// object itself can be outputted later by calling outputIndirect.
func (d *Document) reserveIndirect() (i *indirect) {
	i = &indirect{num: len(d.objs) + 1}
	d.objs = append(d.objs, i)
	return i
}

// outputIndirect writes o as a PDF indirect object to the output.
func (d *Document) outputIndirect(i *indirect, o interface{}) {
	i.off = d.off
	n, err := d.w.Write([]byte(fmt.Sprintf("%d 0 obj\n", i.num)))
	d.off += n
	check(err)
	n, err = d.w.Write(output(o))
	d.off += n
	check(err)
	n, err = d.w.Write([]byte("\nendobj\n"))
	d.off += n
	check(err)
}

// check panics if err is not nil.
func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

// dontPanic handles expectws panics and turns them into errors of type os.Error.
func dontPanic(err *os.Error) {
	if r := recover(); r != nil {
		switch e := r.(type) {
		case string:
			*err = os.NewError("pdf.go: " + e)
		case os.Error:
			*err = e
		default:
			panic(r)
		}
	}
}

// here prints the current line number and file name every time it's called.
func here() {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d", file, line)
}
