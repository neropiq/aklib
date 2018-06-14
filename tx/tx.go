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
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/aklib/arypack"
	"github.com/AidosKuneen/cuckoo"

	"github.com/vmihailenco/msgpack"
)

//max length of tx and fields in a transaction.
const (
	MessageMax     = 255
	TransactionMax = 2000000
	ArrayMax       = 255
)

var (
	//TypeNormal is a type for nomal tx.
	typeNormal = []byte{0xAD, 0xBF, 0x43, 0x01}
	//TypeConfirm is a type for confirmation tx.
	typeConfirm = []byte{0xAD, 0xBF, 0x43, 0x02}
	//TypeVote is a type for voting tx.
	typeVote = []byte{0xAD, 0xBF, 0x00, 0x03}
)

//Hash is a tx hash.
type Hash []byte

//Array converts slice to array
func (h Hash) Array() [32]byte {
	var h32 [32]byte
	copy(h32[:], h)
	return h32
}

//GetTXFunc gets tx.
type GetTXFunc func(hash []byte) (*Body, error)

//HashSlice is a slice of 32 bytes
type HashSlice []Hash

//Address  is an address of xmss.
type Address []byte

//AddressSlice is a slice of 32 bytes
type AddressSlice HashSlice

//Types when hashing a tx.
const (
	//HashTypeExcludeOutputs is for excluding some outputs.
	HashTypeNormal           = 0x0
	HashTypeExcludeOutputs   = 0x10
	HashTypeExcludeTicketOut = 0x20
)

//Type is a tx type.
type Type byte

//type for minable tx.
const (
	TxNormal Type = iota
	TxRewardTicket
	TxRewardFee
)

func encodeMsgpack(b []Hash, enc *msgpack.Encoder, n int) error {
	dat := make([]byte, n*len(b))
	for i, bb := range b {
		copy(dat[i*n:], bb)
	}
	return enc.Encode(dat)
}

func decodeMsgpack(dec *msgpack.Decoder, n int) ([]Hash, error) {
	var dat []byte
	if err := dec.Decode(&dat); err != nil {
		return nil, err
	}
	if len(dat)%n != 0 {
		return nil, errors.New("length of slice must be 3nN")
	}
	b := make([]Hash, len(dat)/n)
	for j := range b {
		b[j] = dat[j*n : (j+1)*n]
	}
	return b, nil
}

//EncodeMsgpack  marshals slice of 32 bytes into valid JSON.
func (b32 *HashSlice) EncodeMsgpack(enc *msgpack.Encoder) error {
	return encodeMsgpack(*b32, enc, 32)
}

//DecodeMsgpack  unmarshals msgpack bin to slice of 32 bytes.
func (b32 *HashSlice) DecodeMsgpack(dec *msgpack.Decoder) error {
	b, err := decodeMsgpack(dec, 32)
	if err != nil {
		return err
	}
	*b32 = b
	return nil
}

//EncodeMsgpack  marshals slice of 32 bytes into valid JSON.
func (b32 *AddressSlice) EncodeMsgpack(enc *msgpack.Encoder) error {
	return (*HashSlice)(b32).EncodeMsgpack(enc)
}

//DecodeMsgpack  unmarshals msgpack bin to slice of 32 bytes.
func (b32 *AddressSlice) DecodeMsgpack(dec *msgpack.Decoder) error {
	return (*HashSlice)(b32).DecodeMsgpack(dec)
}

//Input is an input in transactions.
type Input struct {
	PreviousTX Hash `json:"previous_tx"` //32 bytes
	Index      byte `json:"index"`
}

//Output is an output in transactions.
type Output struct {
	Address Address `json:"address"` //65 bytes
	Value   uint64  `json:"value"`   //8 bytes
}

//MultiSigOut is an multisig output in transactions.
type MultiSigOut struct {
	N         byte         `json:"n"`         //0 means normal payment, or N out of len(Address) multisig.
	Addresses AddressSlice `json:"addresses"` //< 65 * 32 bytes
	Value     uint64       `json:"value"`     //8 bytes
}

