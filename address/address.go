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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"

	"github.com/AidosKuneen/aklib"
	sha256 "github.com/AidosKuneen/sha256-simd"
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
	config *aklib.Config
	merkle *xmss.Merkle
	Seed   []byte
}

func enc(text []byte, pwd []byte) []byte {
	key := sha256.Sum256(pwd)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	iv := sha256.Sum256(key[:])
	ctext := make([]byte, len(text))
	encryptStream := cipher.NewCTR(block, iv[:aes.BlockSize])
	encryptStream.XORKeyStream(ctext, text)
	return ctext
}

//New returns Address struct.
func New(h byte, seed []byte, config *aklib.Config) (*Address, error) {
	if h > Height20 {
		return nil, errors.New("invalid height")
	}
	return &Address{
		merkle: xmss.NewMerkle(heights[h], seed),
		Seed:   seed,
		config: config,
	}, nil
}

//NewFromEncrypted returns Address struct.
func NewFromEncrypted(h byte, eseed, pwd []byte, config *aklib.Config) (*Address, error) {
	return New(h, enc(eseed, pwd), config)
}

//Seed58 returns base58-encoded encrypted seed.
func (a *Address) Seed58(pwd []byte) string {
	out := make([]byte, len(a.Seed)+4)
	copy(out, a.Seed)
	hash := sha256.Sum256(a.Seed)
	hash = sha256.Sum256(hash[:])

	copy(out[len(a.Seed):], hash[0:4])
	eseed := enc(out, pwd)

	pref := a.config.PrefixPriv[a.Height()]
	s := make([]byte, len(eseed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], eseed)
	return prefixPrivString + encode58(s)
}

func from58(seed58 string, pwd []byte, cfg *aklib.Config) (byte, []byte, error) {
	if prefixPrivString != seed58[:len(prefixPrivString)] {
		return 0, nil, errors.New("invalid prefix string in seed")
	}
	eseed, err := decode58(seed58[len(prefixPrivString):])
	if err != nil {
		return 0, nil, err
	}
	var height byte
	for ; height <= Height20; height++ {
		pref := cfg.PrefixPriv[height]
		if bytes.Equal(eseed[:len(pref)], pref) {
			break
		}
	}
	if height > Height20 {
		return 0, nil, errors.New("invalid prefix bytes in seed")
	}
	eseed = eseed[len(cfg.PrefixPriv[height]):]
	seed := enc(eseed, pwd)
	encoded := seed[:len(seed)-4]
	cksum := seed[len(seed)-4:]

	//Perform SHA-256 twice
	hash := sha256.Sum256(encoded)
	hash = sha256.Sum256(hash[:])
	if !bytes.Equal(hash[:4], cksum) {
		return 0, nil, errors.New("invalid password")
	}
	return height, encoded, nil

}

//NewFrom58 returns Address struct with base58 encoded seed.
func NewFrom58(seed58 string, pwd []byte, cfg *aklib.Config) (*Address, error) {
	height, seed, err := from58(seed58, pwd, cfg)
	if err != nil {
		return nil, err
	}
	return New(height, seed, cfg)
}

//PublicKey returns public key.
func (a *Address) PublicKey() []byte {
	return a.merkle.PublicKey()
}

//Height returns height of Merkle Tree..
func (a *Address) Height() byte {
	switch a.merkle.Height {
	case 10:
		return Height10
	case 16:
		return Height16
	case 20:
		return Height20
	}
	return 0
}

//PK58 returns base58 encoded public key.
func (a *Address) PK58() string {
	pref := a.config.PrefixAdrs[a.Height()]
	pub := a.merkle.PublicKey()
	p := make([]byte, len(pub)+len(pref))
	copy(p, pref)
	copy(p[len(pref):], pub)
	return prefixAdrsString + encode58(p)
}

//FromPK58 returns decode public key from base58 encoded string.
func FromPK58(pub58 string, cfg *aklib.Config) ([]byte, error) {
	if prefixAdrsString != pub58[:len(prefixAdrsString)] {
		return nil, errors.New("invalid prefix string in public key")
	}
	pub, err := decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, err
	}
	var height byte
	for ; height <= Height20; height++ {
		pref := cfg.PrefixAdrs[height]
		if bytes.Equal(pub[:len(pref)], pref) {
			break
		}
	}
	if height > Height20 {
		return nil, errors.New("invalid prefix bytes in seed")
	}
	return pub[len(cfg.PrefixAdrs[height]):], nil
}

//MarshalJSON  marshals Address into valid JSON.
func (a *Address) MarshalJSON() ([]byte, error) {
	s := struct {
		Seed   []byte
		Merkle *xmss.Merkle
		Config *aklib.Config
	}{
		Seed:   a.Seed,
		Merkle: a.merkle,
		Config: a.config,
	}
	return json.Marshal(&s)
}

//UnmarshalJSON  unmarshals JSON to Address.
func (a *Address) UnmarshalJSON(b []byte) error {
	s := struct {
		Seed   []byte
		Merkle *xmss.Merkle
		Config *aklib.Config
	}{}
	err := json.Unmarshal(b, &s)
	a.Seed = s.Seed
	a.merkle = s.Merkle
	a.config = s.Config
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
