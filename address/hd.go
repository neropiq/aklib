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
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"

	"github.com/AidosKuneen/aklib"
)

const prefixPrivString = "AKPRI"
const prefixNkeyString = "AKNKE"

//HDseed returns HD seed.
func HDseed(masterkey []byte, indices ...uint32) []byte {
	seed := sha256.Sum256(masterkey)
	mac := hmac.New(sha512.New, masterkey)
	if _, err := mac.Write(seed[:]); err != nil {
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

//HDSeed58 returns base58-encoded encrypted seed..
func HDSeed58(conf *aklib.Config, seed, pwd []byte, isNode bool) string {
	out := make([]byte, len(seed)+4)
	copy(out, seed)
	hash := sha256.Sum256(seed)
	hash = sha256.Sum256(hash[:])

	copy(out[len(seed):], hash[0:4])
	eseed := enc(out, pwd)

	pref := prefixPriv(conf, isNode)
	s := make([]byte, len(eseed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], eseed)
	if !isNode {
		return prefixPrivString + Encode58(s)
	}
	return prefixNkeyString + Encode58(s)
}

//HDFrom58 returns seed bytes from base58-encoded seed and its password.
func HDFrom58(cfg *aklib.Config, seed58 string, pwd []byte) ([]byte, bool, error) {
	if len(seed58) <= len(prefixPrivString) {
		return nil, false, errors.New("invalid length")
	}
	prefixStr := seed58[:len(prefixPrivString)]
	eseed, err := Decode58(seed58[len(prefixPrivString):])
	if err != nil {
		return nil, false, err
	}
	prefix := eseed[:len(cfg.PrefixNode)]
	eseed = eseed[len(cfg.PrefixPriv):]
	seed := enc(eseed, pwd)
	encoded := seed[:len(seed)-4]
	cksum := seed[len(seed)-4:]

	//Perform SHA-256 twice
	hash := sha256.Sum256(encoded)
	hash = sha256.Sum256(hash[:])
	if !bytes.Equal(hash[:4], cksum) {
		return nil, false, errors.New("invalid password")
	}
	if prefixStr == prefixNkeyString && bytes.Equal(prefix, cfg.PrefixNkey) {
		return encoded, true, nil
	}
	if prefixStr == prefixPrivString && bytes.Equal(prefix, cfg.PrefixPriv) {
		return encoded, false, nil
	}
	return nil, false, errors.New("invalid prefix")
}