//MultiSigIn is an multisig input in transactions.
type MultiSigIn struct {
	PreviousTX Hash `json:"previous_tx"` //32 bytes
	Index      byte `json:"index"`
}

//max size=7021 bytes, normally 363 bytes

//Body is a Transactoin except signature.
type Body struct {
	Type         []byte         `json:"type"`                    //4 bytes
	Nonce        []uint32       `json:"nonce"`                   //20*VarInt(<4)
	Gnonce       uint32         `json:"g_nonce"`                 //4 bytes
	Time         time.Time      `json:"time"`                    //4 bytes
	Message      []byte         `json:"message,omitempty"`       //<255 bytes
	Inputs       []*Input       `json:"inputs,omitempty"`        //<33 * 32 = 1056 bytes
	MultiSigIns  []*MultiSigIn  `json:"multisig_ins,omitempty"`  //<1032 * 4 = 4128 bytes
	Outputs      []*Output      `json:"outputs,omitempty"`       //<40 * 32 = 1280 bytes
	MultiSigOuts []*MultiSigOut `json:"multisig_outs,omitempty"` //<1032 * 4 = 4128 bytes
	Previous     HashSlice      `json:"previous"`                //<8*32=256 bytes
	Easiness     uint32         `json:"easiness"`                //1 byte //not used for now
	LockTime     time.Time      `json:"lock_time"`               // 4 bytes
	HashType     uint16         `json:"hash_type"`
	TicketInput  Hash           `json:"ticket_input,omitempty"`
	TicketOutput Address        `json:"ticket_output,omitempty"`
	Scripts      [][]byte       `json:"scripts,omitempty"` //not used
	Reserved     []byte         `json:"-"`                 //not used
}

//Signatures is a slice of Signature
type Signatures []*address.Signature

//Transaction is a transactio in Aidos Kuneen.
type Transaction struct {
	*Body      `json:"body"`
	Signatures `json:"signatures"`
}

//New returns a transaction object.
func New(s *aklib.Config, previous ...Hash) *Transaction {
	return &Transaction{
		Body: &Body{
			Type:     typeNormal,
			Time:     time.Now(),
			Easiness: s.Easiness,
			Previous: previous,
		},
	}
}

//IssueTicket make and does PoW for a  transaction for issuing tx.
func IssueTicket(s *aklib.Config, ticketOut *address.Address, previous ...Hash) (*Transaction, error) {
	tr := &Transaction{
		Body: &Body{
			Type:         typeNormal,
			Time:         time.Now(),
			Easiness:     s.TicketEasiness,
			TicketOutput: ticketOut.Address(),
			Previous:     previous,
		},
	}
	return tr, tr.PoW()
}

//NewMinableFee returns a minable transaction by fee..
func NewMinableFee(s *aklib.Config, previous ...Hash) *Transaction {
	return &Transaction{
		Body: &Body{
			Type:     typeNormal,
			Time:     time.Now(),
			Easiness: s.Easiness,
			HashType: HashTypeExcludeOutputs | 0x1,
			Previous: previous,
		},
	}
}

//NewMinableTicket returns a minable transaction by ticket..
func NewMinableTicket(s *aklib.Config, ticketIn Hash, previous ...Hash) *Transaction {
	return &Transaction{
		Body: &Body{
			Type:        typeNormal,
			Time:        time.Now(),
			Easiness:    s.Easiness,
			HashType:    HashTypeExcludeTicketOut,
			Previous:    previous,
			TicketInput: ticketIn,
		},
	}
}

//AddInput add an input into tx.
func (body *Body) AddInput(h Hash, idx byte) {
	body.Inputs = append(body.Inputs, &Input{
		PreviousTX: h,
		Index:      idx,
	})
}

//AddOutput add an output into tx.
func (body *Body) AddOutput(adr []byte, v uint64) {
	body.Outputs = append(body.Outputs, &Output{
		Address: adr,
		Value:   v,
	})
}

//AddMultisigIn add a multisig input into tx.
func (body *Body) AddMultisigIn(h Hash, idx byte) {
	body.MultiSigIns = append(body.MultiSigIns, &MultiSigIn{
		PreviousTX: h,
		Index:      idx,
	})
}

