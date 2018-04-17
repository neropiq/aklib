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
	"encoding/json"
	"testing"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/vmihailenco/msgpack"
)

func TestNode2(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKNKEYM2", "AKNODEM2", Height40)
	testAddress(t, aklib.TestConfig, "AKNKEYT2", "AKNODET2", Height40)
}
func TestNode3(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKNKEYM3", "AKNODEM3", Height60)
	testAddress(t, aklib.TestConfig, "AKNKEYT3", "AKNODET3", Height60)
}

func testAddress(t *testing.T, net *aklib.Config, priv, adr string, h byte) {
	seed := address.GenerateSeed()
	var a *Address
	var err error
	a, err = New(h, seed, net)
	if err != nil {
		t.Error(err)
	}
	s58 := a.Seed58()
	t.Log(s58)
	_, err = NewFrom58(s58, net)
	if err != nil {
		t.Error(err)
	}
	height, se, errr := from58(s58, net)
	if errr != nil {
		t.Error(errr)
	}
	if height != h {
		t.Error("invalid height")
	}
	if !bytes.Equal(seed, se) {
		t.Error("invalid seed58")
	}
	if s58[:len(priv)] != priv {
		t.Error("invalid seed58 prefix")
	}

	pk58 := a.Address58()
	t.Log(pk58)
	if pk58[:len(adr)] != adr {
		t.Error("invalid address prefix")
	}
	pk, h2, err := FromAddress58(pk58, net)
	if err != nil {
		t.Error(err)
	}
	if h2 != heights[h] {
		t.Error("invalid height")
	}
	pk2 := a.Address()
	if !bytes.Equal(pk, pk2) {
		t.Error("invalid frompk58")
	}
	msg := []byte("This is a test for XMSS.")
	sig := a.Sign(msg)
	if !Verify(sig, msg, a.PublicKey()) {
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
	if !Verify(sig, msg, a.PublicKey()) {
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
	if !Verify(sig, msg, a.PublicKey()) {
		t.Error("signature is invalid")
	}

}
