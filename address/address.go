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
	"encoding/json"
	"errors"

	"github.com/AidosKuneen/xmss"
)

var (
	prefixPrivBytes = [][]byte{
		[]byte{0xbf, 0x9d}, //MainNet, "VM" in base58
		[]byte{0xc0, 0x50}, //TestNet, "VT" in base58
	}
	prefixAdrsBytes = [][]byte{
		[]byte{0xab, 0x55}, //MainNet, "SM" in base58
		[]byte{0xac, 0x8},  //TestNet, "ST" in base58
	}

	prefixPrivString = "AKPRI"
	prefixAdrsString = "AKADR"
)

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
	net    int
	merkle *xmss.Merkle
	Seed   []byte
}

//New returns Address struct.
func New(h uint32, seed []byte, net int) *Address {
	return &Address{
		merkle: xmss.NewMerkle(h, seed),
		Seed:   seed,
		net:    net,
	}
}

//Seed58 returns base58 encoded seed.
func (a *Address) Seed58() string {
	s := make([]byte, len(a.Seed)+2)
	copy(s, prefixPrivBytes[a.net])
	copy(s[2:], a.Seed)
	return prefixPrivString + encode58(s)
}

//NewFrom58 returns Address struct with base58 encoded seed.
func NewFrom58(h uint32, seed58 string, net int) (*Address, error) {
	if prefixPrivString != seed58[:len(prefixPrivString)] {
		return nil, errors.New("invalid prefix string in seed")
	}
	seed, err := decode58(seed58[len(prefixPrivString):])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(seed[:2], prefixPrivBytes[net]) {
		return nil, errors.New("invalid prefix bytes in seed")
	}
	return New(h, seed[2:], net), nil
}

//PK58 returns base58 encoded public key.
func (a *Address) PK58() string {
	pub := a.merkle.PublicKey()
	p := make([]byte, len(pub)+2)
	copy(p, prefixAdrsBytes[a.net])
	copy(p[2:], pub)
	return prefixAdrsString + encode58(p)
}

//FromPK58 returns decode public key from base58 encoded string.
func FromPK58(pub58 string, net int) ([]byte, error) {
	if prefixAdrsString != pub58[:len(prefixAdrsString)] {
		return nil, errors.New("invalid prefix string in public key")
	}
	pub, err := decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(pub[:2], prefixAdrsBytes[net]) {
		return nil, errors.New("invalid prefix bytes in public key")
	}
	return pub[2:], nil
}

//MarshalJSON  marshals Address into valid JSON.
func (a *Address) MarshalJSON() ([]byte, error) {
	s := struct {
		Seed   []byte
		Merkle *xmss.Merkle
		Net    int
	}{
		Seed:   a.Seed,
		Merkle: a.merkle,
		Net:    a.net,
	}
	return json.Marshal(&s)
}

//UnmarshalJSON  unmarshals JSON to Address.
func (a *Address) UnmarshalJSON(b []byte) error {
	s := struct {
		Seed   []byte
		Merkle *xmss.Merkle
		Net    int
	}{}
	err := json.Unmarshal(b, &s)
	a.Seed = s.Seed
	a.merkle = s.Merkle
	a.net = s.Net
	return err
}

//LeafNo returns leaf number we will use next.
func (a *Address) LeafNo() uint32 {
	return a.merkle.Leaf
}

//NextLeaf adds leaf number and refresh auth and stack in MerkleTree..
func (a *Address) NextLeaf(n int) {
	for i := 0; i < n; i++ {
		a.merkle.Traverse()
	}
}

//Sign signs msg.
func (a *Address) Sign(msg []byte) []byte {
	return a.merkle.Sign(msg)
}

//Verify verifies msg with signature bsig.
func Verify(bsig, msg, bpk []byte) bool {
	return xmss.Verify(bsig, msg, bpk)
}

//GenerateSeed generates a new 32 bytes seed.
func GenerateSeed() []byte {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return seed
}
