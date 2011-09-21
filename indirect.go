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

import (
	"fmt"
)

type indirect struct {
	num int // object number, i.e. ID among objects of the document
	off int // offset in bytes in the document
}

// output returns an indirect representation of i.
func (i *indirect) output() []byte {
	return []byte(fmt.Sprintf("%d 0 R", i.num))
}

// ref returns a refrence representation of i ready for the 'xref' section of a
// PDf file.
func (i *indirect) ref() []byte {
	return []byte(fmt.Sprintf("%010d 00000 n\r\n", i.off))
}
