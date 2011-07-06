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

// This file contains code about the PDF objects. Objects in PDF include
//
//  • Boolean values
//  • Integer and real numbers
//  • Strings
//  • Names
//  • Arrays
//  • Dictionaries
//  • Streams
//  • The null object
//
// Each of these eight types are represented by a type implementing object
// interface.

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// object is the interface that each PDF object must implement by providing
// a toBytes method, used for saving the objects to the PDF file
type object interface {
	toBytes() []byte
}


// -----
// Type boolean represents boolean objects of PDF documents.
type boolean bool

// newBoolean creates a new boolean with default value false.
func newBoolean(v bool) *boolean {
	b := new(boolean)
	*b = boolean(v)
	return b
}

// set changes the value of b.
func (b *boolean) set(v bool) {
	*b = boolean(v)
}

// toBytes returns bytes of either "true" or "false", according to the value of
// b.
func (b *boolean) toBytes() []byte {
	if bool(*b) {
		return []byte("true")
	}
	return []byte("false")
}


// -----
// Type number represents integer and real numbers in PDF documents.
type number float64

// newNumber creates a new number with value v.
func newNumber(v float64) *number {
	n := new(number)
	*n = number(v)
	return n
}

// newNumberInt creates a new number with value v.
func newNumberInt(v int) *number {
	return newNumber(float64(v))
}

// set changes the value of n.
func (n *number) set(v float64) {
	*n = number(v)
}

// setInt changes the value of n with given int.
func (n *number) setInt(v int) {
	*n = number(v)
}

// toBytes returns a PDF-ready representation of n.
func (n *number) toBytes() []byte {
	return []byte(strconv.Ftoa64(float64(*n), 'f', -1))
}


// -----
// Type str represents string values in PDF documents.
type str string

// TODO escape special characters (p. 54)
// TODO break long lines (p. 54)
// TODO what about hexadecimal strings? (p. 56)

// newStr creates a new str with default value "".
func newStr(v string) *str {
	s := new(str)
	*s = str(v)
	return s
}

// set changes the value of s.
func (s *str) set(v string) {
	*s = str(v)
}

// toBytes returns a PDF-ready representation of s.
func (s *str) toBytes() []byte {
	// TODO non-ASCII characters?
	// TODO escapes, \n, \t, etc.
	return []byte("(" + string(*s) + ")")
}


// -----
// Type name represents names in PDF documents.
type name string

// TODO escape non-regular characters using # (p. 57)
// TODO check length limit (p. 57)

// newName creates a new name with default value "".
func newName(v string) *name {
	n := new(name)
	*n = name(v)
	return n
}

// set changes the value of n.
func (n *name) set(v string) {
	*n = name(v)
}

// toBytes returns a PDF-ready representation of n.
func (n *name) toBytes() []byte {
	return []byte("/" + string(*n))
}


// -----
// Type array represents array objects in PDF documents.
type array []object

// newArray creates a new empty array.
func newArray() *array {
	return new(array)
}

// add appends a new object at the end of the a.
func (a *array) add(o object) {
	*a = array(append([]object(*a), o))
}

// toBytes returns a PDF-ready representation of a.
func (a *array) toBytes() []byte {
	// Make a new slice to hold each part.
	all := make([][]byte, len([]object(*a))+2)

	// Fill the slice.
	all[0] = []byte{'['}
	for i, v := range []object(*a) {
		all[i+1] = v.toBytes()
	}
	all[len(all)-1] = []byte{']'}

	// Now join all the bytes with space as separator.
	return bytes.Join(all, []byte{' '})
}


// -----
// Type dict represents dictionary objects in PDF documents.
type dict []pair

// newDict creates a new empty dict.
func newDict() *dict {
	return new(dict)
}

// add makes a new key/value pair and appends it at the end of d.
func (d *dict) add(k *name, v object) {
	p := newPair(k, v)
	*d = dict(append([]pair(*d), *p))
}

