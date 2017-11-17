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
	"testing"

	"github.com/AidosKuneen/aklib"
)

func TestAddress(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKPRIVM5", "AKADRSM5", Height10)
	testAddress(t, aklib.TestConfig, "AKPRIVT5", "AKADRST5", Height10)
	testAddress(t, aklib.MainConfig, "AKPRIVM8", "AKADRSM8", Height16)
	testAddress(t, aklib.TestConfig, "AKPRIVT8", "AKADRST8", Height16)
}

func testAddress(t *testing.T, net *aklib.Config, priv, adr string, h byte) {
	seed := GenerateSeed()
	a, err := New(h, seed, net)
	if err != nil {
		t.Error(err)
	}
	s58 := a.Seed58()
	if h == Height10 {
		aa, err := NewFrom58(s58, net)
		if err != nil {
			t.Error(err)
		}
		if a.Height() != h {
			t.Error("invalid height")
		}
		if !bytes.Equal(seed, aa.Seed) {
			t.Error("invalid seed58")
		}
	} else {
		height, se, err := from58(s58, net)
		if err != nil {
			t.Error(err)
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
