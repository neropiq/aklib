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

package tx

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/AidosKuneen/cuckoo"

	"github.com/AidosKuneen/aklib"
	sha256 "github.com/AidosKuneen/sha256-simd"
)

//max length of tx and fields in a transaction.
const (
	MessageMax         = 255
	TransactionMax     = 65535
	SigMax             = 16
	OutputsMax         = 32
	MultisigMax        = 4
	MultisigAddressMax = 32
	PreviousMax        = 8

	nonceLocation = 4
)

var (
	txType = []byte{0xAD, 0xBF, 0x00, 0x01}
	//Genesis tx hash
	Genesis = make([]byte, 32)
)

//GetTXFunc gets tx.
type GetTXFunc func(hash []byte) (*Body, error)

//VerifyFunc verifies bsig with message bms and pubkey bpk..
type VerifyFunc func(bsig []byte, msg []byte, bpk []byte) bool

//Input is an input in transactions.
type Input struct {
	PreviousTX []byte //32 bytes
	Index      byte
}

//Output is an output in transactions.
type Output struct {
	Address []byte //32 bytes
	Value   uint64 //8 bytes
}

//MultiSig is an multisig output in transactions.
type MultiSig struct {
	N         byte     //0 means normal payment, or N out of len(Address) multisig.
	Addresses [][]byte //< 32 * 32 bytes
	Value     uint64   //8 bytes
}

//max size=7021 bytes, normally 363 bytes

//Body is a Transactoin except signature.
type Body struct {
	Type       []byte      //4 bytes
	Nonce      []uint32    //20*VarInt
	Time       uint32      //4 bytes
	Message    []byte      //<255 bytes
	Inputs     []*Input    //<33 * 32 = 1056 bytes
	Outputs    []*Output   //<40 * 32 = 1280 bytes
	MultiSigs  []*MultiSig //<1032 * 4 = 4128 bytes
	Previous   [][]byte    //<8*32=256 bytes
	Difficulty byte        //1 byte //not used for now
	LockTime   uint32      // 4 bytes
}

//Signatures is a signature  part of Transaction.
type Signatures [][]byte //  2852 * N bytes. < 2852*16=.45632 bytes

//Transaction is a transactio in Aidos Kuneen.
type Transaction struct {
	*Body
	Signatures
}

//New returns Transaction struct filled in type field.
func New(cfg *aklib.Config) *Transaction {
	tx := &Transaction{
		Body: &Body{
			Type:       txType,
			Difficulty: cfg.Difficulty,
		},
	}
	return tx
}

//Check checks the tx.
func (tx *Transaction) Check(cfg *aklib.Config) error {
	if tx.Body == nil || tx.Signatures == nil {
		return errors.New("body or signature is null")
	}
	if tx.Type == nil || !bytes.Equal(tx.Type, txType) {
		return errors.New("invalid type")
	}
	if len(tx.Nonce) != cuckoo.ProofSize {
		return fmt.Errorf("nonce must be %d size", cuckoo.ProofSize)
	}
	bs := tx.bytesForPoW()
	hs := sha256.Sum256(bs)
	var nonces [cuckoo.ProofSize]uint32
	copy(nonces[:], tx.Nonce)
	if err := cuckoo.Verify(hs[:], &nonces); err != nil {
		return err
	}
	if time.Unix(int64(tx.Time), 0).After(time.Now()) {
		return errors.New("timestamp is in future")
	}
	if len(tx.Message) > MessageMax {
		return fmt.Errorf("message length must be under %d bytes", MessageMax)
	}
	if len(tx.Inputs) > SigMax {
		return fmt.Errorf("number of inputs must be under %d", SigMax)
	}
	for n, i := range tx.Inputs {
		if len(i.PreviousTX) != 32 {
			return fmt.Errorf("previous tx hash at %d must be 32 bytes", n)
		}
	}
	if len(tx.Outputs) > OutputsMax {
		return fmt.Errorf("number of output must be %d", OutputsMax)
	}
	for n, o := range tx.Outputs {
		if len(o.Address) != 32 {
			return fmt.Errorf("address in outputs %d must be 32 bytes", n)
		}
		if o.Value > aklib.ADKSupply {
			return fmt.Errorf("value in outputs %d must be under %d adk",
				n, aklib.ADKSupply)
		}
	}
	if len(tx.MultiSigs) > MultisigMax {
		return fmt.Errorf("number of multisig must be under %d", MultisigMax)
	}
	for n, o := range tx.MultiSigs {
		if len(o.Addresses) > MultisigAddressMax {
			return fmt.Errorf("number of addresses at multisig %d must be under %d",
				n, MultisigAddressMax)
		}
		for _, a := range o.Addresses {
			if len(a) != 32 {
				return fmt.Errorf("address size at multisig %d must be 32 bytes", n)
			}
		}
		if o.N > byte(len(o.Addresses)) {
			return fmt.Errorf("N at multisig %d must be under number of address %d",
				n, len(o.Addresses))
		}
		if o.Value > aklib.ADKSupply {
			return fmt.Errorf("value in multisig %d must be under %d adk",
				n, aklib.ADKSupply)
		}
	}
	if len(tx.Previous) > PreviousMax {
		return fmt.Errorf("number of previous tx must be under %d", PreviousMax)
	}
	for n, i := range tx.Previous {
		if len(i) != 32 {
			return fmt.Errorf("tx hash size at previous tx %d must be 32 bytes", n)
		}
	}
	if tx.Difficulty < cfg.Difficulty {
		return fmt.Errorf("difficulty must be %d", cfg.Difficulty)
	}
	if tx.LockTime != 0 && time.Unix(int64(tx.LockTime), 0).After(time.Now()) {
		return errors.New("this tx is not unlocked yet")
	}
	if len(tx.Signatures) < len(tx.Inputs) {
		return fmt.Errorf("number signature must be over the total one of input %d",
			len(tx.Inputs))
	}
	if len(tx.Signatures) > SigMax {
		return fmt.Errorf("number of signatures must be under %d", SigMax)
	}
	return tx.hasValidHashes(cfg)
}

