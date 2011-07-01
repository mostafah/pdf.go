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
	"fmt"
	"io"
	"os"
)

// TODO use bytes.Buffer

// type Document holds all the objects of a PDF document.
type Document struct {
	objs []object
	offsets []int64
	w io.WriteCloser
	offset int64
	xrefOffset int64
}

// New initializes a new Document objects and returns a pointer to it. The
// returned Document is ready for adding PDF objects like Page, Text, etc.
// and finally saved by calling WriteTo or Save.
func New() *Document {
	d := new(Document)
	d.objs = make([]object, 0, 10)
	return d
}

// WriteTo saves the PDF document d to w.
func (d *Document) WriteTo(w io.WriteCloser) (n int64, err os.Error) {
	d.w = w
	return d.write()
}

// Save saves the PDF document d into a file with the given file name.
func (d *Document) Save(fname string) (n int64, err os.Error) {
	d.w, err = os.Create(fname)
	if err != nil {
		return 0, err
	}
	defer d.w.Close()
	return d.write()
}

// write saves the PDF document d to d.w.
func (d *Document) write() (n int64, err os.Error) {
	if d.w == nil {
		return 0, error("w is nil; cannot write the document")
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

// writeHeader prints PDF header to d.w and updtes d.offset.
func (d *Document) writeHeader() (err os.Error) {
	n, err := d.w.Write([]byte("%PDF-1.7\n\n"))
	d.offset += int64(n)
	return
}

// writeBody prints PDF objects to d.w and updates d.offset.
func (d *Document) writeBody() (err os.Error) {
	d.offsets = make([]int64, len(d.objs))
	all := make([][]byte, 3)
	for i, o := range d.objs {
		d.offsets[i] = d.offset

		b := []byte(fmt.Sprintf("%d %d obj\n", i, 0))
		all[0] = b
		d.offset += int64(len(b))

		b = o.toBytes()
		all[1] = b
		d.offset += int64(len(b))

		b = []byte(fmt.Sprintf("endobj\n\n", i, 0))
		all[2] = b
		d.offset += int64(len(b))
	}
	_, err = d.w.Write(bytes.Join(all, nil))
	return
}

// writeRefs prints the cross-reference table for the objects.
func (d *Document) writeRefs() (err os.Error) {
	d.xrefOffset = d.offset

	// print number of objects
	n, err := fmt.Fprintf(d.w, "%d %d\n", 0, len(d.objs) + 1)
	d.offset += int64(n)
	if err != nil {
		return
	}

	// TODO find out what it is
	n, err = d.w.Write([]byte("0000000000 65535 f\n"))
	d.offset += int64(n)
	if err != nil {
		return
	}

	// write references of the objects
	for i, _ := range d.objs {
		b := []byte(fmt.Sprintf("%010d %05d n", d.offsets[i], 0))
		n, err := d.w.Write(b)
		d.offset += int64(n)
		if err != nil {
			return
		}
	}
	return
}

// writeTrailer finishes of the PDF document.
func (d *Document) writeTrailer() (err os.Error) {
	n, err := d.w.Write([]byte("trailer\n"))
	d.offset += int64(n)
	if err != nil {
		return
	}

	dic := newDict()
	dic.add(newName("Size"), newNumber(float64(len(d.objs) + 1)))
	dic.add(newName("Root"), newStr("1 0 R"))
	b := dic.toBytes()
	n, err = d.w.Write(b)
	d.offset += int64(n)
	if err != nil {
		return
	}

	n, err = d.w.Write([]byte("%%EOF\n"))
	d.offset += int64(n)
	if err != nil {
		return
	}
	return
}

// error is a convenient function for generating errors in the this package.
func error(s string) os.Error {
	return os.NewError("PDF error:" + s)
}