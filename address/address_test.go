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
	"fmt"
	"log"
	"testing"

	"github.com/AidosKuneen/aklib"
)

func TestAddress(t *testing.T) {
	testAddress(t, aklib.MainNet, "AKPRIVM", "AKADRSM")
	testAddress(t, aklib.TestNet, "AKPRIVT", "AKADRST")
}

func testAddress(t *testing.T, net int, priv, adr string) {
	seed := GenerateSeed()
	a := New(10, seed, net)
	s58 := a.Seed58()
	aa, err := NewFrom58(10, s58, net)
	if err != nil {
		t.Error(err)
	}
	if s58[:7] != priv {
		t.Error("invalid seed58 prefix")
	}
	if !bytes.Equal(seed, aa.Seed) {
		t.Error("invalid seed58")
	}
	fmt.Println(s58)

	pk58 := a.PK58()
	fmt.Println(pk58)
	if pk58[:7] != adr {
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
	sig := aa.Sign(msg)
	if !Verify(sig, msg, pk) {
		log.Println("signature is invalid")
	}
	b, err := aa.MarshalJSON()
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
		log.Println("signature is invalid")
	}
}
