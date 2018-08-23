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
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/aklib/arypack"
	"github.com/AidosKuneen/cuckoo"
	"github.com/AidosKuneen/numcpu"
)

var (
	zero = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	one  = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	two  = []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	a    [5]*address.Address
)

func TestMain(m *testing.M) {
	var err error
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	for i := range a {
		a[i], err = address.New(aklib.DebugConfig, false)
		if err != nil {
			panic(err)
		}
	}
	c := m.Run()
	os.Exit(c)
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
	tr := New(aklib.TestConfig, one, zero)

	if err := tr.Check(aklib.TestConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	t.Log(hex.EncodeToString(tr.Hash()))
	if err := tr.Check(aklib.TestConfig, TypeNormal); err != nil {
		t.Error(err)
	}

	tr = New(aklib.TestConfig, zero, zero)
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.TestConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}
}

func TestTX(t *testing.T) {
	tr := New(aklib.DebugConfig, zero, one)
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}

	tr.Nonce = make([]uint32, cuckoo.ProofSize)
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("must be error")
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.TestConfig, TypeNormal); err == nil {
		t.Error("must be error")
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err == nil {
		t.Error("must be error")
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardTicket); err == nil {
		t.Error("must be error")
	}

	tr = New(aklib.DebugConfig, zero, zero)
	tr.AddInput(zero, 0)
	tr.AddInput(zero, 0)
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("must be error")
	}

	tr = New(aklib.DebugConfig, zero, one)
	sig, err2 := tr.Signature(a[0])
	if err2 != nil {
		t.Error(err2)
	}
	tr.AddSig(sig)
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err != nil {
		t.Error(err)
	}

	tr = New(aklib.DebugConfig, zero, one)
	tr.Time = time.Now().Add(time.Hour)
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("must be error")
	}

	b := arypack.Marshal(tr.Body)
	s := arypack.Marshal(tr.Signatures)
	var b2 Body
	if err := arypack.Unmarshal(b, &b2); err != nil {
		t.Error(err)
	}
	var s2 Signatures
	if err := arypack.Unmarshal(s, &s2); err != nil {
		t.Error(err)
	}
	t2 := &Transaction{
		Body:       &b2,
		Signatures: s2,
	}
	h := tr.Hash()
	h2 := t2.Hash()
	if !bytes.Equal(h, h2) {
		t.Error("pack/unpack is incorrect")
	}
	t2 = tr.Clone()
	h2 = t2.Hash()
	if !bytes.Equal(h, h2) {
		t.Error("clone is incorrect")
	}

	tr = New(aklib.DebugConfig, zero, one)
	tr.LockTime = time.Now().Add(time.Hour)
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}

	tr = New(aklib.DebugConfig, zero, one)
	tr.HashType = 0x1f

	if err := tr.Sign(a[0]); err == nil {
		t.Error("should be error")
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err == nil {
		t.Error("should be error")
	}

	tr = New(aklib.DebugConfig, zero, one)
	tr.HashType = 0x11
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 1234); err != nil {
		t.Error(err)
	}
	hs1, err2 := tr.bytesForSign()
	if err2 != nil {
		t.Error(err2)
	}
	tr.Outputs[0].Address = nil
	tr.Outputs[0].Value = 234
	hs2, err2 := tr.bytesForSign()
	if err2 != nil {
		t.Error(err2)
	}
	if !bytes.Equal(hs1, hs2) {
		t.Error("invalid hashtype")
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err == nil {
		t.Error("should be error")
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}

	tr = NewMinableFee(aklib.DebugConfig, zero, one)
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 1234); err != nil {
		t.Error(err)
	}
	if err := tr.AddOutput(aklib.DebugConfig, "", 2345); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err != nil {
		t.Error(err)
	}

	tr.HashType = 0
	if err := tr.Sign(a[0]); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeRewardFee); err == nil {
		t.Error("invalid isminable")
	}

	tr = New(aklib.DebugConfig, zero, one)
	tr.HashType = 0x20
	tr.TicketInput = one
	tr.TicketInput = two
	hs1, err2 = tr.bytesForSign()
	if err2 != nil {
		t.Error(err2)
	}
	tr.TicketOutput = nil
	hs2, err2 = tr.bytesForSign()
	if err2 != nil {
		t.Error(err2)
	}
	if !bytes.Equal(hs1, hs2) {
		t.Error("invalid hashtype")
	}
	tr = NewMinableTicket(aklib.DebugConfig, one, zero, one)
	if err := tr.Check(aklib.DebugConfig, TypeRewardTicket); err != nil {
		t.Error(err)
	}
}
func TestTicket(t *testing.T) {
	ticket, err := IssueTicket(aklib.DebugConfig, a[0].Address(aklib.DebugConfig), zero)
	if err != nil {
		t.Error(err)
	}
	if err := ticket.Check(aklib.DebugConfig, TypeNormal); err != nil {
		t.Error(err)
	}
	if !isValidHash(ticket.Hash(), aklib.DebugConfig.TicketEasiness) {
		t.Error("invlaid ticket hash")
	}
	ticket.Easiness = aklib.DebugConfig.Easiness
	if err := ticket.PoW(); err != nil {
		t.Error(err)
	}
	if err := ticket.Check(aklib.DebugConfig, TypeNormal); err == nil {
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

func TestTX2(t *testing.T) {
	tr := New(aklib.DebugConfig, one, two)
	tr.AddInput(one, 0)
	tr.AddInput(one, 1)
	tr.AddMultisigIn(one, 1)
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 111); err != nil {
		t.Error(err)
	}
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 222); err != nil {
		t.Error(err)
	}
	if err := tr.AddMultisigOut(aklib.DebugConfig, 3, 2,
		a[0].Address58(aklib.DebugConfig), a[1].Address58(aklib.DebugConfig), a[2].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}
	if err := tr.AddMultisigOut(aklib.DebugConfig, 2, 331,
		a[0].Address58(aklib.DebugConfig), a[1].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}

	var d1, d2 [32]byte
	d1[0] = 0x1
	d2[0] = 0x2
	m := make(store)
	m[d1] = New(aklib.DebugConfig).Body
	if err := m[d1].AddOutput(aklib.DebugConfig, a[1].Address58(aklib.DebugConfig), 543); err != nil {
		t.Error(err)
	}
	if err := m[d1].AddOutput(aklib.DebugConfig, a[1].Address58(aklib.DebugConfig), 0); err != nil {
		t.Error(err)
	}
	if err := m[d1].AddMultisigOut(aklib.DebugConfig, 2, 0,
		a[2].Address58(aklib.DebugConfig), a[3].Address58(aklib.DebugConfig), a[4].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}
	if err := m[d1].AddMultisigOut(aklib.DebugConfig, 2, 123,
		a[2].Address58(aklib.DebugConfig), a[3].Address58(aklib.DebugConfig), a[4].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}
	m[d2] = New(aklib.DebugConfig).Body

	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[3]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[4]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err != nil {
		t.Fatal(err)
	}
	tr.Outputs[0].Value = 110
	tr.Signatures = tr.Signatures[:0]
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[3]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[4]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err == nil {
		t.Fatal("should be error")
	}
	tr.Outputs[0].Value = 111

	tr.Parent = tr.Parent[:1]
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err == nil {
		t.Error("should be error")
	}
	tr.Signatures = tr.Signatures[:0]
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[3]); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[4]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err != nil {
		t.Error(err)
	}

	m[d1].MultiSigOuts[1].Addresses[1] = a[2].Address(aklib.DebugConfig)
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err == nil {
		t.Error("must be error")
	}

	m[d1].MultiSigOuts[1].Addresses[1] = a[3].Address(aklib.DebugConfig)
	sig := tr.Signatures
	tr.Signatures = tr.Signatures[:1]
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err == nil {
		t.Error("must be error")
	}
	tr.Signatures = sig
	tr.Outputs[0].Value = 1
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err == nil {
		t.Error("must be error")
	}
}
func TestTicket2(t *testing.T) {
	tx2 := NewMinableTicket(aklib.DebugConfig, one, one)
	tx2.AddInput(one, 0)
	tx2.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 543)

	a1, err := address.New(aklib.DebugConfig, false)
	if err != nil {
		t.Error(err)
	}
	a2, err := address.New(aklib.DebugConfig, false)
	if err != nil {
		t.Error(err)
	}

	var d [32]byte
	d[0] = 0x1
	m := make(store)
	m[d] = New(aklib.DebugConfig).Body
	if err := m[d].AddOutput(aklib.DebugConfig, a1.Address58(aklib.DebugConfig), 543); err != nil {
		t.Error(err)
	}
	m[d].TicketOutput = a1.Address(aklib.DebugConfig)

	if err := tx2.Sign(a1); err != nil {
		t.Error(err)
	}
	if err := tx2.CheckAll(aklib.DebugConfig, m.GetTX, TypeRewardTicket); err != nil {
		t.Error(err)
	}
	if err := tx2.Sign(a2); err != nil {
		t.Error(err)
	}
	if err := tx2.CheckAll(aklib.DebugConfig, m.GetTX, TypeRewardTicket); err == nil {
		t.Error("should be error")
	}
	tx2.Signatures = tx2.Signatures[:1]

	m[d].TicketOutput = a2.Address(aklib.DebugConfig)
	if err := tx2.CheckAll(aklib.DebugConfig, m.GetTX, TypeRewardTicket); err == nil {
		t.Error("must be error")
	}

	m[d].TicketOutput = a1.Address(aklib.DebugConfig)
	tx2.TicketOutput = a2.Address(aklib.DebugConfig)

	if err := tx2.PoW(); err != nil {
		t.Error(err)
	}
	if err := tx2.CheckAll(aklib.DebugConfig, m.GetTX, TypeNormal); err != nil {
		t.Error(err)
	}
}

