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
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
)

var debugConfig = &aklib.Config{
	Difficulty: 8,
	PrefixPriv: []byte{0xc0, 0x50}, //"VT" in base58
	PrefixAdrs: []byte{0xac, 0x8},  //"ST" in base58
}

var tx = &Transaction{
	Body: &Body{
		Type:    txType,
		Nonce:   make([]byte, 32),
		Time:    uint32(time.Now().Unix()),
		Message: []byte("This is a test for a transaction."),
		Inputs: []*Input{
			&Input{
				PreviousTX: make([]byte, 32),
				Index:      0,
			},
			&Input{
				PreviousTX: make([]byte, 32),
				Index:      3,
			},
		},
		Outputs: []*Output{
			&Output{
				Address: make([]byte, 32),
				Value:   111,
			},
			&Output{
				Address: make([]byte, 32),
				Value:   222,
			},
		},
		MultiSigs: []*MultiSig{
			&MultiSig{
				N: 3,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 2,
			},
			&MultiSig{
				N: 2,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 331,
			},
		},
		Previous: [][]byte{
			make([]byte, 32),
			make([]byte, 32),
		},
		Difficulty: debugConfig.Difficulty,
		LockTime:   0,
	},
	Signatures: [][]byte{
		make([]byte, 32),
		make([]byte, 32),
		make([]byte, 32),
	},
}

func TestTX(t *testing.T) {
	h := make([]byte, 32)
	h[30] = 0xff
	if !isValidHash(h, 8) {
		t.Error("isValidHash is incorrect")
	}
	h[31] = 0xff
	if isValidHash(h, 8) {
		t.Error("isValidHash is incorrect")
	}
	h[31] = 0x07
	if !isValidHash(h, 5) {
		t.Error("isValidHash is incorrect")
	}
	h[31] = 0x10
	if isValidHash(h, 5) {
		t.Error("isValidHash is incorrect")
	}
	h[31] = 0x00
	h[30] = 0x07
	if !isValidHash(h, 13) {
		t.Error("isValidHash is incorrect")
	}
	h[30] = 0x10
	if isValidHash(h, 13) {
		t.Error("isValidHash is incorrect")
	}
	if len(tx.NoExistHashes(m.GetTX, errNotFound)) != 4 {
		t.Error("invalid nonexistshashes")
	}
	if err := tx.Check(debugConfig); err == nil {
		t.Error("must be error")
	}
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(debugConfig); err != nil {
		t.Error(err)
	}

	tx.Time = uint32(time.Now().Add(2 * time.Second).Unix())
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(debugConfig); err == nil {
		t.Error("must be error")
	}

	time.Sleep(2 * time.Second)
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(debugConfig); err != nil {
		t.Error(err)
	}
	b := tx.Body.Pack()
	s := tx.Signatures.Pack()
	b2, err := UnpackBody(b)
	if err != nil {
		t.Error(err)
	}
	s2, err := UnpackSignature(s)
	if err != nil {
		t.Error(err)
	}
	t2 := &Transaction{
		Body:       b2,
		Signatures: s2,
	}
	h = tx.Hash()
	h2 := t2.Hash()
	if !bytes.Equal(h, h2) {
		t.Error("pack/unpack is incorrect")
	}
	tx.Nonce = make([]byte, 32)
}

type store map[[32]byte]*Body

var errNotFound = errors.New("not found")

func (s store) GetTX(hash []byte) (*Body, error) {
	var h [32]byte
	copy(h[:], hash)
	d, ok := s[h]
	if !ok {
		return nil, errNotFound
	}
	return d, nil
}

var m = store{
	[32]byte{0x1}: &Body{
		Outputs: []*Output{
			&Output{
				Address: make([]byte, 32),
				Value:   543,
			},
			&Output{
				Address: make([]byte, 32),
				Value:   0,
			},
		},
		MultiSigs: []*MultiSig{
			&MultiSig{
				N: 3,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 0,
			},
			&MultiSig{
				N: 2,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 123,
			},
		},
	},
	[32]byte{0x2}: &Body{
		Outputs: []*Output{
			&Output{
				Address: make([]byte, 32),
				Value:   0,
			},
		},
		MultiSigs: []*MultiSig{
			&MultiSig{
				N: 3,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 0,
			},
			&MultiSig{
				N: 2,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 0,
			},
		},
	},
}

