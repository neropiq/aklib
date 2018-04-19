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
	"errors"
	"fmt"
	"math/big"

	sha256 "github.com/AidosKuneen/sha256-simd"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var decodeMap [255]byte

func init() {
	for i := 0; i < len(decodeMap); i++ {
		decodeMap[i] = 0xff
	}
	for i := 0; i < len(alphabet); i++ {
		decodeMap[alphabet[i]] = byte(i)
	}
}

//Encode58 encodes byte slice to base58.
func Encode58(encoded []byte) string {
	out := make([]byte, len(encoded)+4)
	copy(out, encoded)
	hash := sha256.Sum256(encoded)
	hash = sha256.Sum256(hash[:])
	copy(out[len(encoded):], hash[0:4])
	bigcksum := new(big.Int).SetBytes(out)
	cksum := encodeBig(bigcksum)
	buffer := make([]byte, 0, len(cksum))

	for i := 0; i < len(encoded) && encoded[i] == 0; i++ {
		buffer = append(buffer, '1')
	}
	buffer = append(buffer, cksum...)
	return string(buffer)
}

//Decode58 decodes base58 value to bytes.
func Decode58(value string) ([]byte, error) {
	if len(value) < 5 {
		return nil, errors.New("invalid input")
	}
	var pk big.Int
	if err := decodeToBig([]byte(value), &pk); err != nil {
		return nil, err
	}

	ecksum := pk.Bytes()
	encoded := ecksum[:len(ecksum)-4]
	cksum := ecksum[len(ecksum)-4:]

	buffer := make([]byte, 0, len(encoded))
	for i := 0; i < len(value) && value[i] == '1'; i++ {
		buffer = append(buffer, 0)
	}
	buffer = append(buffer, encoded...)

	//Perform SHA-256 twice
	hash := sha256.Sum256(buffer)
	hash = sha256.Sum256(hash[:])

	var err error
	if !bytes.Equal(hash[:4], cksum) {
		err = fmt.Errorf("%s checksum did not match to the embeded one:%s, should be :%s",
			value, hex.EncodeToString(hash[:4]), hex.EncodeToString(cksum))
	}
	return buffer, err
}

func decodeToBig(src []byte, n *big.Int) error {
	radix := big.NewInt(58)
	for i, s := range src {
		b := decodeMap[s]
		if b == 0xff {
			return fmt.Errorf("wrong data at %d", i)
		}
		n.Mul(n, radix).Add(n, big.NewInt(int64(b)))
	}
	return nil
}

func encodeBig(src *big.Int) []byte {
	n := new(big.Int)
	n.Set(src)
	radix := big.NewInt(58)
	zero := big.NewInt(0)
	var dst []byte

	for i := 0; n.Cmp(zero) > 0; i++ {
		mod := new(big.Int)
		n.DivMod(n, radix, mod)
		dst = append(dst, alphabet[mod.Int64()])
	}
	for i, j := 0, len(dst)-1; i < j; i, j = i+1, j-1 {
		dst[i], dst[j] = dst[j], dst[i]
	}
	return dst
}
