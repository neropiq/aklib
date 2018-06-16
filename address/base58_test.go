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
	edat := Encode58(dat)
	if edat != to {
		t.Error("incorrect encoding:wrong result=" + edat)
	}

	ddat, err := Decode58(to)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(dat, ddat) {
		t.Error("incorrect decoding:wrong result=" + hex.EncodeToString(ddat))
	}
}

func BenchmarkBase58ASeed(b *testing.B) {
	from := make([]byte, 32+3+4)
	to := make([]byte, 32+3+4)
	for i := range to {
		to[i] = 0xff
	}
	es := []string{
		"VT1", "VT5", "VT8", "VTA",
		"VM1", "VM5", "VM8", "VMA",
	}
loop:
	for i := 0; i <= 0xff; i++ {
		from[0] = byte(i)
		to[0] = byte(i)
		for k := 0; k <= 0xff; k++ {
			from[1] = byte(k)
			to[1] = byte(k)
			for j := 0; j <= 0xff; j++ {
				from[2] = byte(j)
				to[2] = byte(j)
				ef := Encode58(from)
				et := Encode58(to)
				if ef[0] == et[0] && ef[1] == et[1] && ef[2] == et[2] {
					for m := len(es) - 1; m >= 0; m-- {
						e := es[m]
						if ef[0] == e[0] && ef[1] == e[1] && ef[2] == e[2] {
							fmt.Printf("[]byte{0x%02x, 0x%02x, 0x%02x}, //%s\n", i, k, j, e)
							es = append(es[:m], es[m+1:]...)
							if len(es) == 0 {
								break loop
							}
							break
						}
					}
				}
			}
		}
	}
}

func BenchmarkBase58AAddr(b *testing.B) {
	from := make([]byte, 33+3)
	to := make([]byte, 33+3)
	for i := range to {
		to[i] = 0xff
	}
	es := []string{
		"ST1", "ST5", "ST8", "STA",
		"YT1", "YT2", "YT3",
		"ET1", "ET2", "ET3",
		"SM1", "SM5", "SM8", "SMA",
		"YM1", "YM2", "YM3",
		"EM1", "EM2", "EM3",
	}

loop:
	for i := 0; i <= 0xff; i++ {
		from[0] = byte(i)
		to[0] = byte(i)
		for k := 0; k <= 0xff; k++ {
			from[1] = byte(k)
			to[1] = byte(k)
			for j := 0; j <= 0xff; j++ {
				from[2] = byte(j)
				to[2] = byte(j)
				ef := Encode58(from)
				et := Encode58(to)
				if ef[0] == et[0] && ef[1] == et[1] && ef[2] == et[2] {
					for m := len(es) - 1; m >= 0; m-- {
						e := es[m]
						if ef[0] == e[0] && ef[1] == e[1] && ef[2] == e[2] {
							fmt.Printf("[]byte{0x%02x, 0x%02x, 0x%02x}, //%s\n", i, k, j, e)
							es = append(es[:m], es[m+1:]...)
							if len(es) == 0 {
								break loop
							}
							break
						}
					}
				}
			}
		}
	}
}
