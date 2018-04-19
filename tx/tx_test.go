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
	"bytes"
	"encoding/hex"
	"errors"
	"runtime"
	"testing"
	"time"

	"crypto/sha256"

	"github.com/AidosKuneen/numcpu"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/cuckoo"
)

var tx = &Transaction{
	Body: &Body{
		Type:    TxNormal,
		Gnonce:  0,
		Nonce:   make([]uint32, cuckoo.ProofSize),
		Time:    time.Now(),
		Message: []byte("This is a test for a transaction."),
		Inputs: []*Input{
			&Input{
				PreviousTX: make([]byte, 32),
				Index:      0,
			},
			&Input{
				PreviousTX: make([]byte, 32),
				Index:      1,
			},
		},
		MultiSigIns: []*MultiSigIn{
			&MultiSigIn{
				PreviousTX: make([]byte, 32),
				Index:      1,
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
		MultiSigOuts: []*MultiSigOut{
			&MultiSigOut{
				N: 3,
				Addresses: [][]byte{
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
				},
				Value: 2,
			},
			&MultiSigOut{
				N: 2,
				Addresses: [][]byte{
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
					[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
				},
				Value: 331,
			},
		},
		Previous: [][]byte{
			make([]byte, 32),
			make([]byte, 32),
		},
		Easiness:     aklib.TestConfig.Easiness,
		LockTime:     time.Time{},
		HashType:     0,
		TicketInput:  nil,
		TicketOutput: nil,
		Scripts:      nil,
		Reserved:     nil,
	},
	Signatures: []*Signature{
		&Signature{
			PublicKey: make([]byte, 65),
			Sig:       make([]byte, 32),
		},
	},
}

func TestValidHashTX(t *testing.T) {
	h := make([]byte, 32)
	h[31] = 0x1f
	if !isValidHash(h, 0x1fffffff) {
		t.Error("isValidHash is incorrect")
	}
	if isValidHash(h, 0x1effffff) {
		t.Error("isValidHash is incorrect")
	}
	if !isValidHash(h, 0x20ffffff) {
		t.Error("isValidHash is incorrect")
	}
}
func TestPoW(t *testing.T) {
	seed1 := address.GenerateSeed()
	a1, err := address.New(address.Height10, seed1, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	s256 := sha256.Sum256(a1.PublicKey())
	tx.Signatures[0].PublicKey = s256[:]
	dat := tx.BytesForSign()
	s1 := a1.Sign(dat)
	tx.Signatures[0] = &Signature{
		PublicKey: a1.PublicKey(),
		Sig:       s1,
	}

	if err := tx.PoW(); err != nil {
		t.Error(err)
	}
	t.Log(hex.EncodeToString(tx.Hash()))
	if err := tx.Check(aklib.TestConfig); err != nil {
		t.Error(err)
	}
}
func TestTX(t *testing.T) {
	seed1 := address.GenerateSeed()
	a1, err := address.New(address.Height10, seed1, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	dat := tx.BytesForSign()
	s1 := a1.Sign(dat)
	tx.Signatures[0] = &Signature{
		PublicKey: a1.PublicKey(),
		Sig:       s1,
	}

	tx.Nonce = make([]uint32, cuckoo.ProofSize)
	if len(tx.NoExistHashes(m.GetTX)) != 5 {
		t.Error("invalid nonexistshashes")
	}
	if err := tx.Check(aklib.TestConfig); err == nil {
		t.Error("must be error")
	}
	if err := tx.PoW(); err != nil {
		t.Error(err)
	}
	if err := tx.Check(aklib.TestConfig); err != nil {
		t.Error(err)
	}
	tx.Inputs[1].Index = 0
	dat = tx.BytesForSign()
	s1 = a1.Sign(dat)
	tx.Signatures[0] = &Signature{
		PublicKey: a1.PublicKey(),
		Sig:       s1,
	}
	if err := tx.check(aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}
	tx.Inputs[1].Index = 1
	tx.Time = time.Now().Add(time.Hour)
	tx.Signatures[0].Sig = a1.Sign(tx.BytesForSign())
	if err := tx.check(aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}

	tx.Time = time.Now().Add(-1 * time.Minute)
	tx.Signatures[0].Sig = a1.Sign(tx.BytesForSign())

	if err := tx.check(aklib.TestConfig, false); err != nil {
		t.Error(err)
	}
	tx.Time = time.Now()

	b := tx.Body.Pack()
	s := tx.Signatures.Pack()
	var b2 Body
	if err := b2.Unpack(b); err != nil {
		t.Error(err)
	}
	var s2 Signatures
	if err := s2.Unpack(s); err != nil {
		t.Error(err)
	}
	t2 := &Transaction{
		Body:       &b2,
		Signatures: s2,
	}
	h := tx.Hash()
	h2 := t2.Hash()
	if !bytes.Equal(h, h2) {
		t.Error("pack/unpack is incorrect")
	}
	t2 = tx.Clone()
	h2 = t2.Hash()
	if !bytes.Equal(h, h2) {
		t.Error("clone is incorrect")
	}

	tx.LockTime = time.Now().Add(time.Hour)
	if err := tx.check(aklib.TestConfig, false); err == nil {
		t.Error("should be error")
	}
	tx.LockTime = time.Time{}

	tx.HashType = 0x1f
	if err := tx.check(aklib.TestConfig, false); err == nil {
		t.Error(err)
	}
	tx.HashType = 0x11
	hs1 := tx.BytesForSign()
	tx.Outputs[1].Address[31] = 0x1
	tx.Outputs[1].Value = 1357
	hs2 := tx.BytesForSign()
	if !bytes.Equal(hs1, hs2) {
		t.Error("invalid hashtype")
	}
	tx.Outputs[1].Address[31] = 0
	tx.Outputs[1].Value = 222
	tx.Signatures[0].Sig = a1.Sign(tx.BytesForSign())
	typ, err := tx.IsMinable(aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	if typ != RewardFee {
		t.Error("invalid reward type")
	}

	tx.HashType = 0
	tx.Signatures[0].Sig = a1.Sign(tx.BytesForSign())
	if _, err := tx.IsMinable(aklib.TestConfig); err == nil {
		t.Error("invalid isminable")
	}

	tx.Nonce = make([]uint32, cuckoo.ProofSize)
}

func TestTicket(t *testing.T) {
	ticket := &Transaction{
		Body: &Body{
			Type:         TxNormal,
			Time:         time.Now(),
			Easiness:     aklib.TestConfig.TicketEasiness,
			TicketOutput: make([]byte, 32),
			Previous: [][]byte{
				make([]byte, 32),
				make([]byte, 32),
			},
		},
	}
	if err := ticket.PoW(); err != nil {
		t.Error(err)
	}
	if err := ticket.Check(aklib.TestConfig); err != nil {
		t.Error(err)
	}
	if !isValidHash(ticket.Hash(), aklib.TestConfig.TicketEasiness) {
		t.Error("invlaid ticket hash")
	}
	ticket.Easiness = aklib.TestConfig.Easiness
	if err := ticket.PoW(); err != nil {
		t.Error(err)
	}
	if err := ticket.Check(aklib.TestConfig); err == nil {
		t.Error("should be error")
	}

}

type store map[[32]byte]*Body

func (s store) GetTX(hash []byte) (*Body, error) {
	var h [32]byte
	copy(h[:], hash)
	d, ok := s[h]
	if !ok {
		return nil, errors.New(hex.EncodeToString(h[:]) + " not found")
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
		MultiSigOuts: []*MultiSigOut{
			&MultiSigOut{
				N: 3,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 0,
			},
			&MultiSigOut{
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
		MultiSigOuts: []*MultiSigOut{
			&MultiSigOut{
				N: 3,
				Addresses: [][]byte{
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
					make([]byte, 32),
				},
				Value: 0,
			},
			&MultiSigOut{
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
	tx.MultiSigIns[0].PreviousTX[0] = 0x1

	seed1 := address.GenerateSeed()
	seed2 := address.GenerateSeed()
	seed3 := address.GenerateSeed()
	seed4 := address.GenerateSeed()
	a1, err := address.New(address.Height10, seed1, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	a2, err := address.New(address.Height10, seed2, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	a3, err := address.New(address.Height10, seed3, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	a4, err := address.New(address.Height10, seed4, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}

	var d [32]byte
	d[0] = 0x1
	s2561 := sha256.Sum256(a1.PublicKey())
	m[d].Outputs[0].Address = s2561[:]
	m[d].Outputs[1].Address = s2561[:]
	s2562 := sha256.Sum256(a2.PublicKey())
	m[d].MultiSigOuts[1].Addresses[0] = s2562[:]
	s2563 := sha256.Sum256(a3.PublicKey())
	m[d].MultiSigOuts[1].Addresses[1] = s2563[:]
	s2564 := sha256.Sum256(a4.PublicKey())
	m[d].MultiSigOuts[1].Addresses[2] = s2564[:]

	dat := tx.BytesForSign()
	s1 := a1.Sign(dat)
	s3 := a3.Sign(dat)
	s4 := a4.Sign(dat)
	tx.Signatures = []*Signature{
		&Signature{
			PublicKey: a1.PublicKey(),
			Sig:       s1,
		},
		&Signature{
			PublicKey: a3.PublicKey(),
			Sig:       s3,
		},
		&Signature{
			PublicKey: a4.PublicKey(),
			Sig:       s4,
		},
	}
	preh0 := tx.PreHash()
	if err := tx.PoW(); err != nil {
		t.Error(err)
	}
	if err := tx.CheckAll(m.GetTX, aklib.TestConfig); err != nil {
		t.Error(err)
	}
	preh1 := tx.PreHash()
	if !bytes.Equal(preh0, preh1) {
		t.Error("invalid pre hash")
	}
	if len(tx.NoExistHashes(m.GetTX)) != 0 {
		t.Error("invalid nonexistshashes")
	}

	m[d].MultiSigOuts[1].Addresses[1] = a2.PublicKey()
	if err := tx.checkAll(m.GetTX, aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}

	m[d].MultiSigOuts[1].Addresses[1] = a3.PublicKey()
	sig := tx.Signatures
	tx.Signatures = tx.Signatures[:1]
	if err := tx.checkAll(m.GetTX, aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}

	tx.Nonce = make([]uint32, cuckoo.ProofSize)
	tx.Inputs[0].PreviousTX[0] = 0
	tx.Inputs[1].PreviousTX[0] = 0
	tx.Previous[0][0] = 0
	tx.Previous[1][0] = 0
	tx.Signatures = sig
}
func TestTicket2(t *testing.T) {
	tx2 := &Transaction{
		Body: &Body{
			Type:   TxNormal,
			Gnonce: 0,
			Nonce:  make([]uint32, cuckoo.ProofSize),
			Time:   time.Now(),
			Inputs: []*Input{
				&Input{
					PreviousTX: make([]byte, 32),
					Index:      0,
				},
			},
			Outputs: []*Output{
				&Output{
					Address: make([]byte, 32),
					Value:   543,
				},
			},
			Previous: [][]byte{
				make([]byte, 32),
				make([]byte, 32),
			},
			TicketInput:  make([]byte, 32),
			TicketOutput: make([]byte, 32),
		},
		Signatures: []*Signature{
			&Signature{
				PublicKey: make([]byte, 65),
				Sig:       make([]byte, 32),
			},
			&Signature{
				PublicKey: make([]byte, 65),
				Sig:       make([]byte, 32),
			},
		},
	}
	tx2.Inputs[0].PreviousTX[0] = 0x1
	tx2.Previous[0][0] = 0x1
	tx2.Previous[1][0] = 0x2
	tx2.TicketInput[0] = 0x1
	tx2.TicketOutput[0] = 0x1

	seed1 := address.GenerateSeed()
	seed2 := address.GenerateSeed()
	a1, err := address.New(address.Height10, seed1, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	a2, err := address.New(address.Height10, seed2, aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}

	var d [32]byte
	d[0] = 0x1
	s256 := sha256.Sum256(a1.PublicKey())
	m[d].Outputs[0].Address = s256[:]
	s2561 := sha256.Sum256(a2.PublicKey())
	m[d].TicketOutput = s2561[:]

	dat := tx2.BytesForSign()
	s1 := a1.Sign(dat)
	s2 := a2.Sign(dat)
	tx2.Signatures[0] = &Signature{
		PublicKey: a1.PublicKey(),
		Sig:       s1,
	}
	tx2.Signatures[1] = &Signature{
		PublicKey: a2.PublicKey(),
		Sig:       s2,
	}
	if err := tx2.checkAll(m.GetTX, aklib.TestConfig, false); err != nil {
		t.Error(err)
	}
	typ, err := tx2.IsMinable(aklib.TestConfig)
	if err != nil {
		t.Error(err)
	}
	if typ != RewardTicket {
		t.Error("invalid reward type")
	}

	m[d].TicketOutput = nil
	if err := tx.checkAll(m.GetTX, aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}

	m[d].TicketOutput = a2.PublicKey()
	tx2.Signatures = tx2.Signatures[:1]
	if err := tx2.checkAll(m.GetTX, aklib.TestConfig, false); err == nil {
		t.Error("must be error")
	}
}

func BenchPoWMain0(b *testing.B) {
	n := numcpu.NumCPU()
	p := runtime.GOMAXPROCS(n)
	tx.Easiness = aklib.MainConfig.Easiness
	tx.Time = time.Time{}
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	if err := tx.PoW(); err != nil {
		b.Error(err)
	}
	if err := tx.Check(aklib.MainConfig); err != nil {
		b.Error(err)
	}
	b.Log(hex.EncodeToString(tx.Hash()))
	b.Log(tx.Nonce)
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	tx.Easiness = aklib.TestConfig.Easiness
	runtime.GOMAXPROCS(p)
}
func BenchPoWMainRand(b *testing.B) {
	n := numcpu.NumCPU()
	p := runtime.GOMAXPROCS(n)
	tx.Easiness = aklib.MainConfig.Easiness
	tx.Time = time.Now()
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	if err := tx.PoW(); err != nil {
		b.Error(err)
	}
	if err := tx.Check(aklib.MainConfig); err != nil {
		b.Error(err)
	}
	b.Log(hex.EncodeToString(tx.Hash()))
	b.Log(tx.Nonce)
	for i := range tx.Nonce {
		tx.Nonce[i] = 0
	}
	tx.Easiness = aklib.TestConfig.Easiness
	runtime.GOMAXPROCS(p)
}
