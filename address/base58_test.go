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
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestBase58(t *testing.T) {
	from := "00010966776006953D5567439E5E39F86A0D273BEE"
	to := "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM"
	dat, err := hex.DecodeString(from)
	if err != nil {
		t.Error(err)
	}
	edat := encode58(dat)
	if edat != to {
		t.Error("incorrect encoding:wrong result=" + edat)
	}

	ddat, err := decode58(to)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(dat, ddat) {
		t.Error("incorrect decoding:wrong result=" + hex.EncodeToString(ddat))
	}
}

func BenchmarkBase58A(b *testing.B) {
	from := make([]byte, 35+4)
	to := make([]byte, 35+4)
	// from := make([]byte, 35)
	// to := make([]byte, 35)
	for i := range to {
		to[i] = 0xff
	}
	for i := 0; i <= 0xff; i++ {
		from[0] = byte(i)
		to[0] = byte(i)
		for k := 0; k <= 0xff; k++ {
			from[1] = byte(k)
			to[1] = byte(k)
			for j := 0; j <= 0xff; j++ {
				from[2] = byte(j)
				to[2] = byte(j)
				ef := encode58(from)
				et := encode58(to)
				if ef[0] == et[0] && ef[1] == et[1] && ef[2] == et[2] {
					if ef[0] == 'S' || ef[0] == 'V' {
						if ef[1] == 'T' || ef[1] == 'M' {
							if ef[2] == '5' || ef[2] == '8' || ef[2] == 'A' {
								fmt.Printf("%x %x %x:%s %s\n", i, k, j, ef, et)
							}
						}
					}
				}
			}
		}
	}
}