// toBytes retunrs a PDF-ready representation of d.
func (d *dict) toBytes() []byte {
	// Make a new slice to hold each part.
	all := make([][]byte, len([]pair(*d))+2)

	// Fill the slice.
	all[0] = []byte{'<', '<'}
	for i, p := range []pair(*d) {
		all[i+1] = p.toBytes()
	}
	all[len(all)-1] = []byte{'>', '>'}

	// Now join all the bytes with space as separator.
	return bytes.Join(all, []byte{'\n'})
}

// Type pair holds key/value pairs for using in dict.
type pair struct {
	key   name
	value object
}

// newPair returns a new pair made with the given key/value as k and v.
func newPair(k *name, v object) *pair {
	return &pair{key: *k, value: v}
}
// toBytes returns a PDF-ready representation of p. By this method type pair
// implements interface object, but it's only used by type dict and is not one
// of PDF's eight object types.
func (p *pair) toBytes() []byte {
	all := [][]byte{p.key.toBytes(), p.value.toBytes()}
	return bytes.Join(all, []byte{' '})
}


// -----
// Type stream represents stream objects on PDF documents.
type stream bytes.Buffer

// newStream() returns a new stream filled with the given bytes.
func newStream(v []byte) *stream {
	b := bytes.NewBuffer(v)
	return (*stream)(b)
}

// append gets a new slice of bytes and adds that to the end of s.
func (s *stream) append(v []byte) (err os.Error) {
	_, err = (*bytes.Buffer)(s).Write(v)
	return
}

// toBytes returns a PDF-ready representation of s.
func (s *stream) toBytes() []byte {
	// PDF streams start with a dictionary, then the word "stream", then
	// the stream itself, and finally the world "endstream". The slice all
	// holds []byte version of each of these four parts.
	all := make([][]byte, 4)

	all[1] = []byte("stream")
	all[2] = (*bytes.Buffer)(s).Bytes()
	all[3] = []byte("endstream")

	// The dictionary part is added at the end because it should have the
	// length of the stream in it.
	d := newDict()
	d.add(newName("Length"), newNumber(float64(len(all[2]))))
	all[0] = d.toBytes()

	return bytes.Join(all, []byte{'\n'})
}


// -----
// Type null represents null objects in PDF documents.
type null bool

// newNull creates a new null.
func newNull() *null {
	return new(null)
}

// toBytes returns a PDF-ready representation of n.
func (n *null) toBytes() []byte {
	return []byte("null")
}


// -----
// Type indirect holds a PDF object and represents it as a PDF indirect object.
// In PDF terminology, it's indirect vernion of its object.
type indirect struct {
	obj    object
	num    uint32
	offset uint64
}

// newIndirect gets a PDF object and returns an indirect
func newIndirect(o object) *indirect {
	return &indirect{obj: o, num: 0, offset: 0}
}

// setNum assigns an object number to i. It should be called after i was added
// to the objects of the document, i.e. as soon as it's object number is found.
func (i *indirect) setNum(n uint32) {
	i.num = n
}

// setOffset gives the byte offset of i in document to it. It's necessary for
// calling ref later.
func (i *indirect) setOffset(o uint64) {
	i.offset = o
}

// toBytes returns an indirect representation of i.
func (i *indirect) toBytes() []byte {
	return []byte(fmt.Sprintf("%d 0 R", i.num))
}

// body returns a representation of i ready for the 'body' section of a PDF file.
func (i *indirect) body() []byte {
	head := fmt.Sprintf("%d 0 obj\n", i.num)
	buf := bytes.NewBufferString(head)
	// TODO set i.obj to nil to free memory after it's been written.
	buf.Write(i.obj.toBytes())
	buf.WriteString("\nendobj\n")
	return buf.Bytes()
}

// ref returns a refrence representation of i ready for the 'xref' section of a
// PDf file.
func (i *indirect) ref() []byte {
	return []byte(fmt.Sprintf("%010d 00000 n\r\n", i.offset))
}
