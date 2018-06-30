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

package tx

import (
	"errors"
	"fmt"

	"github.com/AidosKuneen/cuckoo"
)

const maxcnt = 5000

//PoW does PoW.
func (tx *Transaction) PoW() error {
	cu := cuckoo.NewCuckoo()
	for tx.Gnonce = 0; tx.Gnonce < maxcnt; tx.Gnonce++ {
		hs := tx.hashForPoW()
		nonces, found := cu.PoW(hs)
		if found {
			if len(nonces) != cuckoo.ProofSize {
				return fmt.Errorf("nonce from cuckoo PoW is wrong %d", len(tx.Nonce))
			}
			tx.Nonce = nonces
			txh := tx.Hash()
			if isValidHash(txh, tx.Easiness) {
				break
			}
		}
	}
	if tx.Gnonce == maxcnt {
		return errors.New("failed to PoW")
	}
	if len(tx.Nonce) != cuckoo.ProofSize {
		return fmt.Errorf("tx nonce is wrong %d", len(tx.Nonce))
	}
	return nil
}