//AddMultisigOut add a mulsig output into tx.
func (body *Body) AddMultisigOut(n byte, v uint64, adrs ...[]byte) error {
	if len(adrs) < int(n) {
		return errors.New("length of adrs is less than n")
	}
	as := make(AddressSlice, len(adrs))
	for i, adr := range adrs {
		as[i] = adr
	}
	body.MultiSigOuts = append(body.MultiSigOuts, &MultiSigOut{
		N:         n,
		Addresses: as,
		Value:     v,
	})
	return nil
}

//Sign sings the tx.
func (tr *Transaction) Sign(a *address.Address) error {
	dat, err := tr.bytesForSign()
	if err != nil {
		return err
	}
	tr.Signatures = append(tr.Signatures, a.Sign(dat))
	return nil
}

//Signature returns  singture of the tx.
func (tr *Transaction) Signature(a *address.Address) (*address.Signature, error) {
	dat, err := tr.bytesForSign()
	if err != nil {
		return nil, err
	}
	return a.Sign(dat), nil
}

//AddSig adds a signature  to tx.
func (tr *Transaction) AddSig(sig *address.Signature) {
	tr.Signatures = append(tr.Signatures, sig)
}

//Size returns tx size.
func (tr *Transaction) Size() int {
	return len(arypack.Marshal(tr))
}

