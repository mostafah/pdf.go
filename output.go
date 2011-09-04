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

// This file deals with representing data in PDF file. It's core functionality
// is about the output function which returns a []byte representation of variables.
// This []byte output is ready to be put in the PDF file.

import (
	"bytes"
	"reflect"
	"strconv"
)

// Names in PDF have a different representation than normal strings. Casting strings
// to this type makes them to appear as names in the output.
type name string

type outputter interface {
	output() []byte
}

// output gives out the PDF representation of v.
func output(v interface{}) []byte {
	// check for nil
	if v == nil {
		return []byte("null")
	}

	// check simple types with simple type assertion
	switch t := v.(type) {
	case outputter:
		return t.output()
	case bool:
		if t {
			return []byte("true")
		} else {
			return []byte("false")
		}
	case int:
		return []byte(strconv.Itoa(t))
	case float32:
		return []byte(strconv.Ftoa32(t, 'f', -1))
	case float64:
		// TODO 2.3 prints 2.299999952316284. Is it OK with PDF?
		return []byte(strconv.Ftoa64(t, 'f', -1))
	case string:
		// TODO non-ASCII characters?
		// TODO escapes, \n, \t, etc. (p. 54)
		// TODO break long lines (p. 54)
		// TODO what about hexadecimal strings? (p. 56)
		return []byte("(" + t + ")")
	case name:
		// TODO escape non-regular characters using # (p. 57)
		// TODO check length limit (p. 57)
		return []byte("/" + string(t))
	case []byte:
		return outputStream(t)
	case *bytes.Buffer:
		return outputStream(t.Bytes())
	case reflect.Value:
		return output(t.Interface())
	}

	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Invalid:
		panic(error("unsupported type passed to output"))
	case reflect.Array, reflect.Slice:
		buf := bytes.NewBufferString("[ ")

		for i := 0; i < r.Len(); i++ {
			buf.Write(output(r.Index(i)))
			buf.WriteString(" ")
		}

		buf.WriteString("]")

		return buf.Bytes()
	case reflect.Map:
		buf := bytes.NewBufferString("<<\n")

		for _, k := range r.MapKeys() {
			if k.Kind() != reflect.String {
				panic(("key of map passed to output is not string"))
			}
			buf.Write(output(name(k.String())))
			buf.WriteString(" ")
			buf.Write(output(r.MapIndex(k)))
			buf.WriteString("\n")
		}

		buf.WriteString(">>")

		return buf.Bytes()
	}

	return []byte("null")
}

// outputStream returns the given buffer as PDF stream.
func outputStream(b []byte) []byte {
	// TODO add filters

	// PDF streams start with a dictionary, then the word "stream", then
	// the stream itself, and finally the world "endstream". The slice all
	// holds []byte version of each of these four parts.
	all := make([][]byte, 4)

	all[0] = output(map[string]int{"Length": len(b)})
	all[1] = []byte("stream")
	all[2] = b
	all[3] = []byte("endstream")

	return bytes.Join(all, []byte{'\n'})
}