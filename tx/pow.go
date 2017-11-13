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
	"time"

	"github.com/AidosKuneen/cuckoo"
	sha256 "github.com/AidosKuneen/sha256-simd"
)

const maxcnt = 10

//PoW does PoW.
func (tx *Transaction) PoW() error {
	bs := tx.bytesForPoW()
	hs := sha256.Sum256(bs)
	bs2 := make([]byte, 0, len(bs))
	var nonces *[cuckoo.ProofSize]uint32
	found := false
	for cnt := 0; !found && cnt < maxcnt; cnt++ {
		nonces, found = cuckoo.PoW(hs[:], func(nonces *[cuckoo.ProofSize]uint32) bool {
			bs2 = append(bs2, bs[:4+1]...)
			for _, n := range nonces {
				bs2 = append(bs2, int2Varint(uint64(n))...)
			}
			bs2 = append(bs2, bs[4+1+cuckoo.ProofSize:]...)
			h := sha256.Sum256(bs2)
			bs2 = bs2[:0]
			return isValidHash(h[:], tx.Difficulty)
		})
		if !found {
			tx.Time = uint32(time.Now().Unix())
			bs = tx.bytesForPoW()
		}
	}
	if !found {
		return errors.New("failed to PoW")
	}
	copy(tx.Nonce, nonces[:])
	return nil
}
