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
	"bufio"
	"fmt"
	"os"
)

// TODO Use bytes.Buffer
// TODO Instead of making everything in memory and writing at the end, get the
// filename or writer at the beginning and write as we go. This reduces memory
// usage, specially for cases that include a lot of images.

// type Document holds all the objects of a PDF document.
type Document struct {
	objects    []indirect
	w          *bufio.Writer
	offset     uint64
	xrefOffset uint64
}

// New initializes a new Document objects and returns a pointer to it. The
// returned Document is ready for adding PDF objects like page, text, etc.
// and finally saving by calling Save.
func New() *Document {
	d := new(Document)
	d.objects = make([]indirect, 0, 10)
	return d
}

// TODO add a WriteTo method.

// Save writes document d into a PDF file.
func (d *Document) Save(fname string) (n uint64, err os.Error) {
	w, err := os.Create(fname)
	if err != nil {
		return 0, err
	}
	defer w.Close()
	d.w = bufio.NewWriter(w)
	defer d.w.Flush()
	return d.write()
}

// write saves the PDF document d to d.w.
func (d *Document) write() (n uint64, err os.Error) {
	if d.w == nil {
		return 0, error("writer is nil; cannot write the document")
	}

	d.offset = 0

	// TODO see if it error-handling can be done in another way
	if err = d.writeHeader(); err != nil {
		n = d.offset
		return
	}

	if err = d.writeBody(); err != nil {
		n = d.offset
		return
	}

	if err = d.writeRefs(); err != nil {
		n = d.offset
		return
	}

	if err = d.writeTrailer(); err != nil {
		n = d.offset
		return
	}

	n = d.offset
	return
}

// writeHeader writes the PDF header to d.w.
func (d *Document) writeHeader() (err os.Error) {
	b := []byte("%PDF-1.7\n\n")
	// TODO add comment with binary bytes (over 127)
	n, err := d.w.Write(b)
	d.offset += uint64(n)
	return
}

// writeBody prints PDF objects to d.w.
func (d *Document) writeBody() (err os.Error) {
	// Writing the objects to body and saving their offsets at the same time.
	for _, o := range d.objects {
		o.setOffset(d.offset)
		n, err := d.w.Write(o.body())
		d.offset += uint64(n)
		if err != nil {
			return
		}
	}
	return nil
}

// writeRefs prints the cross-reference table for the objects.
func (d *Document) writeRefs() (err os.Error) {
	d.xrefOffset = d.offset

	// Print the beginning 'xref' and number of objects
	n, err := fmt.Fprintf(d.w, "xref\n%d %d\n", 0, len(d.objects)+1)
	d.offset += uint64(n)
	if err != nil {
		return
	}

	// Print the first line in xref
	n, err = d.w.Write([]byte("0000000000 65535 f\r\n"))
	d.offset += uint64(n)
	if err != nil {
		return
	}

	// write references of the objects
	for _, object := range d.objects {
		n, err := d.w.Write(object.ref())
		d.offset += uint64(n)
		if err != nil {
			return
		}
	}
	return
}

// writeTrailer finishes of the PDF document.
func (d *Document) writeTrailer() (err os.Error) {
	n, err := d.w.Write([]byte("trailer\n"))
	d.offset += uint64(n)
	if err != nil {
		return
	}

	dic := newDict()
	dic.add(newName("Size"), newNumberInt(len(d.objects)+1))
	dic.add(newName("Root"), newStr("1 0 R"))
	b := dic.toBytes()
	n, err = d.w.Write(b)
	d.offset += uint64(n)
	if err != nil {
		return
	}

	n, err = d.w.Write([]byte("%%EOF\n"))
	d.offset += uint64(n)
	if err != nil {
		return
	}
	return
}

// error is a convenient function for generating errors in the this package.
func error(s string) os.Error {
	return os.NewError("PDF error:" + s)
}
