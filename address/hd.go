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

package address

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
)

//HDseed returns HD seed.
func HDseed(seed, masterkey []byte, indices ...uint32) []byte {
	mac := hmac.New(sha512.New, masterkey)
	if _, err := mac.Write(seed); err != nil {
		panic(err)
	}
	key := mac.Sum(nil)
	private := key[:32]
	chain := key[32:]
	data := make([]byte, 32+4)
	for _, index := range indices {
		mac = hmac.New(sha512.New, chain)
		copy(data, private)
		binary.BigEndian.PutUint32(data[32:], index)
		if _, err := mac.Write(data); err != nil {
			panic(err)
		}
		key := mac.Sum(nil)
		private = key[:32]
		chain = key[32:]
	}
	return private
}
