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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/AidosKuneen/aklib"
	"github.com/vmihailenco/msgpack"
)

func TestAddress1(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKADRSM1", Height2)
	testAddress(t, aklib.TestConfig, "AKADRST1", Height2)
}

func TestAddress5(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKADRSM5", Height10)
	testAddress(t, aklib.TestConfig, "AKADRST5", Height10)
}
func TestAddressM8(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKADRSM8", Height16)
}
func TestAddressT8(t *testing.T) {
	testAddress(t, aklib.TestConfig, "AKADRST8", Height16)
}

func testAddress(t *testing.T, net *aklib.Config, adr string, h byte) {
	pwd1 := make([]byte, 15)
	if _, err := rand.Read(pwd1); err != nil {
		panic(err)
	}
	pwd2 := make([]byte, 15)
	if _, err := rand.Read(pwd2); err != nil {
		panic(err)
	}
	seed := GenerateSeed()
	var a *Address
	var err error
	a, err = New(h, seed, net)
	if err != nil {
		t.Error(err)
	}

	pk58 := a.Address58()
	t.Log(pk58)
	if pk58[:len(adr)] != adr {
		t.Error("invalid address prefix")
	}
	pk, h2, err := ParseAddress58(pk58, net)
	if err != nil {
		t.Error(err)
	}
	if h2 != h {
		t.Error("invalid height")
	}
	pk581 := To58(a.Address())
	if pk58 != pk581 {
		t.Error("invalid To58")
	}
	pk2 := a.Address()
	if !bytes.Equal(pk, pk2) {
		t.Error("invalid frompk58", hex.EncodeToString(pk), hex.EncodeToString(pk2))
	}
	msg := []byte("This is a test for XMSS.")
	sig := a.Sign(msg)
	if !Verify(sig, msg) {
		t.Error("signature is invalid")
	}

	b, err := json.Marshal(a)
	if err != nil {
		t.Error(err)
	}
	var c Address
	err = json.Unmarshal(b, &c)
	if err != nil {
		t.Error(err)
	}
	if c.LeafNo() != 1 {
		t.Error("invalid unmarshal")
	}
	sig = c.Sign(msg)
	if !Verify(sig, msg) {
		t.Error("signature is invalid")
	}

	mb, err := msgpack.Marshal(a)
	if err != nil {
		t.Error(err)
	}
	var mc Address
	if err := msgpack.Unmarshal(mb, &mc); err != nil {
		t.Error(err)
	}
	if mc.LeafNo() != 1 {
		t.Error("invalid unmarshal")
	}
	sig = mc.Sign(msg)
	if !Verify(sig, msg) {
		t.Error("signature is invalid")
	}
}

func TestAddress2(t *testing.T) {
	seed := GenerateSeed()
	a, err := New(Height10, seed, aklib.MainConfig)
	if err != nil {
		t.Error(err)
	}
	msg := []byte("This is a test for XMSS.")
	sig := a.Sign(msg)
	if !Verify(sig, msg) {
		t.Error("signature is invalid")
	}
}