func TestTX3(t *testing.T) {
	tr := New(aklib.DebugConfig, zero, one)
	tr.AddMultisigIn(zero, 1)
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err != nil {
		t.Error(err)
	}

	tr.AddMultisigIn(zero, 1)
	tr.Signatures = tr.Signatures[:0]
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}

	tr = New(aklib.DebugConfig, zero, one)
	if err := tr.AddMultisigOut(aklib.DebugConfig, 2, 0, a[1].Address58(aklib.DebugConfig)); err == nil {
		t.Error("should be error")
	}
	tr.MultiSigOuts = append(tr.MultiSigOuts, &MultiSigOut{
		M:         2,
		Addresses: []Address{zero},
		Value:     1,
	})
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	if err := tr.Check(aklib.DebugConfig, TypeNormal); err == nil {
		t.Error("should be error")
	}

}
func TestMarshal(t *testing.T) {
	tr := New(aklib.DebugConfig, one, two)
	tr.AddInput(one, 0)
	tr.AddInput(one, 1)
	tr.AddMultisigIn(one, 1)
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 111); err != nil {
		t.Error(err)
	}
	if err := tr.AddOutput(aklib.DebugConfig, a[0].Address58(aklib.DebugConfig), 222); err != nil {
		t.Error(err)
	}
	if err := tr.AddMultisigOut(aklib.DebugConfig, 3, 2,
		a[0].Address58(aklib.DebugConfig), a[1].Address58(aklib.DebugConfig), a[2].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}
	if err := tr.AddMultisigOut(aklib.DebugConfig, 2, 331,
		a[0].Address58(aklib.DebugConfig), a[1].Address58(aklib.DebugConfig)); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a[1]); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	dat, err := json.Marshal(tr)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(dat))
	var tr2 Transaction
	if err := json.Unmarshal(dat, &tr2); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(tr.Hash(), tr2.Hash()) {
		t.Error("invalid json")
	}
}

func BenchmarkPoWMain0(b *testing.B) {
	benchPoW(b, false)
}
func BenchmarkPoWMainRand(b *testing.B) {
	benchPoW(b, true)
}

func benchPoW(b *testing.B, r bool) {
	n := numcpu.NumCPU()
	p := runtime.GOMAXPROCS(n)
	tr := New(aklib.MainConfig, zero)
	var seed1 []byte
	var err error
	if r {
		tr.Time = time.Time{}
		seed1 = make([]byte, 32)
	} else {
		seed1 = address.GenerateSeed32()
	}
	a1, err := address.NewFromSeed(aklib.MainConfig, seed1, false)
	if err != nil {
		b.Error(err)
	}
	if err := tr.Sign(a1); err != nil {
		b.Error(err)
	}
	if err := tr.PoW(); err != nil {
		b.Error(err)
	}
	if err := tr.Check(aklib.MainConfig, TypeNormal); err != nil {
		b.Error(err)
	}
	b.Log(hex.EncodeToString(tr.Hash()))
	b.Log(tr.Nonce)
	for i := range tr.Nonce {
		tr.Nonce[i] = 0
	}
	runtime.GOMAXPROCS(p)
}
