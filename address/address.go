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
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"sort"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/glyph"
	"github.com/vmihailenco/msgpack"
)

var (
	prefixAdrsString = "AKADR"
	prefixMsigString = "AKMSI"
	prefixNodeString = "AKNOD"
)

//Signature is a signature for hashed-address..
type Signature struct {
	PublicKey []byte `json:"public_key"`
	Sig       []byte `json:"sig"`
}

//Address represents an address an assciated Merkle Tree in ADK.
type Address struct {
	IsNode     bool
	privateKey *glyph.SigningKey
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

func prefixAdrs(config *aklib.Config, isNode bool) []byte {
	if isNode {
		return config.PrefixNode
	}
	return config.PrefixAdrs
}

func prefixPriv(config *aklib.Config, isNode bool) []byte {
	if isNode {
		return config.PrefixNkey
	}
	return config.PrefixPriv
}

//FromSeed makes Address struct from the seed.
func FromSeed(config *aklib.Config, seed []byte) (*Address, error) {
	isNode := false
	switch {
	case bytes.Equal(seed[:2], config.PrefixPriv):
		isNode = false
	case bytes.Equal(seed[:2], config.PrefixNkey):
		isNode = true
	default:
		return nil, errors.New("unknown seed")
	}

	sk, err := glyph.NewSigningKey(seed[2:])
	if err != nil {
		return nil, err
	}

	return &Address{
		IsNode:     isNode,
		privateKey: sk,
	}, nil
}

//New returns a new Address struct.
func New(config *aklib.Config, isNode bool) (*Address, error) {
	return NewFromSeed(config, GenerateSeed32(), isNode)
}

//NewFromSeed returns a new Address struct built from the seed.
func NewFromSeed(config *aklib.Config, seed []byte, isNode bool) (*Address, error) {
	if len(seed) != 32 {
		return nil, errors.New("invalid length of seed")
	}
	sk := glyph.NewSKFromSeed(seed)
	return &Address{
		IsNode:     isNode,
		privateKey: sk,
	}, nil
}

//PublicKey returns public key.
func (a *Address) PublicKey() []byte {
	return a.privateKey.PK().Bytes()
}

//Address returns the address in binary..
func (a *Address) Address(cfg *aklib.Config) []byte {
	pub := a.PublicKey()
	hpub := sha256.Sum256(pub)
	return append(prefixAdrs(cfg, a.IsNode), hpub[:]...)
}

//Address58 returns base58 encoded address.
func (a *Address) Address58(cfg *aklib.Config) string {
	a58, err := Address58(cfg, a.Address(cfg))
	if err != nil {
		panic(err)
	}
	return a58
}

//Address58 converts address bytes to encode58 format.
func Address58(config *aklib.Config, p []byte) (string, error) {
	var prefix string
	switch {
	case bytes.Equal(p[:2], config.PrefixAdrs):
		prefix = prefixAdrsString
	case bytes.Equal(p[:2], config.PrefixNode):
		prefix = prefixNodeString
	default:
		return "", errors.New("unknown prefix")
	}

	return prefix + Encode58(p), nil
}

//Address58ForAddress converts address bytes to encode58 format.
func Address58ForAddress(p []byte) string {
	return prefixAdrsString + Encode58(p)
}

//ParseAddress58 parses and checks base58 encoded address
//and returns binary public key and its height.
func ParseAddress58(cfg *aklib.Config, pub58 string) ([]byte, bool, error) {
	if len(pub58) <= len(prefixAdrsString) {
		return nil, false, errors.New("invalid length")
	}
	prefix := pub58[:len(prefixAdrsString)]

	pub, err := Decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, false, err
	}
	if len(pub) != 32+2 {
		return nil, false, errors.New("invalid length")
	}
	if bytes.Equal(pub[:2], cfg.PrefixAdrs) && prefix == prefixAdrsString {
		return pub, false, nil
	}
	if bytes.Equal(pub[:2], cfg.PrefixNode) && prefix == prefixNodeString {
		return pub, true, nil
	}
	return nil, false, errors.New("invalid prefix")
}

//ParseAddress58ForAddress parses and checks base58 encoded address
//and returns binary public key and its height.
func ParseAddress58ForAddress(pub58 string) ([]byte, error) {
	if len(pub58) <= len(prefixAdrsString) {
		return nil, errors.New("invalid length")
	}
	prefix := pub58[:len(prefixAdrsString)]
	if prefix != prefixAdrsString {
		return nil, errors.New("invalid prefix")
	}
	pub, err := Decode58(pub58[len(prefixAdrsString):])
	if err != nil {
		return nil, err
	}
	if len(pub) != 32+2 {
		return nil, errors.New("invalid length")
	}
	return pub, nil
}

//Address returns address bytes from sig.
func (sig *Signature) Address(cfg *aklib.Config, isNode bool) []byte {
	hadr := sha256.Sum256(sig.PublicKey)
	return append(prefixAdrs(cfg, isNode), hadr[:]...)
}

type address struct {
	IsNode     bool
	PrivateKey *glyph.SigningKey
}

func (a *Address) exports() *address {
	return &address{
		IsNode:     a.IsNode,
		PrivateKey: a.privateKey,
	}
}

func (a *Address) imports(s *address) {
	a.IsNode = s.IsNode
	a.privateKey = s.PrivateKey
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

//Sign signs msg.
func (a *Address) Sign(msg []byte) (*Signature, error) {
	sig, err := a.privateKey.Sign(msg)
	if err != nil {
		return nil, err
	}
	return &Signature{
		PublicKey: a.PublicKey(),
		Sig:       sig.Bytes(),
	}, nil
}

//Verify verifies msg from a node with node key..
func Verify(bsig *Signature, msg []byte) error {
	pk, err := glyph.NewPublickey(bsig.PublicKey)
	if err != nil {
		return err
	}
	sig, err := glyph.NewSignature(bsig.Sig)
	if err != nil {
		log.Println(err)
		return err
	}
	return pk.Verify(sig, msg)
}

//GenerateSeed32 generates a new 32 bytes seed.
func GenerateSeed32() []byte {
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		panic(err)
	}
	return seed
}

//MultisigAddress returns an multisig address.
func MultisigAddress(cfg *aklib.Config, m byte, address ...[]byte) string {
	s := sha256.New()
	s.Write([]byte{m})
	sort.Slice(address, func(i, j int) bool {
		return bytes.Compare(address[i], address[j]) < 0
	})
	for _, a := range address {
		s.Write(a)
	}
	h := append(cfg.PrefixMsig, s.Sum(nil)...)
	return prefixMsigString + Encode58(h)
}
