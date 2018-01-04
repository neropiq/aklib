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
	"crypto/rand"
	"testing"

	"github.com/AidosKuneen/aklib"
)

func TestAddress1(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKPRIVM1", "AKADRSM1", Height2)
	testAddress(t, aklib.TestConfig, "AKPRIVT1", "AKADRST1", Height2)
}

func TestAddress5(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKPRIVM5", "AKADRSM5", Height10)
	testAddress(t, aklib.TestConfig, "AKPRIVT5", "AKADRST5", Height10)
}
func TestAddressM8(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKPRIVM8", "AKADRSM8", Height16)
}
func TestAddressT8(t *testing.T) {
	testAddress(t, aklib.TestConfig, "AKPRIVT8", "AKADRST8", Height16)
}

func testAddress(t *testing.T, net *aklib.Config, priv, adr string, h byte) {
	pwd1 := make([]byte, 15)
	if _, err := rand.Read(pwd1); err != nil {
		panic(err)
	}
	pwd2 := make([]byte, 15)
	if _, err := rand.Read(pwd2); err != nil {
		panic(err)
	}
	seed := GenerateSeed()
	a, err := New(h, seed, net)
	if err != nil {
		t.Error(err)
	}
	s58 := a.Seed58(pwd1)
	t.Log(s58)
	_, err = NewFrom58(s58, pwd2, net)
	if err == nil {
		t.Error("should be error")
	}
	if h == Height10 {
		aa, errr := NewFrom58(s58, pwd1, net)
		if errr != nil {
			t.Error(errr)
		}
		if a.Height() != h {
			t.Error("invalid height")
		}
		if !bytes.Equal(seed, aa.Seed) {
			t.Error("invalid seed58")
		}
	} else {
		height, se, errr := from58(s58, pwd1, net)
		if errr != nil {
			t.Error(errr)
		}
		if height != h {
			t.Error("invalid height")
		}
		if !bytes.Equal(seed, se) {
			t.Error("invalid seed58")
		}
	}
	if s58[:len(priv)] != priv {
		t.Error("invalid seed58 prefix")
	}

	pk58 := a.PK58()
	t.Log(pk58)
	if pk58[:len(adr)] != adr {
		t.Error("invalid address prefix")
	}
	pk, err := FromPK58(pk58, net)
	if err != nil {
		t.Error(err)
	}
	pk2 := a.merkle.PublicKey()
	if !bytes.Equal(pk, pk2) {
		t.Error("invalid frompk58")
	}
	msg := []byte("This is a test for XMSS.")
	sig := a.Sign(msg)
	if !Verify(sig, msg, pk) {
		t.Error("signature is invalid")
	}
	b, err := a.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	c := Address{}
	if err := c.UnmarshalJSON(b); err != nil {
		t.Error(err)
	}
	if c.LeafNo() != 1 {
		t.Error("invalid unmarshal")
	}
	sig = c.Sign(msg)
	if !Verify(sig, msg, pk) {
		t.Error("signature is invalid")
	}
}

func TestAddress2(t *testing.T) {
	seed := GenerateSeed()
	a, err := New(Height10, seed, aklib.MainConfig)
	if err != nil {
		t.Error(err)
	}
	pk2 := a.merkle.PublicKey()
	msg := []byte("This is a test for XMSS.")
	sig := a.Sign(msg)
	if !Verify(sig, msg, pk2) {
		t.Error("signature is invalid")
	}
}
