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
	"errors"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	sha256 "github.com/AidosKuneen/sha256-simd"
	"github.com/AidosKuneen/xmss"
	"github.com/vmihailenco/msgpack"
)

var (
	prefixNodeString = "AKNOD"
	prefixNkeyString = "AKNKE"
)

//Height represents height of Merkle Tree for XMSS
const (
	Height20 = iota
	Height40
	Height60
)

var heights = []byte{20, 40, 60}

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
	prefixPriv []byte
	prefixPub  []byte
	privateKey *xmss.PrivKeyMT
	Seed       []byte
	Height     byte
}

//New returns Address struct for node ID..
func New(h byte, seed []byte, config *aklib.Config) (*Address, error) {
	if h > Height60 {
		return nil, errors.New("invalid height")
	}
	mt, err := xmss.NewPrivKeyMT(seed, uint32(heights[h]), uint32(heights[h]/10))
	if err != nil {
		return nil, err
	}
	return &Address{
		privateKey: mt,
		Seed:       seed,
		prefixPriv: config.PrefixNkey[h],
		prefixPub:  config.PrefixNode[h],
		Height:     heights[h],
	}, nil
}

//Seed58 returns base58-encoded seed..
func (a *Address) Seed58() string {
	out := make([]byte, len(a.Seed)+4)
	copy(out, a.Seed)
	hash := sha256.Sum256(a.Seed)
	hash = sha256.Sum256(hash[:])

	copy(out[len(a.Seed):], hash[0:4])

	pref := a.prefixPriv
	s := make([]byte, len(a.Seed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], a.Seed)
	return prefixNkeyString + address.Encode58(s)
}

func from58(seed58 string, cfg *aklib.Config) (byte, []byte, error) {
	if seed58[:len(prefixNkeyString)] != prefixNkeyString {
		return 0, nil, errors.New("invalid prefix string in seed")
	}
	eseed, err := address.Decode58(seed58[len(prefixNkeyString):])
	if err != nil {
		return 0, nil, err
	}
	var height byte
	for ; height <= Height60; height++ {
		pref := cfg.PrefixNkey[height]
		if bytes.Equal(eseed[:len(pref)], pref) {
			break
		}
	}
	if height > Height60 {
		return 0, nil, errors.New("invalid prefix bytes in seed")
	}
	eseed = eseed[len(cfg.PrefixNkey[height]):]
	return height, eseed, nil

}

//NewFrom58 returns Address struct with base58 encoded seed.
func NewFrom58(seed58 string, cfg *aklib.Config) (*Address, error) {
	height, seed, err := from58(seed58, cfg)
	if err != nil {
		return nil, err
	}
	return New(height, seed, cfg)
}

//PublicKey returns public key.
func (a *Address) PublicKey() []byte {
	return a.privateKey.PublicKey()
}

//PK58 returns base58 encoded public key.
func (a *Address) PK58() string {
	pref := a.prefixPub
	pub := a.privateKey.PublicKey()
	p := make([]byte, len(pub)-1+len(pref))
	copy(p, pref)
	copy(p[len(pref):], pub[1:])
	return prefixNodeString + address.Encode58(p)
}

//FromPK58 returns decode public key from base58 encoded string.
func FromPK58(pub58 string, cfg *aklib.Config) ([]byte, error) {
	if pub58[:len(prefixNodeString)] != prefixNodeString {
		return nil, errors.New("invalid prefix string in public key")
	}
	pub, err := address.Decode58(pub58[len(prefixNodeString):])
	if err != nil {
		return nil, err
	}
	var height byte
	for ; height <= Height60; height++ {
		pref := cfg.PrefixNode[height]
		if bytes.Equal(pub[:len(pref)], pref) {
			break
		}
	}
	if height > Height60 {
		return nil, errors.New("invalid prefix bytes in seed")
	}
	pub = pub[len(cfg.PrefixAdrs[height]):]
	var header byte
	header, err = xmss.PublickeyMTHeader(uint32(heights[height]), uint32(heights[height]/10))
	if err != nil {
		return nil, err
	}
	pub = append([]byte{header}, pub...)
	return pub, nil
}

type adr struct {
	Seed       []byte
	PrivateKey *xmss.PrivKeyMT
	PrefixPriv []byte
	PrefixPub  []byte
	Height     byte
}

func (a *Address) exports() *adr {
	return &adr{
		Seed:       a.Seed,
		PrivateKey: a.privateKey,
		PrefixPriv: a.prefixPriv,
		PrefixPub:  a.prefixPub,
		Height:     a.Height,
	}
}

func (a *Address) imports(s *adr) {
	a.Seed = s.Seed
	a.privateKey = s.PrivateKey
	a.prefixPriv = s.PrefixPriv
	a.prefixPub = s.PrefixPub
	a.Height = s.Height
}

//MarshalJSON  marshals Merkle into valid JSON.
func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.exports())
}

//UnmarshalJSON  unmarshals JSON to Merkle.
func (a *Address) UnmarshalJSON(b []byte) error {
	var s adr
	err := json.Unmarshal(b, &s)
	if err == nil {
		a.imports(&s)
	}
	return err
}

//EncodeMsgpack  marshals Address into valid JSON.
func (a *Address) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(a.exports())
}

//DecodeMsgpack  unmarshals JSON to Address.
func (a *Address) DecodeMsgpack(dec *msgpack.Decoder) error {
	var s adr
	err := dec.Decode(&s)
	if err == nil {
		a.imports(&s)
	}
	return err
}

//LeafNo returns leaf number we will use next.
func (a *Address) LeafNo() uint64 {
	return a.privateKey.LeafNo()
}

//SetLeafNo adds sets leaf no
//it returns an error if tried to set past index.
func (a *Address) SetLeafNo(n uint64) error {
	return a.privateKey.SetLeafNo(n)
}

//Sign signs msg.
func (a *Address) Sign(msg []byte) []byte {
	return a.privateKey.Sign(msg)
}

//Verify verifies msg from a node with node key..
func Verify(bsig, msg, bpk []byte) bool {
	return xmss.VerifyMT(bsig, msg, bpk)
}
