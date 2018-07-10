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
	"crypto/sha512"
	"encoding/binary"
	"errors"

	"github.com/AidosKuneen/aklib"

	sha256 "github.com/AidosKuneen/sha256-simd"
)

const prefixPrivString = "AKPRI"

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
func HDSeed58(conf *aklib.Config, seed, pwd []byte) string {
	out := make([]byte, len(seed)+4)
	copy(out, seed)
	hash := sha256.Sum256(seed)
	hash = sha256.Sum256(hash[:])

	copy(out[len(seed):], hash[0:4])
	eseed := enc(out, pwd)

	pref := conf.PrefixPriv
	s := make([]byte, len(eseed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], eseed)
	return prefixPrivString + Encode58(s)
}

//HDFrom58 returns seed bytes from base58-encoded seed and its password.
func HDFrom58(seed58 string, pwd []byte, cfg *aklib.Config) ([]byte, error) {
	if seed58[:len(prefixPrivString)] != prefixPrivString {
		return nil, errors.New("invalid prefix string in seed")
	}
	eseed, err := Decode58(seed58[len(prefixPrivString):])
	if err != nil {
		return nil, err
	}
	pref := cfg.PrefixPriv
	if !bytes.Equal(eseed[:len(pref)], pref) {
		return nil, errors.New("invalid prefix bytes")
	}
	eseed = eseed[len(cfg.PrefixPriv):]
	seed := enc(eseed, pwd)
	encoded := seed[:len(seed)-4]
	cksum := seed[len(seed)-4:]

	//Perform SHA-256 twice
	hash := sha256.Sum256(encoded)
	hash = sha256.Sum256(hash[:])
	if !bytes.Equal(hash[:4], cksum) {
		return nil, errors.New("invalid password")
	}
	return encoded, nil
}
