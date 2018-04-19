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
	"testing"
)

func TestHD(t *testing.T) {
	masterkey := make([]byte, 32)
	seed := make([]byte, 32)
	if _, err := rand.Read(masterkey); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(seed); err != nil {
		t.Fatal(err)
	}
	seed200 := HDseed(seed, masterkey, Height2, 0, 0)
	seed253 := HDseed(seed, masterkey, Height2, 5, 3)
	seed800 := HDseed(seed, masterkey, Height16, 0, 0)
	seedA00 := HDseed(seed, masterkey, Height20, 0, 0)
	seedA002 := HDseed(seed, masterkey, Height20, 0, 0)
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