//isValidHash reteurns true if  hash bytes h meets difficulty.
func isValidHash(h []byte, dif byte) bool {
	var i byte
	for i = 0; i < dif>>3; i++ {
		if h[31-i] != 0x00 {
			return false
		}
	}
	d := dif - (i << 3)
	if d == 0 {
		return true
	}
	b := (1 << (8 - d)) - 1
	if h[31-i] > byte(b) {
		return false
	}
	return true
}

//hasValidHashes reteurns true if  hashes in tx and tx hash  meets difficulty.
func (tx *Transaction) hasValidHashes(cfg *aklib.Config) error {
	h := tx.Hash()
	if !isValidHash(h, tx.Difficulty) {
		return errors.New("tx hash doesn't not match difficulty")
	}
	for _, i := range tx.Inputs {
		if !isValidHash(i.PreviousTX, cfg.Difficulty) {
			return errors.New("inputs txs' hash doesn't not match difficulty")
		}
	}
	for _, p := range tx.Previous {
		if !isValidHash(p, cfg.Difficulty) {
			return errors.New("previous txs' hash doesn't not match difficulty")
		}
	}
	return nil
}

//Hash reteurns hash of tx.
func (tx *Transaction) Hash() []byte {
	bd := tx.Body.Pack()
	sig2 := tx.Signatures.Pack()
	h := sha256.Sum256(sig2)
	bd = append(bd, h[:]...)
	hh := sha256.Sum256(bd)
	return hh[:]
}

//NoExistHashes returns tx hashes which is not found.
//getTx must return er if tx is not found.
func (tx *Transaction) NoExistHashes(getTX GetTXFunc, errNotFound error) [][]byte {
	hs := make([][]byte, 0, len(tx.Previous)+len(tx.Inputs))
	for _, i := range tx.Previous {
		if _, err := getTX(i); err == errNotFound {
			hs = append(hs, i)
		}
	}
	for _, i := range tx.Inputs {
		if _, err := getTX(i.PreviousTX); err == errNotFound {
			hs = append(hs, i.PreviousTX)
		}
	}
	return hs
}

//CheckAll checks the tx, including other txs refered by the tx..
//Genesis block must be saved in the store.
func (tx *Transaction) CheckAll(getTX GetTXFunc, verify VerifyFunc,
	cfg *aklib.Config) error {
	if err := tx.Check(cfg); err != nil {
		return err
	}
	for _, i := range tx.Previous {
		if _, err := getTX(i); err != nil {
			return err
		}
	}
	var totalout uint64
	for _, o := range tx.Outputs {
		totalout += o.Value
	}
	for _, o := range tx.MultiSigs {
		totalout += o.Value
	}
	dat := tx.BytesForSign()
	var totalin uint64
	var nsig byte
	for n, inp := range tx.Inputs {
		inTX, err := getTX(inp.PreviousTX)
		if err != nil {
			return err
		}
		if byte(len(inTX.Outputs)+len(inTX.MultiSigs)) <= inp.Index {
			return fmt.Errorf("index at input %d must be under number of previous transaction outputs %d",
				n, len(inTX.Outputs)+len(inTX.MultiSigs))
		}
		if byte(len(inTX.Outputs)) > inp.Index {
			totalin += inTX.Outputs[inp.Index].Value
			if !verify(tx.Signatures[nsig], dat,
				inTX.Outputs[inp.Index].Address) {
				return fmt.Errorf("failed to verify a signature at %d", n)
			}
			nsig++
		} else {
			mul := inTX.MultiSigs[inp.Index-byte(len(inTX.Outputs))]
			totalin += mul.Value
			j := 0
		loop:
			for i := byte(0); i < mul.N; i++ {
				for ; j < len(mul.Addresses); j++ {
					if verify(tx.Signatures[nsig], dat, mul.Addresses[j]) {
						nsig++
						continue loop
					}
				}
				return fmt.Errorf("cannot verify a multisig signature at input %d", n)
			}
		}
	}
	if int(nsig) != len(tx.Signatures) {
		return fmt.Errorf("number of signatures %d must be same as one of input %d",
			len(tx.Signatures), nsig)
	}
	if totalin != totalout {
		return fmt.Errorf("total input ADK %v does not equal to one of output %v",
			totalin, totalout)
	}
	return nil
}

//BytesForSign returns byte slice for  signinig
func (tx *Transaction) BytesForSign() []byte {
	bd2 := *(tx.Body)
	bd2.Nonce = make([]uint32, cuckoo.ProofSize)
	return bd2.Pack()
}
func (tx *Transaction) bytesForPoW() []byte {
	bd2 := *(tx.Body)
	bd2.Nonce = make([]uint32, cuckoo.ProofSize)
	bd := bd2.Pack()
	sig2 := tx.Signatures.Pack()
	h := sha256.Sum256(sig2)
	bd = append(bd, h[:]...)
	return bd
}

//Clone clones tx.
func (tx *Transaction) Clone() *Transaction {
	var err error
	tx2 := &Transaction{}
	tx2.Body, err = UnpackBody(tx.Body.Pack())
	if err != nil {
		panic(err)
	}
	tx2.Signatures, err = UnpackSignature(tx.Signatures.Pack())
	if err != nil {
		panic(err)
	}
	return tx2
}
