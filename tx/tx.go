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
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/aklib/arypack"
)

//max length of tx and fields in a transaction.
const (
	MessageMax          = 255
	TransactionMax      = 2000000
	ArrayMax            = 255
	DefaultPreviousSize = 2
)

var (
	//TypeNormal is a type for nomal tx.
	typeNormal = []byte{0xAD, 0xBF, 0x43, 0x01}
)

//Hash is a tx hash.
type Hash []byte

//Array converts slice to array
func (h Hash) Array() [32]byte {
	var h32 [32]byte
	copy(h32[:], h)
	return h32
}

func (h Hash) String() string {
	return hex.EncodeToString(h)
}

//GetTXFunc gets tx.
type GetTXFunc func(hash []byte) (*Body, error)

//Types when hashing a tx.
const (
	//HashTypeExcludeOutputs is for excluding some outputs.
	HashTypeNormal           = 0x0
	HashTypeExcludeOutputs   = 0x10
	HashTypeExcludeTicketOut = 0x20
)

//Type is a tx type.
type Type byte

//ByteSlice is a byte slice marsharing to hex string.
type ByteSlice []byte

//type for minable tx.
const (
	TypeNormal Type = iota
	TypeRewardTicket
	TypeRewardFee
	TypeNotPoWed
)

//Input is an input in transactions.
type Input struct {
	PreviousTX Hash `json:"previous_tx"` //32 bytes
	Index      byte `json:"index"`
}

//Output is an output in transactions.
type Output struct {
	Address address.Bytes `json:"address"` //65 bytes
	Value   uint64        `json:"value"`   //8 bytes
}

//MultisigStruct is a structure of  multisig.
type MultisigStruct struct {
	M         byte            `json:"n"`         //0 means normal payment, or M out of len(Address) multisig.
	Addresses []address.Bytes `json:"addresses"` //< 65 * 32 bytes

}

//MultiSigOut is an multisig output in transactions.
type MultiSigOut struct {
	MultisigStruct
	Value uint64 `json:"value"` //8 bytes
}

//AddressByte returns a multisig address in binary form.
func (mout *MultisigStruct) AddressByte(cfg *aklib.Config) address.Bytes {
	adr := make([]address.Bytes, len(mout.Addresses))
	for i, a := range mout.Addresses {
		adr[i] = a
	}
	return address.MultisigAddressByte(cfg, mout.M, adr...)
}

//Address returns a multisig address.
func (mout *MultisigStruct) Address(cfg *aklib.Config) string {
	adr := make([]address.Bytes, len(mout.Addresses))
	for i, a := range mout.Addresses {
		adr[i] = a
	}
	return address.MultisigAddress(cfg, mout.M, adr...)
}

//MultiSigIn is an multisig input in transactions.
type MultiSigIn struct {
	PreviousTX Hash `json:"previous_tx"` //32 bytes
	Index      byte `json:"index"`
}

//Body is a Transactoin except signature.
type Body struct {
	Type         ByteSlice      `json:"type"`                    //4 bytes
	Nonce        []uint32       `json:"nonce"`                   //20*VarInt(<4)
	Gnonce       uint32         `json:"g_nonce"`                 //4 bytes
	Time         time.Time      `json:"time"`                    //8 bytes
	Message      ByteSlice      `json:"message,omitempty"`       //<255 bytes
	Inputs       []*Input       `json:"inputs,omitempty"`        //<33 * 32 = 1056 bytes
	MultiSigIns  []*MultiSigIn  `json:"multisig_ins,omitempty"`  //<1032 * 4 = 4128 bytes
	Outputs      []*Output      `json:"outputs,omitempty"`       //<40 * 32 = 1280 bytes
	MultiSigOuts []*MultiSigOut `json:"multisig_outs,omitempty"` //<1032 * 4 = 4128 bytes
	Parent       []Hash         `json:"parent"`                  //<8*32=256 bytes
	Easiness     uint32         `json:"easiness"`                //1 byte //not used for now
	LockTime     time.Time      `json:"lock_time"`               // 4 bytes
	HashType     uint16         `json:"hash_type"`
	TicketInput  Hash           `json:"ticket_input,omitempty"`
	TicketOutput address.Bytes  `json:"ticket_output,omitempty"`
	Scripts      [][]byte       `json:"-"` //not used
	Reserved     []byte         `json:"-"` //not used
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
			Time:     time.Now().Truncate(time.Second),
			Easiness: s.Easiness,
			Parent:   previous,
		},
	}
}

