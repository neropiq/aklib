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

import "github.com/vmihailenco/msgpack"

//Pack returns tx in msgpack format.
func (tx *Transaction) Pack() []byte {
	bd, err := msgpack.Marshal(tx)
	if err != nil {
		panic(err)
	}
	return bd
}

//Pack returns tx body in msgpack format.
func (body *Body) Pack() []byte {
	bd, err := msgpack.Marshal(body)
	if err != nil {
		panic(err)
	}
	return bd
}

//Pack returns tx bSignaturesody in msgpack format.
func (sig *Signatures) Pack() []byte {
	bd, err := msgpack.Marshal(sig)
	if err != nil {
		panic(err)
	}
	return bd
}

//Unpack returns tx from msgpack bin data.
func (tx *Transaction) Unpack(dat []byte) error {
	return msgpack.Unmarshal(dat, tx)
}

//Unpack returns tx from msgpack bin data.
func (body *Body) Unpack(dat []byte) error {
	return msgpack.Unmarshal(dat, body)
}

//Unpack returns tx from msgpack bin data.
func (sig *Signatures) Unpack(dat []byte) error {
	return msgpack.Unmarshal(dat, sig)
}
