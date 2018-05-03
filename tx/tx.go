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

package tx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/AidosKuneen/aklib/arypack"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/cuckoo"
	sha256 "github.com/AidosKuneen/sha256-simd"

	"github.com/vmihailenco/msgpack"
)

//max length of tx and fields in a transaction.
const (
	MessageMax         = 255
	TransactionMax     = 65535
	SigMax             = 16
	OutputsMax         = 32
	MultisigMax        = 4
	MultisigAddressMax = 32
	PreviousMax        = 2
)

var (
	//TxNormal is a type for nomal tx.
	TxNormal = []byte{0xAD, 0xBF, 0x43, 0x01}
	//TxConfirm is a type for confirmation tx.
	TxConfirm = []byte{0xAD, 0xBF, 0x43, 0x02}
	//TxVote is a type for voting tx.
	TxVote = []byte{0xAD, 0xBF, 0x00, 0x03}
	//Genesis tx hash
	Genesis = make([]byte, 32)
)

//GetTXFunc gets tx.
type GetTXFunc func(hash []byte) (*Body, error)

//Bytes32 is a slice of 32 bytes
type Bytes32 [][]byte

const (
	//HashTypeExcludeOutputs is for excluding some outputs.
	HashTypeExcludeOutputs = 0x10
)

//type for minable tx.
const (
	RewardTicket = iota + 1
	RewardFee
)

func encodeMsgpack(b [][]byte, enc *msgpack.Encoder, n int) error {
	dat := make([]byte, n*len(b))
	for i, bb := range b {
		copy(dat[i*n:], bb)
	}
	return enc.Encode(dat)
}

func decodeMsgpack(dec *msgpack.Decoder, n int) ([][]byte, error) {
	var dat []byte
	if err := dec.Decode(&dat); err != nil {
		return nil, err
	}
	if len(dat)%n != 0 {
		return nil, errors.New("length of slice must be 3nN")
	}
	b := make([][]byte, len(dat)/n)
	for j := range b {
		b[j] = dat[j*n : (j+1)*n]
	}
	return b, nil
}

//EncodeMsgpack  marshals slice of 32 bytes into valid JSON.
func (b32 *Bytes32) EncodeMsgpack(enc *msgpack.Encoder) error {
	return encodeMsgpack(*b32, enc, 32)
}

//DecodeMsgpack  unmarshals msgpack bin to slice of 32 bytes.
func (b32 *Bytes32) DecodeMsgpack(dec *msgpack.Decoder) error {
	b, err := decodeMsgpack(dec, 32)
	if err != nil {
		return err
	}
	*b32 = b
	return nil
}

//Input is an input in transactions.
type Input struct {
	PreviousTX []byte `json:"previous_tx"` //32 bytes
	Index      byte   `json:"index"`
}

//Output is an output in transactions.
type Output struct {
	Address []byte `json:"address"` //65 bytes
	Value   uint64 `json:"value"`   //8 bytes
}

//MultiSigOut is an multisig output in transactions.
type MultiSigOut struct {
	N         byte    `json:"n"`         //0 means normal payment, or N out of len(Address) multisig.
	Addresses Bytes32 `json:"addresses"` //< 65 * 32 bytes
	Value     uint64  `json:"value"`     //8 bytes
}

//MultiSigIn is an multisig input in transactions.
type MultiSigIn struct {
	PreviousTX []byte `json:"previous_tx"` //32 bytes
	Index      byte   `json:"index"`
}

//max size=7021 bytes, normally 363 bytes

//Body is a Transactoin except signature.
type Body struct {
	Type         []byte         `json:"type"`          //4 bytes
	Nonce        []uint32       `json:"nonce"`         //20*VarInt(<4)
	Gnonce       uint32         `json:"g_nonce"`       //4 bytes
	Time         time.Time      `json:"time"`          //4 bytes
	Message      []byte         `json:"message"`       //<255 bytes
	Inputs       []*Input       `json:"inputs"`        //<33 * 32 = 1056 bytes
	MultiSigIns  []*MultiSigIn  `json:"multisig_ins"`  //<1032 * 4 = 4128 bytes
	Outputs      []*Output      `json:"outputs"`       //<40 * 32 = 1280 bytes
	MultiSigOuts []*MultiSigOut `json:"multisig_outs"` //<1032 * 4 = 4128 bytes
	Previous     Bytes32        `json:"previous"`      //<8*32=256 bytes
	Easiness     uint32         `json:"easiness"`      //1 byte //not used for now
	LockTime     time.Time      `json:"lock_time"`     // 4 bytes
	HashType     byte           `json:"hash_type"`
	TicketInput  []byte         `json:"ticket_input"`
	TicketOutput []byte         `json:"ticket_output"`
	Scripts      [][]byte       `json:"scripts"` //not used
	Reserved     []byte         `json:"-"`       //not used
}