//IssueTicket make and does PoW for a  transaction for issuing tx.
func IssueTicket(s *aklib.Config, ticketOut []byte, previous ...Hash) (*Transaction, error) {
	tr := &Transaction{
		Body: &Body{
			Type:         typeNormal,
			Time:         time.Now().Truncate(time.Second),
			Easiness:     s.TicketEasiness,
			TicketOutput: ticketOut,
			Parent:       previous,
		},
	}
	return tr, tr.PoW()
}

//NewMinableFee returns a minable transaction by fee..
func NewMinableFee(s *aklib.Config, previous ...Hash) *Transaction {
	return &Transaction{
		Body: &Body{
			Type:     typeNormal,
			Time:     time.Now().Truncate(time.Second),
			Easiness: s.Easiness,
			HashType: HashTypeExcludeOutputs | 0x1,
			Parent:   previous,
		},
	}
}

//NewMinableTicket returns a minable transaction by ticket..
func NewMinableTicket(s *aklib.Config, ticketIn Hash, previous ...Hash) *Transaction {
	return &Transaction{
		Body: &Body{
			Type:        typeNormal,
			Time:        time.Now().Truncate(time.Second),
			Easiness:    s.Easiness,
			HashType:    HashTypeExcludeTicketOut,
			Parent:      previous,
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
func (body *Body) AddOutput(cfg *aklib.Config, adr string, v uint64) error {
	var pub []byte
	var err error
	if adr != "" {
		pub, _, err = address.ParseAddress58(cfg, adr)
		if err != nil {
			return err
		}
	}
	body.Outputs = append(body.Outputs, &Output{
		Address: pub,
		Value:   v,
	})
	return nil
}

//AddMultisigIn add a multisig input into tx.
func (body *Body) AddMultisigIn(h Hash, idx byte) {
	body.MultiSigIns = append(body.MultiSigIns, &MultiSigIn{
		PreviousTX: h,
		Index:      idx,
	})
}

//AddMultisigOut add a mulsig output into tx.
func (body *Body) AddMultisigOut(cfg *aklib.Config, m byte, v uint64, adrs ...string) error {
	if len(adrs) < int(m) {
		return errors.New("length of adrs is less than n")
	}
	as := make([]address.Bytes, len(adrs))
	for i, adr := range adrs {
		pub, _, err := address.ParseAddress58(cfg, adr)
		if err != nil {
			return err
		}
		if !checkAdrsPrefix(cfg, pub) {
			return errors.New("invalid address for this network")
		}
		as[i] = pub
	}
	body.MultiSigOuts = append(body.MultiSigOuts, &MultiSigOut{
		MultisigStruct: MultisigStruct{
			M:         m,
			Addresses: as,
		},
		Value: v,
	})
	return nil
}

//Sign sings the tx.
func (tr *Transaction) Sign(a *address.Address) error {
	dat, err := tr.bytesForSign()
	if err != nil {
		return err
	}
	sig, err := a.Sign(dat)
	if err != nil {
		return err
	}
	tr.Signatures = append(tr.Signatures, sig)
	return nil
}

//Signature returns  singture of the tx.
func (tr *Transaction) Signature(a *address.Address) (*address.Signature, error) {
	dat, err := tr.bytesForSign()
	if err != nil {
		return nil, err
	}
	return a.Sign(dat)
}

//AddSig adds a signature  to tx.
func (tr *Transaction) AddSig(sig *address.Signature) {
	tr.Signatures = append(tr.Signatures, sig)
}

//Size returns tx size.
func (tr *Transaction) Size() int {
	return len(arypack.Marshal(tr))
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

//UnmarshalJSON sets *bs to a copy of data.
func (bs *ByteSlice) UnmarshalJSON(b []byte) error {
	h := ""
	if err := json.Unmarshal(b, &h); err != nil {
		return err
	}
	var err error
	*bs, err = hex.DecodeString(h)
	return err
}

//MarshalJSON returns m as the JSON encoding of m.
func (bs *ByteSlice) MarshalJSON() ([]byte, error) {
	h := hex.EncodeToString(*bs)
	return json.Marshal(&h)
}

//UnmarshalJSON sets *bs to a copy of data.
func (h *Hash) UnmarshalJSON(b []byte) error {
	he := ""
	if err := json.Unmarshal(b, &he); err != nil {
		return err
	}
	var err error
	*h, err = hex.DecodeString(he)
	return err
}

//MarshalJSON returns m as the JSON encoding of m.
func (h *Hash) MarshalJSON() ([]byte, error) {
	he := hex.EncodeToString(*h)
	return json.Marshal(&he)
}
