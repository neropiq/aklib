// Copyright (c) 2017 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tx

import (
	"bytes"
	"log"

	"github.com/vmihailenco/msgpack"
)

//Pack returns tx in msgpack format.
func (tx *Transaction) Pack() []byte {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).StructAsArray(true)
	if err := enc.Encode(tx); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

//Pack returns tx body in msgpack format.
func (body *Body) Pack() []byte {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).StructAsArray(true)
	if err := enc.Encode(body); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

//Pack returns tx bSignaturesody in msgpack format.
func (sig *Signatures) Pack() []byte {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf).StructAsArray(true)
	if err := enc.Encode(sig); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

//Unpack returns tx from msgpack bin data.
func (tx *Transaction) Unpack(dat []byte) error {
	buf := bytes.NewBuffer(dat)
	dec := msgpack.NewDecoder(buf)
	return dec.Decode(tx)
}

//Unpack returns tx from msgpack bin data.
func (body *Body) Unpack(dat []byte) error {
	buf := bytes.NewBuffer(dat)
	dec := msgpack.NewDecoder(buf)
	return dec.Decode(body)
}

//Unpack returns tx from msgpack bin data.
func (sig *Signatures) Unpack(dat []byte) error {
	buf := bytes.NewBuffer(dat)
	dec := msgpack.NewDecoder(buf)
	return dec.Decode(sig)
}
