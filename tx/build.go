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
	"fmt"
	"log"
	"sort"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
)

//RawOutput is an output in a tx.
type RawOutput struct {
	Address string
	Value   uint64
}

//UTXO  is an candidate of inputs in a tx.
type UTXO struct {
	Address AddressIF
	*InoutHash
	Value uint64
}

//AddressIF is an interface for sign.
type AddressIF interface {
	Sign(*Transaction) error
	String() string
}

//Wallet is a wallet interface for getting UTXOs and a new address.
type Wallet interface {
	GetUTXO(uint64) ([]*UTXO, error)
	NewChangeAddress() (*address.Address, error)
	GetLeaves() ([]Hash, error)
}

//Wallet2 is a wallet interface for getting a ticket out tx.
type Wallet2 interface {
	Wallet
	GetTicketout() (Hash, *address.Address, error)
}

//Build builds a tx for sending coins.
func Build(conf *aklib.Config, w Wallet, tag []byte, outputs []*RawOutput,
	beforeSignFunc func(*Transaction) error) (*Transaction, error) {
	ls, err := w.GetLeaves()
	if err != nil {
		return nil, err
	}
	tr := New(conf, ls...)
	tr.Message = tag
	var outtotal uint64
	for _, o := range outputs {
		outtotal += o.Value
	}
	utxos, err := w.GetUTXO(outtotal)
	if err != nil {
		return nil, err
	}
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value < utxos[j].Value
	})
	i := sort.Search(len(utxos), func(i int) bool {
		return utxos[i].Value >= outtotal
	})
	change := int64(outtotal)
	var adrs []AddressIF
	if i == len(utxos) {
		i--
	}
	for ; i >= 0 && change > 0; i-- {
		log.Println(utxos[i])
		tr.AddInput(utxos[i].Hash, utxos[i].Index)
		adrs = append(adrs, utxos[i].Address)
		change -= int64(utxos[i].Value)
	}
	if change > 0 {
		return nil, fmt.Errorf("insufficient balance %v", change)
	}
	if change != 0 {
		adr, err := w.NewChangeAddress()
		if err != nil {
			return nil, err
		}
		if err := tr.AddOutput(conf, adr.Address58(conf), uint64(-change)); err != nil {
			return nil, err
		}
	}
	//don't move it. ticket fee must be at the last of outputs.
	for _, o := range outputs {
		if err2 := tr.AddOutput(conf, o.Address, o.Value); err2 != nil {
			return nil, err2
		}
	}
	if beforeSignFunc != nil {
		if err := beforeSignFunc(tr); err != nil {
			return nil, err
		}
	}
	for _, a := range adrs {
		if err := a.Sign(tr); err != nil {
			return nil, err
		}
	}

	return tr, nil
}

//BuildParam is a param for building a tx..
type BuildParam struct {
	Comment string
	Dest    []*RawOutput
	PoWType Type
	Fee     uint64
}

//Build2 builds a tx for sending coins with fee or ticket..
func Build2(conf *aklib.Config, w Wallet2, p *BuildParam) (*Transaction, error) {
	var err error
	if p.PoWType == TypeRewardFee {
		if p.Fee == 0 {
			return nil, errors.New("fee is zero")
		}
		p.Dest = append(p.Dest, &RawOutput{
			Address: "",
			Value:   p.Fee,
		})
	}
	var ticketadr *address.Address
	f := func(tr *Transaction) error {
		switch p.PoWType {
		case TypeRewardFee:
			tr.Body.HashType = HashTypeExcludeOutputs | 0x1
		case TypeRewardTicket:
			tr.Body.HashType = HashTypeExcludeTicketOut
			var h Hash
			h, ticketadr, err = w.GetTicketout()
			if err != nil {
				return err
			}
			tr.TicketInput = h
		}
		return nil
	}
	tr, err := Build(conf, w, []byte(p.Comment), p.Dest, f)
	if err != nil {
		return nil, err
	}
	if p.PoWType == TypeRewardTicket {
		if err := tr.Sign(ticketadr); err != nil {
			return nil, err
		}
	}
	return tr, nil
}