//Check checks the tx.
func (tr *Transaction) Check(cfg *aklib.Config, typ Type) error {
	if typ != TxNormal &&
		typ != TxRewardFee && typ != TxRewardTicket {
		return errors.New("invalid reward type")
	}
	powed := true
	if typ == TxRewardFee || typ == TxRewardTicket {
		powed = false
	}
	if tr.Size() > TransactionMax {
		return errors.New("tx size is too big")
	}
	if tr.Body == nil {
		return errors.New("body is null")
	}
	if !bytes.Equal(tr.Type, typeNormal) {
		return errors.New("invalid type")
	}
	switch powed {
	case true:
		if len(tr.Nonce) != cuckoo.ProofSize {
			return fmt.Errorf("nonce must be %d size", cuckoo.ProofSize)
		}
		if err := cuckoo.Verify(tr.hashForPoW(), tr.Nonce); err != nil {
			return err
		}
	case false:
		if len(tr.Nonce) != 0 {
			return fmt.Errorf("nonce must be 0 size")
		}
	}
	if tr.Time.After(time.Now()) {
		return errors.New("timestamp is in future")
	}
	if len(tr.Message) > MessageMax {
		return fmt.Errorf("message length must be under %d bytes", MessageMax)
	}
	if len(tr.Inputs) > ArrayMax {
		return errors.New("length of inputs is over 255")
	}
	for n, i := range tr.Inputs {
		if len(i.PreviousTX) != 32 {
			return fmt.Errorf("previous tx hash at %d must be 32 bytes", n)
		}
		for j := 0; j < n; j++ {
			if tr.Inputs[j].Index == i.Index && bytes.Equal(tr.Inputs[j].PreviousTX, i.PreviousTX) {
				return fmt.Errorf("input %d has a same previous and index at input %d", n, j)
			}
		}
	}
	if len(tr.MultiSigIns) > ArrayMax {
		return errors.New("length of MultiSigIns is over 255")
	}
	for n, i := range tr.MultiSigIns {
		if len(i.PreviousTX) != 32 {
			return fmt.Errorf("previous tx hash at %d must be 32 bytes", n)
		}
		for j := 0; j < n; j++ {
			if tr.MultiSigIns[j].Index == i.Index && bytes.Equal(tr.MultiSigIns[j].PreviousTX, i.PreviousTX) {
				return fmt.Errorf("input %d has a same previous and index at input %d", n, j)
			}
		}
	}
	if len(tr.Outputs) > ArrayMax {
		return errors.New("length of Outputs is over 255")
	}
	for n, o := range tr.Outputs {
		if !(typ == TxRewardFee &&
			n == len(tr.Outputs)-1) && len(o.Address) != 32 {
			return fmt.Errorf("address in outputs %d must be 32 bytes", n)
		}
		if o.Value > aklib.ADKSupply {
			return fmt.Errorf("value in outputs %d must be under %d adk",
				n, aklib.ADKSupply)
		}
	}
	if typ == TxRewardFee &&
		(len(tr.Outputs) == 0 || tr.Outputs[len(tr.Outputs)-1].Address != nil) {
		return errors.New("last address of inputs must be nil")
	}
	if len(tr.MultiSigOuts) > ArrayMax {
		return errors.New("length of MultiSigOuts is over 255")
	}
	for n, o := range tr.MultiSigOuts {
		for i, a := range o.Addresses {
			if len(a) != 32 {
				return fmt.Errorf("address size at multisig %d must be 32 bytes", n)
			}
			if len(o.Addresses) > ArrayMax {
				return errors.New("length of MultiSigOut Addresses is over 255")
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
	if len(tr.Previous) == 0 {
		return fmt.Errorf("number of previous tx must be over 0")
	}
	if len(tr.Previous) > ArrayMax {
		return errors.New("length of Previous is over 255")
	}
	for n, i := range tr.Previous {
		if len(i) != 32 {
			return fmt.Errorf("tx hash size at previous tx %d must be 32 bytes", n)
		}
		for j := n + 1; j < len(tr.Previous); j++ {
			if bytes.Equal(i, tr.Previous[j]) {
				return fmt.Errorf("previous tx %d is same as %d", n, j)
			}
		}
	}
	if tr.Easiness > cfg.Easiness {
		return fmt.Errorf("Easiness must be %d", cfg.Easiness)
	}
	if !tr.LockTime.IsZero() && tr.LockTime.After(time.Now()) {
		return errors.New("this tx is not unlocked yet")
	}
	if tr.HashType != 0 &&
		(tr.HashType&HashTypeExcludeOutputs == 0 && tr.HashType&HashTypeExcludeTicketOut == 0) {
		return fmt.Errorf("invalid hashtype %d", tr.HashType)
	}
	n := int(tr.HashType) & 0xf
	switch typ {
	case TxRewardFee:
		if tr.HashType&0xfff0 != HashTypeExcludeOutputs {
			return errors.New("hashtype of reward with fee must be 0x1X")
		}
		if n != 1 {
			return errors.New("hashtype must be 0x11")
		}
	case TxRewardTicket:
		if tr.HashType&0xfff0 != 0x20 {
			return errors.New("hashtype of reward with Ticket must be 0x2X")
		}
	}
	if tr.HashType&HashTypeExcludeOutputs != 0 && n > len(tr.Outputs) {
		return fmt.Errorf("number of outputs  is too large for hashtype %d", n)
	}
	if tr.TicketInput != nil && len(tr.TicketInput) != 32 {
		return errors.New("ticket intput must  be 32 bytes")
	}
	if tr.TicketOutput != nil && len(tr.TicketOutput) != 32 {
		return errors.New("ticket output must  be 32 bytes")
	}
	switch typ {
	case TxRewardTicket:
		if tr.TicketOutput != nil {
			return errors.New("ticket outtput must not be filled for RewardTicket")
		}
		if tr.TicketInput == nil {
			return errors.New("ticket intput must  be filled for RewardTicket")
		}
	case TxRewardFee:
		if tr.TicketInput != nil || tr.TicketOutput != nil {
			return errors.New("cannot use ticket")
		}
	case TxNormal:
		if tr.TicketInput != nil && tr.TicketOutput == nil {
			return errors.New("ticket_output is nil but ticket_input is not nil")
		}
		if tr.TicketInput == nil && tr.TicketOutput != nil {
			//Issuing a ticket
			if len(tr.Inputs) > 0 || len(tr.MultiSigIns) > 0 || len(tr.Outputs) > 0 || len(tr.MultiSigOuts) > 0 || !tr.LockTime.IsZero() ||
				tr.HashType != 0 || len(tr.Signatures) != 0 {
				return errors.New("tx content for ticket must be empty")
			}
			if tr.Easiness > cfg.TicketEasiness {
				return errors.New("PoW doesn't meet ticket difficulty")
			}
		}
	}

	if len(tr.Scripts) > 0 {
		return errors.New("cannot use scriptsd")
	}
	if len(tr.Reserved) > 0 {
		return errors.New("cannot use reserved field")
	}

	dat, err := tr.bytesForSign()
	if err != nil {
		return err
	}
	for n, sig := range tr.Signatures {
		if !address.Verify(sig, dat) {
			return fmt.Errorf("failed to verify a signature at %d", n)
		}
		for nn := n + 1; nn < len(tr.Signatures); nn++ {
			if bytes.Equal(sig.PublicKey, tr.Signatures[nn].PublicKey) {
				return fmt.Errorf("there are same publik keys in signature at %d and %d", n, nn)
			}
		}
	}
	if powed && !isValidHash(tr.Hash(), tr.Easiness) {
		return errors.New("tx does not match easiness")
	}
	return nil
}

//isValidHash reteurns true if  hash bytes h meets difficulty.
func isValidHash(h []byte, e uint32) bool {
	ea := binary.LittleEndian.Uint32(h[len(h)-4:])
	return ea <= e
}

//Hash reteurns hash of tx.
func (tr *Transaction) Hash() Hash {
	hh := sha256.Sum256(arypack.Marshal(tr))
	return hh[:]
}

//Hash reteurns hash of signature.
func (sig *Signatures) Hash() Hash {
	hh := sha256.Sum256(arypack.Marshal(sig))
	return hh[:]
}

//Hash reteurns hash of body.
func (body *Body) Hash() Hash {
	hh := sha256.Sum256(arypack.Marshal(body))
	return hh[:]
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

//CheckAll checks the tx, including other txs refered by the tx..
//Genesis block must be saved in the store
func (tr *Transaction) CheckAll(getTX GetTXFunc, cfg *aklib.Config, typ Type) error {
	if err := tr.Check(cfg, typ); err != nil {
		return err
	}
	for _, i := range tr.Previous {
		if _, err := getTX(i); err != nil {
			return err
		}
	}
	var totalout uint64
	for _, o := range tr.Outputs {
		totalout += o.Value
	}
	for _, o := range tr.MultiSigOuts {
		totalout += o.Value
	}
	var totalin uint64
	adrs := make([]*addresses, 0, len(tr.Inputs)+1)
	for _, sig := range tr.Signatures {
		hadr := sha256.Sum256(sig.PublicKey)
		adrs = append(adrs, &addresses{
			adr: hadr[:],
		})
	}
	for n, inp := range tr.Inputs {
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
	for n, inp := range tr.MultiSigIns {
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
	if len(tr.TicketInput) > 0 {
		inTX, err := getTX(tr.TicketInput)
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

//bytesForSign returns a hash slice for  signinig
func (tr *Transaction) bytesForSign() ([]byte, error) {
	tx2 := tr.Clone()
	tx2.Gnonce = 0
	tx2.Nonce = nil
	if tr.HashType&HashTypeExcludeTicketOut != 0 {
		tx2.TicketOutput = nil
	}
	if tr.HashType&0xf0 == HashTypeExcludeOutputs {
		exclude := int(tr.HashType & 0x0f)
		if len(tx2.Outputs) < exclude {
			return nil, errors.New("output length is less than one specified by hash_type")
		}
		tx2.Outputs = tx2.Outputs[:len(tr.Outputs)-exclude]
	}
	return arypack.Marshal(tx2.Body), nil
}

func (tr *Transaction) hashForPoW() []byte {
	nonce := tr.Nonce
	tr.Nonce = nil
	btx := arypack.Marshal(tr)
	h := sha256.Sum256(btx)
	tr.Nonce = nonce
	return h[:]
}

//Clone clones tx.
func (tr *Transaction) Clone() *Transaction {
	var tx2 Transaction
	if err := arypack.Unmarshal(arypack.Marshal(tr), &tx2); err != nil {
		panic(err)
	}
	return &tx2
}
