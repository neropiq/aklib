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

// 171 85 : SM143bKyeipzswzkX1mU4qSFP93tnpDwGRrJqQyfAKaH7YtZjwdQ   SM2znmZRgaFg75pLDvM64nkrA4AWWTvn9DVNwf8M4KEHW3SuAcx7
// 172 8 : ST11JAqKNFz4R5qmmHXZvZqa14U8LBGn1HFz75nndFsUXtEKPidN   ST2x3M4mQ7QjeDfMUC7BvXAAmyak3pyct4u4DKwUXFXUvNootRag
// 191 157 : VM1d5XEHT9nSg1KxMoQ1vzRA7mHHBzM9ZvL1n4h5XTpzfhfPMoxy   VM3ZphTjV1D7u99Y4hydvwjktgPtue3zShy5tJqmRTV14CE3qkDm
// 192 80 : VT1aL6jdAgwWD9Ayc5A7nipUjghWjMPzJmjh3jWCzQ8C633668b2   VT3X5Gy5CYNBSGzZJyjjng95Wbp8T16qBZNm9yettPnCUXWcASYf
func TestBase58A(t *testing.T) {
	from := make([]byte, 34)
	to := make([]byte, 34)
	for i := range to {
		to[i] = 0xff
	}
	var foundS, foundV bool
	for i := 0; i <= 0xff; i++ {
		from[0] = byte(i)
		to[0] = byte(i)
		for k := 0; k <= 0xff; k++ {
			from[1] = byte(k)
			to[1] = byte(k)
			ef := encode58(from)
			et := encode58(to)
			if ef[0] == et[0] && ef[0] == 'S' && ef[1] == et[1] && ef[1] == 'M' && !foundS {
				fmt.Println(i, k, ":", ef, " ", et)
			}
			if ef[0] == et[0] && ef[0] == 'S' && ef[1] == et[1] && ef[1] == 'T' && !foundS {
				fmt.Println(i, k, ":", ef, " ", et)
			}
			if ef[0] == et[0] && ef[0] == 'V' && ef[1] == et[1] && ef[1] == 'M' && !foundV {
				fmt.Println(i, k, ":", ef, " ", et)
			}
			if ef[0] == et[0] && ef[0] == 'V' && ef[1] == et[1] && ef[1] == 'T' && !foundV {
				fmt.Println(i, k, ":", ef, " ", et)
			}
		}
	}
}
