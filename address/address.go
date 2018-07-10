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
)

//Heights converts height byte to real height.
var Heights = []byte{2, 10, 16, 20}

//Signature is a signature for hashed-address..
type Signature struct {
	PublicKey []byte `json:"public_key"`
	Sig       []byte `json:"sig"`
}

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
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
		privateKey: xmss.NewMerkle(Heights[h], seed),
		Seed:       seed,
		prefixPub:  config.PrefixAdrs[h],
		Height:     Heights[h],
	}, nil
}

//PublicKey returns public key.
func (a *Address) PublicKey() []byte {
	return a.privateKey.PublicKey()
}

//Address returns the address in binary..
func (a *Address) Address() []byte {
	pub := a.privateKey.PublicKey()
	hpub := sha256.Sum256(pub)
	r := make([]byte, 32+3)
	copy(r, a.prefixPub)
	copy(r[3:], hpub[:])
	return r
}

//Address58 returns base58 encoded address.
func (a *Address) Address58() string {
	return To58(a.Address())
}

//To58 converts address bytes to encode58 format.
func To58(p []byte) string {
	return prefixAdrsString + Encode58(p)
}

//ParseAddress58 parses and checks base58 encoded address
//and returns binary public key and its height.
func ParseAddress58(pub58 string, cfg *aklib.Config) ([]byte, byte, error) {
	if pub58[:len(prefixAdrsString)] != prefixAdrsString {
		return nil, 0, errors.New("invalid prefix string in public key")
	}
	pub, err := Decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, 0, err
	}
	if len(pub) != 32+3 {
		return nil, 0, errors.New("invalid length")
	}
	if cfg == nil {
		return pub, 255, nil
	}
	for height := byte(Height2); height <= Height20; height++ {
		pref := cfg.PrefixAdrs[height]
		if bytes.Equal(pub[:len(pref)], pref) {
			return pub, height, nil
		}
	}
	return nil, 0, errors.New("invalid public key")
}

//Address returns address bytes from sig.
func (sig *Signature) Address(cfg *aklib.Config) ([]byte, error) {
	r := make([]byte, 32+3)
	hadr := sha256.Sum256(sig.PublicKey)
	h := 0
	switch sig.PublicKey[0] {
	case 2:
		h = 0
	case 10:
		h = 1
	case 16:
		h = 2
	case 20:
		h = 3
	default:
		return nil, errors.New("invalid public key")
	}
	copy(r, cfg.PrefixAdrs[h])
	copy(r[3:], hadr[:])
	return r, nil
}

//Index returns signature indes(idx_sig).
func (sig *Signature) Index() (uint32, error) {
	return xmss.IndexFromSig(sig.Sig)
}

type address struct {
	Seed       []byte
	PrivateKey *xmss.Merkle
	PrefixPub  []byte
	Height     byte
}

func (a *Address) exports() *address {
	return &address{
		Seed:       a.Seed,
		PrivateKey: a.privateKey,
		PrefixPub:  a.prefixPub,
		Height:     a.Height,
	}
}

func (a *Address) imports(s *address) {
	a.Seed = s.Seed
	a.privateKey = s.PrivateKey
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
func (a *Address) Sign(msg []byte) *Signature {
	return &Signature{
		PublicKey: a.PublicKey(),
		Sig:       a.privateKey.Sign(msg),
	}
}

//Verify verifies msg from a node with node key..
func Verify(bsig *Signature, msg []byte) bool {
	return xmss.Verify(bsig.Sig, msg, bsig.PublicKey)
}

//GenerateSeed generates a new 64 bytes seed.
func GenerateSeed() []byte {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return seed
}
