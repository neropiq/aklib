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
	"errors"

	"github.com/AidosKuneen/cuckoo"
	sha256 "github.com/AidosKuneen/sha256-simd"
)

const maxcnt = 5000

//PoW does PoW.
func (tx *Transaction) PoW() error {
	cu := cuckoo.NewCuckoo()
	for tx.Gnonce = 0; tx.Gnonce < maxcnt; tx.Gnonce++ {
		bs := tx.bytesForPoW()
		hs := sha256.Sum256(bs)
		nonces, found := cu.PoW(hs[:])
		if found {
			tx.Nonce = nonces
			txh := tx.Hash()
			if isValidHash(txh, tx.Difficulty) {
				break
			}
		}
	}
	if tx.Gnonce == maxcnt {
		return errors.New("failed to PoW")
	}
	return nil
}
