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
	"context"
	"errors"

	"github.com/AidosKuneen/aklib/rand"
	"github.com/AidosKuneen/cuckoo"
)

//PoW does PoW.
func (tx *Transaction) PoW() error {
	return tx.PoWContext(context.Background())
}

//ErrCanceled means the PoW was canceled by context.
var ErrCanceled = errors.New("PoW was canceled")

//PoWContext does PoW with context..
func (tx *Transaction) PoWContext(ctx context.Context) error {
	cu := cuckoo.NewCuckoo()
	start := rand.R.Uint32()
	for tx.Gnonce = start; tx.Gnonce != start-1; tx.Gnonce++ {
		hs := tx.hashForPoW()
		nonces, found := cu.PoW(hs)
		if found {
			tx.Nonce = nonces
			if isValidHash(tx.Hash(), tx.Easiness) {
				break
			}
		}
		select {
		case <-ctx.Done():
			return ErrCanceled
		default:
		}
	}
	if tx.Gnonce == start-1 {
		return errors.New("failed to PoW")
	}
	return nil
}
