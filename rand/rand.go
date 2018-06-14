// Copyright (c) 2018 Aidos Developer

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

//This code is partially copied  from https://github.com/wadey/cryptorand (MIT license)

package rand

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type source struct{}

//R is a source of random numbers in math/rand.
var R = rand.New(&source{})

func (s source) Int63() int64 {
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		panic(err)
	}
	r := binary.BigEndian.Uint64(b)
	return int64(r & 0x7fffffffffffffff)
}

func (s source) Uint64() uint64 {
	b := make([]byte, 8)
	if _, err := crand.Read(b); err != nil {
		panic(err)
	}
	r := binary.BigEndian.Uint64(b)
	return r
}

func (source) Seed(int64) {
	panic("Seed() is not allowed")
}
