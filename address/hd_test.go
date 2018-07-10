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
	"log"
	"strings"
	"testing"

	"github.com/AidosKuneen/aklib"
)

func TestHD(t *testing.T) {
	masterkey := make([]byte, 32)
	if _, err := rand.Read(masterkey); err != nil {
		t.Fatal(err)
	}
	seed200 := HDseed(masterkey, Height2, 0, 0)
	seed253 := HDseed(masterkey, Height2, 5, 3)
	seed800 := HDseed(masterkey, Height16, 0, 0)
	seedA00 := HDseed(masterkey, Height20, 0, 0)
	seedA002 := HDseed(masterkey, Height20, 0, 0)
	log.Println(hex.EncodeToString(seed200))
	log.Println(hex.EncodeToString(seed253))
	log.Println(hex.EncodeToString(seed800))
	log.Println(hex.EncodeToString(seedA00))
	if bytes.Equal(seed200, seed253) {
		t.Error("should not be equal")
	}
	if bytes.Equal(seed200, seed800) {
		t.Error("should not be equal")
	}
	if bytes.Equal(seed200, seedA00) {
		t.Error("should not be equal")
	}
	if bytes.Equal(seed253, seed800) {
		t.Error("should not be equal")
	}
	if bytes.Equal(seed253, seedA00) {
		t.Error("should not be equal")
	}
	if bytes.Equal(seed800, seedA00) {
		t.Error("should not be equal")
	}
	if !bytes.Equal(seedA00, seedA002) {
		t.Error("should be equal")
	}

}

func TestHDSeed(t *testing.T) {
	seed := GenerateSeed()
	seed200 := HDseed(seed, Height2, 0, 0)
	pwd1 := []byte("qewrty123")
	s58 := HDSeed58(aklib.TestConfig, seed200, pwd1)
	t.Log(s58)
	if !strings.HasPrefix(s58, "AKPRIVT1") {
		t.Error("invaild prefix")
	}
	if _, err := HDFrom58(s58, []byte("wrong"), aklib.DebugConfig); err == nil {
		t.Error("invalid from58")
	}
	rec, err := HDFrom58(s58, pwd1, aklib.DebugConfig)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(rec, seed200) {
		t.Error("invalid from58")
	}

	s58 = HDSeed58(aklib.MainConfig, seed200, pwd1)
	t.Log(s58)
	if !strings.HasPrefix(s58, "AKPRIVM1") {
		t.Error("invaild prefix")
	}
}
