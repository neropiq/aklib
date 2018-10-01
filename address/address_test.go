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
	"strings"
	"testing"

	"github.com/AidosKuneen/aklib"
)

func TestMultisigAddress(t *testing.T) {
	testMultisigAddress(t, aklib.MainConfig, "AKMSIGM")
	testMultisigAddress(t, aklib.TestConfig, "AKMSIGT")
	testMultisigAddress(t, aklib.DebugConfig, "AKMSIGD")
}

func TestNodeAddress(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKNODEM", true)
	testAddress(t, aklib.TestConfig, "AKNODET", true)
	testAddress(t, aklib.DebugConfig, "AKNODED", true)
}
func TestAddress(t *testing.T) {
	testAddress(t, aklib.MainConfig, "AKADRSM", false)
	testAddress(t, aklib.TestConfig, "AKADRST", false)
	testAddress(t, aklib.DebugConfig, "AKADRSD", false)
}

func testMultisigAddress(t *testing.T, net *aklib.Config, adr string) {
	pwd1 := make([]byte, 15)
	if _, err := rand.Read(pwd1); err != nil {
		panic(err)
	}
	pwd2 := make([]byte, 15)
	if _, err := rand.Read(pwd2); err != nil {
		panic(err)
	}
	var a1, a2, a3 *Address
	var err error
	seed1 := GenerateSeed32()
	seed2 := GenerateSeed32()
	seed3 := GenerateSeed32()
	a1, err = New(net, seed1)
	if err != nil {
		t.Error(err)
	}
	a2, err = New(net, seed2)
	if err != nil {
		t.Error(err)
	}
	a3, err = New(net, seed3)
	if err != nil {
		t.Error(err)
	}
	msig := MultisigAddress(net, 2, a1.Address(net), a2.Address(net), a3.Address(net))
	if !strings.HasPrefix(msig, adr) {
		t.Error("invalid msig adr", msig)
	}
	t.Log(msig)
	badr, err := ParseMultisigAddress(net, msig)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(badr, MultisigAddressByte(net, 2, a1.Address(net), a2.Address(net), a3.Address(net))) {
		t.Error("invalid madr")
	}
}

func testAddress(t *testing.T, net *aklib.Config, adr string, isNode bool) {
	pwd1 := make([]byte, 15)
	if _, err := rand.Read(pwd1); err != nil {
		panic(err)
	}
	pwd2 := make([]byte, 15)
	if _, err := rand.Read(pwd2); err != nil {
		panic(err)
	}
	var a *Address
	var err error
	seed := GenerateSeed32()

	if isNode {
		a, err = NewNode(net, seed)
	} else {
		a, err = New(net, seed)
	}
	if err != nil {
		t.Error(err)
	}
	t.Log(a.PrivateKey)

	pk58 := a.Address58(net)
	t.Log(pk58)
	if pk58[:len(adr)] != adr {
		t.Error("invalid address prefix")
	}
	pk, fr2, err := ParseAddress58(net, pk58)
	if err != nil {
		t.Error(err)
	}
	if fr2 != isNode {
		t.Error("invalid address type")
	}
	pk581, err := Address58(net, a.Address(net))
	if err != nil {
		t.Error(err)
	}
	if pk58 != pk581 {
		t.Error("invalid To58")
	}
	pk2 := a.Address(net)
	if !bytes.Equal(pk, pk2) {
		t.Error("invalid frompk58", hex.EncodeToString(pk), hex.EncodeToString(pk2))
	}
	msg := []byte("This is a test for glyph.")
	sig, err := a.Sign(msg)
	if err != nil {
		t.Log(a.PrivateKey)
		t.Error(err)
	}
	if err = sig.Verify(msg); err != nil {
		t.Error(err)
	}
}

func TestAddress2(t *testing.T) {
	msg := []byte("This is a test for XMSS.")

	for _, fr := range []bool{true, false} {
		var a *Address
		var err error
		seed := GenerateSeed32()
		if fr {
			a, err = NewNode(aklib.MainConfig, seed)
		} else {
			a, err = New(aklib.MainConfig, seed)
		}
		if err != nil {
			t.Error(err)
		}
		sig, err := a.Sign(msg)
		if err != nil {
			t.Error(err)
		}
		if err := sig.Verify(msg); err != nil {
			t.Error(err)
		}
	}
}
