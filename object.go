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
	"reflect"
	"os"
	"strconv"
)

// pObject is the interface that each PDF object must implement by providing
// a toBytes method, used for saving the objects to the PDF file
type pObject interface {
	// toBytes returns a PDF-ready representation of pObject.
	toBytes() []byte
}

// -----
// boolean
type pBoolean bool

func newPBoolean(v bool) *pBoolean {
	b := new(pBoolean)
	*b = pBoolean(v)
	return b
}

func (b *pBoolean) set(v bool) {
	*b = pBoolean(v)
}

func (b *pBoolean) toBytes() []byte {
	if bool(*b) {
		return []byte("true")
	}
	return []byte("false")
}

// -----
// integer and real numbers
type pNumber float64

func newPNumber(v float64) *pNumber {
	n := new(pNumber)
	*n = pNumber(v)
	return n
}

func newPNumberInt(v int) *pNumber {
	return newPNumber(float64(v))
}

func (n *pNumber) set(v float64) {
	*n = pNumber(v)
}

func (n *pNumber) setInt(v int) {
	*n = pNumber(v)
}

func (n *pNumber) toBytes() []byte {
	return []byte(strconv.Ftoa64(float64(*n), 'f', -1))
}

// -----
// string
type pString string

func newPString(v string) *pString {
	s := new(pString)
	*s = pString(v)
	return s
}

func (s *pString) set(v string) {
	*s = pString(v)
}

func (s *pString) toBytes() []byte {
	// TODO non-ASCII characters?
	// TODO escapes, \n, \t, etc. (p. 54)
	// TODO break long lines (p. 54)
	// TODO what about hexadecimal strings? (p. 56)
	return []byte("(" + string(*s) + ")")
}

// -----
// name
type pName string

// TODO escape non-regular characters using # (p. 57)
// TODO check length limit (p. 57)

func newPName(v string) *pName {
	n := new(pName)
	*n = pName(v)
	return n
}

func (n *pName) set(v string) {
	*n = pName(v)
}

func (n *pName) toBytes() []byte {
	return []byte("/" + string(*n))
}

// -----
// array
type pArray []pObject

func newPArray() *pArray {
	return new(pArray)
}

func (a *pArray) add(o interface{}) {
	*a = pArray(append([]pObject(*a), pobj(o)))
}

func (a *pArray) toBytes() []byte {
	// Make a new slice to hold each part.
	all := make([][]byte, len([]pObject(*a))+2)

	// Fill the slice.
	all[0] = []byte{'['}
	for i, v := range []pObject(*a) {
		all[i+1] = v.toBytes()
	}
	all[len(all)-1] = []byte{']'}

	// Now join all the bytes with space as separator.
	return bytes.Join(all, []byte{' '})
}

// -----
// dictionary
type pDict []pair

func newPDict() *pDict {
	return new(pDict)
}

// newPDictType makes a new pDict like newPDict, except that it also adds a
// new pair to it the key "Type" and value typ.
func newPDictType(typ string) *pDict {
	d := new(pDict)
	d.put("Type", newPName(typ))
	return d
}

// put makes a new key/value pair and appends it at the end of d. If there is
// already a pair with the given key put updates that.
func (d *pDict) put(k string, v interface{}) {
	// search for a pair with this key
	for i, _ := range []pair(*d) {
		p := &([]pair(*d))[i]
		if p.key == k {
			// found; update the pair and return
			p.val = pobj(v)
			return
		}
	}
	// no pair found with the given key; make a new pair
	p := newPair(k, pobj(v))
	d.add(*p)
}

func (d *pDict) add(p pair) {
	*d = pDict(append([]pair(*d), p))
}

func (d *pDict) toBytes() []byte {
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

// Type pair holds key/value pairs for using in pDict. It implements interface
// pObject, but it's only used by type dict and is not one of PDF's eight
// object types.
type pair struct {
	key string
	val pObject
}

func newPair(k string, v pObject) *pair {
	return &pair{key: k, val: v}
}

func (p *pair) toBytes() []byte {
	all := [][]byte{newPName(p.key).toBytes(), p.val.toBytes()}
	return bytes.Join(all, []byte{' '})
}

// -----
// stream
type pStream bytes.Buffer

// TODO add filters

func newPStream(v []byte) *pStream {
	b := bytes.NewBuffer(v)
	return (*pStream)(b)
}

func (s *pStream) append(v []byte) (err os.Error) {
	_, err = (*bytes.Buffer)(s).Write(v)
	return
}

func (s *pStream) toBytes() []byte {
	// PDF streams start with a dictionary, then the word "stream", then
	// the stream itself, and finally the world "endstream". The slice all
	// holds []byte version of each of these four parts.
	all := make([][]byte, 4)

	b := (*bytes.Buffer)(s)
	d := newPDict()
	d.put("Length", newPNumberInt(b.Len()))

	all[0] = d.toBytes()
	all[1] = []byte("stream")
	all[2] = b.Bytes()
	all[3] = []byte("endstream")

	return bytes.Join(all, []byte{'\n'})
}

// -----
// null
type pNull byte

func newPNull() *pNull {
	return new(pNull)
}

func (n *pNull) toBytes() []byte {
	return []byte("null")
}

// -----
// Type indirect holds a pObject and represents it as a PDF indirect object.
type indirect struct {
	obj pObject
	num int // object number, i.e. ID among objects of the document
	off int // offset in bytes in the document
}

func newIndirect(o pObject) *indirect {
	return &indirect{obj: o, num: 0, off: 0}
}

// set updates pObject of i.
func (i *indirect) set(o pObject) {
	i.obj = o
}

// setNum assigns an object number to i. It should be called after i was added
// to the objects of the document, i.e. as soon as it's object number is found.
func (i *indirect) setNum(n int) {
	i.num = n
}

// setOffset gives the byte offset of i in document to it. It's necessary for
// calling ref later.
func (i *indirect) setOffset(o int) {
	i.off = o
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
	return []byte(fmt.Sprintf("%010d 00000 n\r\n", i.off))
}

// -----
// Functions that make working with pObject types easier and cleaner.

type pObjectable interface {
	pObject() pObject
}

// pobj makes a new pObject out of the given value.
func pobj(v interface{}) pObject {
	// check for nil
	if v == nil {
		return newPNull()
	}

	// check simple types with simple type assertion
	switch t := v.(type) {
	case bool:
		return newPBoolean(t)
	case int:
		return newPNumberInt(t)
	case float32:
		return newPNumber(float64(t))
	case float64:
		return newPNumber(t)
	case string:
		return newPString(t)
	case pObject:
		return t
	case []byte:
		return newPStream(t)
	case reflect.Value:
		return pobj(t.Interface())
	case pObjectable:
		return t.pObject()
		//	case bytes.Buffer, *bytes.Buffer:
		//		return newPStream(t.Bytes())
	}

	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Invalid:
		panic(error("unsupported type passed to pobj"))
	case reflect.Array, reflect.Slice:
		a := newPArray()
		for i := 0; i < r.Len(); i++ {
			a.add(pobj(r.Index(i)))
		}
		return a
	case reflect.Map:
		d := newPDict()
		for _, k := range r.MapKeys() {
			if k.Kind() != reflect.String {
				panic(("key of map passed to pobj is not string"))
			}
			d.put(k.String(), pobj(r.MapIndex(k)))
		}
		return d
	}

	return newPNull()
}