func TestTX2(t *testing.T) {
	tx.Inputs[0].PreviousTX[0] = 0x1
	tx.Inputs[1].PreviousTX[0] = 0x1
	tx.Previous[0][0] = 0x1
	tx.Previous[1][0] = 0x2

	seed1 := address.GenerateSeed()
	seed2 := address.GenerateSeed()
	seed3 := address.GenerateSeed()
	seed4 := address.GenerateSeed()
	a1 := address.New(2, seed1, debugConfig)
	a2 := address.New(2, seed2, debugConfig)
	a3 := address.New(2, seed3, debugConfig)
	a4 := address.New(2, seed4, debugConfig)
	var d [32]byte
	d[0] = 0x1
	m[d].Outputs[0].Address = a1.PublicKey()
	m[d].MultiSigs[1].Addresses[0] = a2.PublicKey()
	m[d].MultiSigs[1].Addresses[1] = a3.PublicKey()
	m[d].MultiSigs[1].Addresses[2] = a4.PublicKey()

	dat := tx.BytesForSign()
	s1 := a1.Sign(dat)
	s3 := a3.Sign(dat)
	s4 := a4.Sign(dat)
	tx.Signatures[0] = s1
	tx.Signatures[1] = s3
	tx.Signatures[2] = s4
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.CheckAll(m.GetTX, address.Verify, debugConfig); err != nil {
		t.Error(err)
	}
	if len(tx.NoExistHashes(m.GetTX, errNotFound)) != 0 {
		t.Error("invalid nonexistshashes")
	}
	tx.Nonce = make([]byte, 32)
	tx.Inputs[0].PreviousTX[0] = 0
	tx.Inputs[1].PreviousTX[0] = 0
	tx.Previous[0][0] = 0
	tx.Previous[1][0] = 0
	tx.Signatures[0] = make([]byte, 32)
	tx.Signatures[1] = make([]byte, 32)
	tx.Signatures[2] = make([]byte, 32)
}

func TestTX3(t *testing.T) {
	tx.Inputs[0].PreviousTX[0] = 0x1
	tx.Inputs[1].PreviousTX[0] = 0x1
	tx.Previous[0][0] = 0x1
	tx.Previous[1][0] = 0x2

	seed1 := address.GenerateSeed()
	seed2 := address.GenerateSeed()
	seed3 := address.GenerateSeed()
	seed4 := address.GenerateSeed()
	a1 := address.New(2, seed1, debugConfig)
	a2 := address.New(2, seed2, debugConfig)
	a3 := address.New(2, seed3, debugConfig)
	a4 := address.New(2, seed4, debugConfig)
	var d [32]byte
	d[0] = 0x1
	m[d].Outputs[0].Address = a1.PublicKey()
	m[d].MultiSigs[1].Addresses[0] = a2.PublicKey()
	m[d].MultiSigs[1].Addresses[1] = a2.PublicKey()
	m[d].MultiSigs[1].Addresses[2] = a4.PublicKey()

	dat := tx.BytesForSign()
	s1 := a1.Sign(dat)
	s3 := a3.Sign(dat)
	s4 := a4.Sign(dat)
	tx.Signatures[0] = s1
	tx.Signatures[1] = s3
	tx.Signatures[2] = s4
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.CheckAll(m.GetTX, address.Verify, debugConfig); err == nil {
		t.Error("must be error")
	}
	tx.Nonce = make([]byte, 32)
	tx.Inputs[0].PreviousTX[0] = 0
	tx.Inputs[1].PreviousTX[0] = 0
	tx.Previous[0][0] = 0
	tx.Previous[1][0] = 0
	tx.Signatures[0] = make([]byte, 32)
	tx.Signatures[1] = make([]byte, 32)
	tx.Signatures[2] = make([]byte, 32)
}

func TestTX4(t *testing.T) {
	if err := tx.generalPoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(debugConfig); err != nil {
		t.Error(err)
	}
	tx.Nonce = make([]byte, 32)
	if err := tx.PoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(debugConfig); err != nil {
		t.Error(err)
	}
	tx.Nonce = make([]byte, 32)
}

func BenchmarkGeneralPoW(b *testing.B) {
	tx.Difficulty = 30
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for i := range tx.Nonce {
			tx.Nonce[i] = 0
		}
		b.StartTimer()
		if err := tx.generalPoW(); err != nil {
			b.Error(err)
		}
	}
	b.Log(tx.Hash())
	b.Log(tx.Nonce)
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	tx.Difficulty = 8
}

func BenchmarkPoW(b *testing.B) {
	tx.Difficulty = 30
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for i := range tx.Nonce {
			tx.Nonce[i] = 0
		}
		b.StartTimer()
		if err := tx.PoW(); err != nil {
			b.Error(err)
		}
	}
	b.Log(tx.Hash())
	b.Log(tx.Nonce)
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	tx.Difficulty = 8
}
