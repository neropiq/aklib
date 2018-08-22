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
	"errors"
)

//InOutHashType is a type in InOutHash struct.
type InOutHashType byte

//Types for in/out txs.
const (
	TypeIn InOutHashType = iota
	TypeMulin
	TypeTicketin
	TypeOut
	TypeMulout
	TypeTicketout
)

func (t InOutHashType) String() string {
	switch t {
	case TypeIn:
		return "input"
	case TypeMulin:
		return "multisig_input"
	case TypeTicketin:
		return "ticket_input"
	case TypeOut:
		return "output"
	case TypeMulout:
		return "multisig_output"
	case TypeTicketout:
		return "ticket_output"
	default:
		return ""
	}
}

//InoutHash represents in/out tx hashess.
type InoutHash struct {
	Hash  Hash          `json:"hash"`
	Type  InOutHashType `json:"type"`
	Index byte          `json:"index"`
}

//NewInoutHash returns a InoutHash object from serialized inout data.
func NewInoutHash(dat []byte) (*InoutHash, error) {
	if len(dat) != 34 {
		return nil, errors.New("invalid dat length")
	}
	ih := &InoutHash{
		Hash:  make(Hash, 32),
		Type:  InOutHashType(dat[32]),
		Index: dat[33],
	}
	copy(ih.Hash, dat[:32])
	return ih, nil
}

//Bytes returns a serialized inout data as a slice.
func (ih *InoutHash) Bytes() []byte {
	ary := Inout2key(ih.Hash, ih.Type, ih.Index)
	return ary[:]
}

//Serialize returns serialized a 34 bytes array.
func (ih *InoutHash) Serialize() [34]byte {
	return Inout2keyArray(ih.Hash, ih.Type, ih.Index)
}

//Inout2keyArray returns a serialized inout data as an array.
func Inout2keyArray(h Hash, typ InOutHashType, no byte) [34]byte {
	var r [34]byte
	copy(r[:], h)
	r[32] = byte(typ)
	r[33] = no
	return r
}

//Inout2key returns a slice of serialized inout data.
func Inout2key(h Hash, typ InOutHashType, no byte) []byte {
	r := Inout2keyArray(h, typ, no)
	return r[:]
}

//InputHashes returns inputs and outputs from tr.
func InputHashes(tr *Body) []*InoutHash {
	prevs := make([]*InoutHash, 0, 1+
		len(tr.Inputs)+len(tr.MultiSigIns))
	if tr.TicketInput != nil {
		prevs = append(prevs, &InoutHash{
			Type: TypeTicketin,
			Hash: tr.TicketInput,
		})
	}
	for _, prev := range tr.Inputs {
		prevs = append(prevs, &InoutHash{
			Type: TypeIn,
			Hash: prev.PreviousTX,
		})
	}
	for _, prev := range tr.MultiSigIns {
		prevs = append(prevs, &InoutHash{
			Type: TypeMulin,
			Hash: prev.PreviousTX,
		})
	}
	return prevs
}

//HashWithType is hash with tx type.
type HashWithType struct {
	Hash Hash
	Type Type
}
