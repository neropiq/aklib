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
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/aklib/arypack"
)

//Vbody is a Validator Transaction except signature.
//Message is used for broadcasting validator PKs and updating the PK of a validator
type Vbody struct {
	Type       ByteSlice `json:"type"`              //4 bytes
	Time       time.Time `json:"time"`              //8 bytes
	Message    ByteSlice `json:"message,omitempty"` //<255 bytes
	Parent     [2]Hash    `json:"parent"`            //<8*32=256 bytes
	Validators Hash      `json:"validators"`        //32 byte
	Vote       []byte    `json:"-"`                 //not used
}

//Vtransaction is a validator transaction in Aidos Kuneen.
type Vtransaction struct {
	*Vbody     `json:"vbody"`
	Signatures `json:"signatures"`
}

//NewV returns a validator transaction object.
func NewV(s *aklib.Config, validators Hash, previous [2]Hash) *Vtransaction {
	return &Vtransaction{
		Vbody: &Vbody{
			Type:       typeValidator,
			Time:       time.Now().Truncate(time.Second),
			Parent:     previous,
			Validators: validators,
		},
	}
}

//Sign sings the tx.
func (tr *Vtransaction) Sign(a *address.Address) error {
	dat, err := tr.bytesForSign()
	if err != nil {
		return err
	}
	tr.Signatures = append(tr.Signatures, a.Sign(dat))
	return nil
}

//Signature returns  singture of the tx.
func (tr *Vtransaction) Signature(a *address.Address) (*address.Signature, error) {
	dat, err := tr.bytesForSign()
	if err != nil {
		return nil, err
	}
	return a.Sign(dat), nil
}

//AddSig adds a signature  to tx.
func (tr *Vtransaction) AddSig(sig *address.Signature) {
	tr.Signatures = append(tr.Signatures, sig)
}

//Size returns tx size.
func (tr *Vtransaction) Size() int {
	return len(arypack.Marshal(tr))
}

//Hash returns hash of tx.
func (tr *Vtransaction) Hash() Hash {
	hh := sha256.Sum256(arypack.Marshal(tr))
	return hh[:]
}

//bytesForSign returns a hash slice for  signinig
func (tr *Vtransaction) bytesForSign() ([]byte, error) {
	tx2 := tr.Clone()
	return arypack.Marshal(tx2.Vbody), nil
}

//Clone clones tx.
func (tr *Vtransaction) Clone() *Vtransaction {
	var tx2 Vtransaction
	if err := arypack.Unmarshal(arypack.Marshal(tr), &tx2); err != nil {
		panic(err)
	}
	return &tx2
}