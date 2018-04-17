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
	"github.com/vmihailenco/msgpack"
)

//Height represents height of Merkle Tree for XMSS
const (
	Height2 = iota
	Height10
	Height16
	Height20
)

var (
	prefixAdrsString = "AKADR"
	prefixPrivString = "AKPRI"
)

var heights = []byte{2, 10, 16, 20}

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
	prefixPriv []byte
	prefixPub  []byte
	privateKey *xmss.Merkle
	Seed       []byte
	Height     byte
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
		privateKey: xmss.NewMerkle(heights[h], seed),
		Seed:       seed,
		prefixPriv: config.PrefixPriv[h],
		prefixPub:  config.PrefixAdrs[h],
		Height:     heights[h],
	}, nil
}

//Seed58 returns base58-encoded encrypted seed..
func (a *Address) Seed58(pwd []byte) string {
	out := make([]byte, len(a.Seed)+4)
	copy(out, a.Seed)
	hash := sha256.Sum256(a.Seed)
	hash = sha256.Sum256(hash[:])

	copy(out[len(a.Seed):], hash[0:4])
	eseed := enc(out, pwd)

	pref := a.prefixPriv
	s := make([]byte, len(eseed)+len(pref))
	copy(s, pref)
	copy(s[len(pref):], eseed)
	return prefixPrivString + Encode58(s)
}

func from58(seed58 string, pwd []byte, cfg *aklib.Config) (byte, []byte, error) {
	if seed58[:len(prefixPrivString)] != prefixPrivString {
		return 0, nil, errors.New("invalid prefix string in seed")
	}
	eseed, err := Decode58(seed58[len(prefixPrivString):])
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
	return a.privateKey.PublicKey()
}

//Address returns the address in binary..
func (a *Address) Address() []byte {
	pub := a.privateKey.PublicKey()
	hpub := sha256.Sum256(pub)
	return hpub[:]
}

//Address58 returns base58 encoded address.
func (a *Address) Address58() string {
	pref := a.prefixPub
	hpub := a.Address()
	p := make([]byte, len(hpub)+len(pref))
	copy(p, pref)
	copy(p[len(pref):], hpub)
	return prefixAdrsString + Encode58(p)
}

//FromAddress58 returns decode public key from base58 encoded string.
func FromAddress58(pub58 string, cfg *aklib.Config) ([]byte, byte, error) {
	if pub58[:len(prefixAdrsString)] != prefixAdrsString {
		return nil, 0, errors.New("invalid prefix string in public key")
	}
	pub, err := Decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, 0, err
	}
	var height byte
	for ; height < Height20; height++ {
		pref := cfg.PrefixAdrs[height]
		if bytes.Equal(pub[:len(pref)], pref) {
			break
		}
	}
	if height > Height20 {
		return nil, 0, errors.New("invalid prefix bytes in seed")
	}
	pub = pub[len(cfg.PrefixAdrs[height]):]
	return pub, heights[height], nil
}

type address struct {
	Seed       []byte
	PrivateKey *xmss.Merkle
	PrefixPriv []byte
	PrefixPub  []byte
	Height     byte
}

func (a *Address) exports() *address {
	return &address{
		Seed:       a.Seed,
		PrivateKey: a.privateKey,
		PrefixPriv: a.prefixPriv,
		PrefixPub:  a.prefixPub,
		Height:     a.Height,
	}
}

func (a *Address) imports(s *address) {
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
	var s address
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
	var s address
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
//xmsss  should refreshe auth and stack in MerkleTree..
//it returns an error if tried to set past index.
func (a *Address) SetLeafNo(n uint64) error {
	return a.privateKey.SetLeafNo(n)
}

//Sign signs msg.
func (a *Address) Sign(msg []byte) []byte {
	return a.privateKey.Sign(msg)
}

//GenerateSeed generates a new 64 bytes seed.
func GenerateSeed() []byte {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return seed
}
