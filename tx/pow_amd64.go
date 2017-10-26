// +build amd64

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
	"runtime"
	"sync"

	pow64 "github.com/AidosKuneen/aklib/tx/internal/pow_amd64"
)

//PoW does PoW for AMD64(with SSE or AVX)..
func (tx *Transaction) PoW() error {
	if tx.Nonce == nil || len(tx.Nonce) != 32 {
		tx.Nonce = make([]byte, 32)
	}
	var wg sync.WaitGroup
	var nonce []byte
	var mutex sync.RWMutex
	for j := 0; j < runtime.NumCPU(); j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			tx2 := tx.Clone()
			add(tx2.Nonce, j)
			b := tx2.bytesForPoW()
			var n []byte
			mutex.RLock()
			n = nonce
			mutex.RUnlock()
			if pow64.PoW(b, tx2.Difficulty) {
				mutex.Lock()
				nonce = b[nonceLocation:]
				mutex.Unlock()
				return
			}
		}(j)
	}
	wg.Wait()
	if nonce == nil {
		return errors.New("cannot find PoW solution")
	}
	copy(tx.Nonce, nonce)
	return nil
}