//Signatures is a slice of Signature
type Signatures []*address.Signature

//Transaction is a transactio in Aidos Kuneen.
type Transaction struct {
	*Body      `json:"body"`
	Signatures `json:"signatures"`
}

//IsMinable returns reward type (fee or ticket) if tx is minable.
func (tx *Transaction) IsMinable(cfg *aklib.Config) (int, error) {
	if err := tx.check(cfg, false); err != nil {
		return 0, err
	}
	hs := tx.Hash()
	if err := cuckoo.Verify(hs, tx.Nonce); err == nil {
		return 0, errors.New("already mined")
	}
	n := tx.HashType & 0xf
	if tx.HashType&HashTypeExcludeOutputs == HashTypeExcludeOutputs && n == 1 {
		return RewardFee, nil
	}
	if tx.TicketOutput != nil {
		return RewardTicket, nil
	}
	return 0, errors.New("incorrect minable tx")
}

//Check checks the tx.
func (tx *Transaction) Check(cfg *aklib.Config) error {
	return tx.check(cfg, true)
}

func (tx *Transaction) check(cfg *aklib.Config, includePow bool) error {
	if tx.Body == nil {
		return errors.New("body is null")
	}
	if tx.Type == nil || !bytes.Equal(tx.Type, TxNormal) {
		return errors.New("invalid type")
	}
	if len(tx.Nonce) != cuckoo.ProofSize {
		return fmt.Errorf("nonce must be %d size", cuckoo.ProofSize)
	}
	if err := cuckoo.Verify(tx.hashForPoW(), tx.Nonce); includePow && err != nil {
		return err
	}
	if tx.Time.After(time.Now()) {
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
		for j := 0; j < n; j++ {
			if tx.Inputs[j].Index == i.Index && bytes.Equal(tx.Inputs[j].PreviousTX, i.PreviousTX) {
				return fmt.Errorf("input %d has a same previous and index at input %d", n, j)
			}
		}
	}
	for n, i := range tx.MultiSigIns {
		if len(i.PreviousTX) != 32 {
			return fmt.Errorf("previous tx hash at %d must be 32 bytes", n)
		}
		for j := 0; j < n; j++ {
			if tx.Inputs[j].Index == i.Index && bytes.Equal(tx.MultiSigIns[j].PreviousTX, i.PreviousTX) {
				return fmt.Errorf("input %d has a same previous and index at input %d", n, j)
			}
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
	if len(tx.MultiSigOuts) > MultisigMax {
		return fmt.Errorf("number of multisig must be under %d", MultisigMax)
	}
	for n, o := range tx.MultiSigOuts {
		if len(o.Addresses) > MultisigAddressMax {
			return fmt.Errorf("number of addresses at multisig %d must be under %d",
				n, MultisigAddressMax)
		}
		for i, a := range o.Addresses {
			if len(a) != 32 {
				return fmt.Errorf("address size at multisig %d must be 32 bytes", n)
			}
			for j := i + 1; j < len(o.Addresses); j++ {
				if bytes.Equal(a, o.Addresses[j]) {
					return fmt.Errorf("multisig %d has same address in %d and %d", n, i, j)
				}
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
	if len(tx.Previous) > PreviousMax || len(tx.Previous) == 0 {
		return fmt.Errorf("number of previous tx must be under %d and must not be 0", PreviousMax)
	}
	for n, i := range tx.Previous {
		if len(i) != 32 {
			return fmt.Errorf("tx hash size at previous tx %d must be 32 bytes", n)
		}
	}
	if tx.Easiness > cfg.Easiness {
		return fmt.Errorf("Easiness must be %d", cfg.Easiness)
	}
	if !tx.LockTime.IsZero() && tx.LockTime.After(time.Now()) {
		return errors.New("this tx is not unlocked yet")
	}
	if len(tx.Signatures) > SigMax {
		return fmt.Errorf("number of signatures must be under %d", SigMax)
	}
	if tx.HashType != 0 && tx.HashType&HashTypeExcludeOutputs != HashTypeExcludeOutputs {
		return fmt.Errorf("invalid hashtype %d", tx.HashType)
	}
	n := int(tx.HashType) & 0xf
	if tx.HashType&HashTypeExcludeOutputs == HashTypeExcludeOutputs && n >= len(tx.Outputs) {
		return fmt.Errorf("number of outputs of hashtype %d is too large", n)
	}
	if len(tx.Scripts) > 0 {
		return errors.New("cannot use scriptsd")
	}
	if len(tx.Reserved) > 0 {
		return errors.New("cannot use reserved field")
	}
	if len(tx.TicketInput) > 0 && len(tx.TicketOutput) == 0 {
		return errors.New("ticket_output is 0")
	}
	if len(tx.TicketInput) == 0 && len(tx.TicketOutput) > 0 {
		//ticket
		if len(tx.Inputs) > 0 || len(tx.MultiSigIns) > 0 || len(tx.Outputs) > 0 || len(tx.MultiSigOuts) > 0 || !tx.LockTime.IsZero() ||
			tx.HashType != 0 {
			return errors.New("tx content for ticket must be empty")
		}
		if tx.Easiness > cfg.TicketEasiness {
			return errors.New("PoW doesn't meet ticket difficulty")
		}
	}
	dat := tx.BytesForSign()
	for n, sig := range tx.Signatures {
		if !address.Verify(sig, dat) {
			return fmt.Errorf("failed to verify a signature at %d", n)
		}
		for nn := n + 1; nn < len(tx.Signatures); nn++ {
			if bytes.Equal(sig.PublicKey, tx.Signatures[nn].PublicKey) {
				return fmt.Errorf("there are same publik keys in signature at %d and %d", n, nn)
			}
		}
	}
	return tx.hasValidHashes(cfg, includePow)
}

//isValidHash reteurns true if  hash bytes h meets difficulty.
func isValidHash(h []byte, e uint32) bool {
	ea := binary.LittleEndian.Uint32(h[len(h)-4:])
	return ea <= e
}

//hasValidHashes reteurns true if  hashes in tx and tx hash  meets difficulty.
func (tx *Transaction) hasValidHashes(cfg *aklib.Config, includePow bool) error {
	h := tx.Hash()
	if !isValidHash(h, tx.Easiness) && includePow {
		return errors.New("tx hash doesn't not match easiness")
	}
	for _, i := range tx.Inputs {
		if !isValidHash(i.PreviousTX, cfg.Easiness) {
			return errors.New("inputs txs' hash doesn't not match easiness")
		}
	}
	for _, p := range tx.Previous {
		if !isValidHash(p, cfg.Easiness) {
			return errors.New("previous txs' hash doesn't not match easiness")
		}
	}
	return nil
}

//Hash reteurns hash of tx.
func (tx *Transaction) Hash() []byte {
	hh := sha256.Sum256(arypack.Marshal(tx))
	return hh[:]
}

//NoExistHashes returns tx hashes which are not found.
//getTx must return err if tx is not found.
func (tx *Transaction) NoExistHashes(getTX GetTXFunc) [][]byte {
	hs := make([][]byte, 0, len(tx.Previous)+len(tx.Inputs))
	for _, i := range tx.Previous {
		if _, err := getTX(i); err != nil {
			hs = append(hs, i)
		}
	}
	for _, i := range tx.Inputs {
		if _, err := getTX(i.PreviousTX); err != nil {
			hs = append(hs, i.PreviousTX)
		}
	}
	for _, i := range tx.MultiSigIns {
		if _, err := getTX(i.PreviousTX); err != nil {
			hs = append(hs, i.PreviousTX)
		}
	}
	if tx.TicketInput != nil {
		hs = append(hs, tx.TicketInput)
	}
	return hs
}

//CheckAll checks the tx, including other txs refered by the tx..
//Genesis block must be saved in the store
func (tx *Transaction) CheckAll(getTX GetTXFunc, cfg *aklib.Config) error {
	return tx.checkAll(getTX, cfg, true)
}

type addresses struct {
	used bool
	adr  []byte
}

func hasAddress(adrs []*addresses, adr []byte) bool {
	for _, a := range adrs {
		if bytes.Equal(a.adr, adr) {
			a.used = true
			return true
		}
	}
	return false
}

func hasUunused(adrs []*addresses) bool {
	for _, a := range adrs {
		if !a.used {
			return true
		}
	}
	return false
}

func (tx *Transaction) checkAll(getTX GetTXFunc, cfg *aklib.Config, includePow bool) error {
	if err := tx.check(cfg, includePow); err != nil {
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
	for _, o := range tx.MultiSigOuts {
		totalout += o.Value
	}
	var totalin uint64
	adrs := make([]*addresses, 0, len(tx.Inputs)+1)
	for _, sig := range tx.Signatures {
		hadr := sha256.Sum256(sig.PublicKey)
		adrs = append(adrs, &addresses{
			adr: hadr[:],
		})
	}
	for n, inp := range tx.Inputs {
		inTX, err := getTX(inp.PreviousTX)
		if err != nil {
			return err
		}
		if len(inTX.Outputs) <= int(inp.Index) {
			return fmt.Errorf("invalid input index, should be under  %d", len(inTX.Outputs))
		}
		totalin += inTX.Outputs[inp.Index].Value
		inTXAdr := inTX.Outputs[inp.Index].Address
		if !hasAddress(adrs, inTXAdr) {
			return fmt.Errorf("no signature for input %d", n)
		}
	}
	for n, inp := range tx.MultiSigIns {
		inTX, err := getTX(inp.PreviousTX)
		if err != nil {
			return err
		}
		if len(inTX.MultiSigOuts) <= int(inp.Index) {
			return fmt.Errorf("invalid multisig index, should be under  %d", len(inTX.MultiSigOuts))
		}
		mul := inTX.MultiSigOuts[inp.Index]
		totalin += mul.Value
		exist := 0
		for _, adr := range mul.Addresses {
			if hasAddress(adrs, adr) {
				exist++
			}
		}
		if exist != int(mul.N) {
			return fmt.Errorf("invalid number of valid signatures %d in multisig %d, should be %d", exist, n, mul.N)
		}
	}
	if len(tx.TicketInput) > 0 {
		inTX, err := getTX(tx.TicketInput)
		if err != nil {
			return err
		}
		if !hasAddress(adrs, inTX.TicketOutput) {
			return errors.New("cannot verify the ticket input")
		}
	}
	if hasUunused(adrs) {
		return errors.New("there are(is) unsed signature")
	}
	if totalin != totalout {
		return fmt.Errorf("total input ADK %v does not equal to one of output %v",
			totalin, totalout)
	}
	return nil
}

//BytesForSign returns a hash slice for  signinig
func (tx *Transaction) BytesForSign() []byte {
	return tx.partialbytes(true)
}

//PreHash returns a hash before PoW.
func (tx *Transaction) PreHash() []byte {
	bytes := tx.partialbytes(false)
	hs := sha256.Sum256(bytes)
	return hs[:]
}

func (tx *Transaction) partialbytes(isBodyOnly bool) []byte {
	tx2 := tx.Clone()
	tx2.Gnonce = 0
	tx2.Nonce = nil
	tx2.TicketOutput = nil
	if tx.HashType&0xf0 == HashTypeExcludeOutputs {
		exclude := int(tx.HashType & 0x0f)
		for i := 0; i < exclude; i++ {
			tx2.Outputs[len(tx.Outputs)-1-i].Address = nil
			tx2.Outputs[len(tx.Outputs)-1-i].Value = 0
		}
	}
	if isBodyOnly {
		return arypack.Marshal(tx2.Body)
	}
	return arypack.Marshal(tx2)
}

func (tx *Transaction) hashForPoW() []byte {
	nonce := tx.Nonce
	tx.Nonce = nil
	btx := arypack.Marshal(tx)
	h := sha256.Sum256(btx)
	tx.Nonce = nonce
	return h[:]
}

//Clone clones tx.
func (tx *Transaction) Clone() *Transaction {
	var tx2 Transaction
	if err := arypack.Unmarshal(arypack.Marshal(tx), &tx2); err != nil {
		panic(err)
	}
	return &tx2
}
