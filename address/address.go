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
	"log"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/xmss"
)

var (
	prefixPrivString = "AKPRI"
	prefixAdrsString = "AKADR"
)

//Height represents height of Merkle Tree for XMSS
const (
	Height10 = iota
	Height16
	Height20
)

var heights = []uint32{10, 16, 20}

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
	height byte
	config *aklib.Config
	merkle *xmss.Merkle
	Seed   []byte
}

//New returns Address struct.
func New(h byte, seed []byte, config *aklib.Config) (*Address, error) {
	if h > Height20 {
		return nil, errors.New("invalid height")
	}
	return &Address{
		height: h,
		merkle: xmss.NewMerkle(heights[h], seed),
		Seed:   seed,
		config: config,
	}, nil
}

//Seed58 returns base58 encoded seed.
func (a *Address) Seed58() string {
	pref := a.config.PrefixPriv[a.height]
	s := make([]byte, len(a.Seed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], a.Seed)
	return prefixPrivString + encode58(s)
}

//NewFrom58 returns Address struct with base58 encoded seed.
func NewFrom58(height byte, seed58 string, cfg *aklib.Config) (*Address, error) {
	if height > Height20 {
		return nil, errors.New("invalid height")
	}
	if prefixPrivString != seed58[:len(prefixPrivString)] {
		return nil, errors.New("invalid prefix string in seed")
	}
	seed, err := decode58(seed58[len(prefixPrivString):])
	if err != nil {
		return nil, err
	}
	pref := cfg.PrefixPriv[height]
	if !bytes.Equal(seed[:len(pref)], pref) {
		log.Println((seed[:len(pref)]), pref)
		return nil, errors.New("invalid prefix bytes in seed")
	}
	return New(height, seed[len(pref):], cfg)
}

//PublicKey returns public key.
func (a *Address) PublicKey() []byte {
	return a.merkle.PublicKey()
}

//PK58 returns base58 encoded public key.
func (a *Address) PK58() string {
	pref := a.config.PrefixAdrs[a.height]
	pub := a.merkle.PublicKey()
	p := make([]byte, len(pub)+len(pref))
	copy(p, pref)
	copy(p[len(pref):], pub)
	return prefixAdrsString + encode58(p)
}

//FromPK58 returns decode public key from base58 encoded string.
func FromPK58(height byte, pub58 string, cfg *aklib.Config) ([]byte, error) {
	if height > Height20 {
		return nil, errors.New("invalid height")
	}
	if prefixAdrsString != pub58[:len(prefixAdrsString)] {
		return nil, errors.New("invalid prefix string in public key")
	}
	pub, err := decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, err
	}
	pref := cfg.PrefixAdrs[height]
	if !bytes.Equal(pub[:len(pref)], pref) {
		return nil, errors.New("invalid prefix bytes in public key")
	}
	return pub[len(pref):], nil
}

//MarshalJSON  marshals Address into valid JSON.
func (a *Address) MarshalJSON() ([]byte, error) {
	s := struct {
		Height byte
		Seed   []byte
		Merkle *xmss.Merkle
		Config *aklib.Config
	}{
		Height: a.height,
		Seed:   a.Seed,
		Merkle: a.merkle,
		Config: a.config,
	}
	return json.Marshal(&s)
}

//UnmarshalJSON  unmarshals JSON to Address.
func (a *Address) UnmarshalJSON(b []byte) error {
	s := struct {
		Height byte
		Seed   []byte
		Merkle *xmss.Merkle
		Config *aklib.Config
	}{}
	err := json.Unmarshal(b, &s)
	a.Seed = s.Seed
	a.merkle = s.Merkle
	a.config = s.Config
	a.height = s.Height
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
