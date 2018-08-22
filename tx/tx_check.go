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
	"errors"
	"fmt"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/cuckoo"
)

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

func checkAdrsPrefix(cfg *aklib.Config, adr Address) bool {
	if len(adr) != 35 {
		return false
	}
	for _, p := range cfg.PrefixAdrs {
		if bytes.Equal(p, adr[:len(p)]) {
			return true
		}
	}
	return false
}

//Check checks the tx.
func (tr *Transaction) Check(cfg *aklib.Config, typ Type) error {
	powed := true
	if typ == TypeRewardFee || typ == TypeRewardTicket || typ == TypeNotPoWed {
		powed = false
	}
	if typ == TypeNotPoWed {
		typ = TypeNormal
	}
	if typ != TypeNormal &&
		typ != TypeRewardFee && typ != TypeRewardTicket {
		return errors.New("invalid reward type")
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
			return fmt.Errorf("nonce must be %d size, but %d", cuckoo.ProofSize, len(tr.Nonce))
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
		if !(typ == TypeRewardFee &&
			n == len(tr.Outputs)-1) && !checkAdrsPrefix(cfg, o.Address) {
			return fmt.Errorf("incorrect address bytes in outputs %d", n)
		}
		if o.Value > aklib.ADKSupply {
			return fmt.Errorf("value in outputs %d must be under %d adk",
				n, aklib.ADKSupply)
		}
	}
	if typ == TypeRewardFee &&
		(len(tr.Outputs) == 0 || tr.Outputs[len(tr.Outputs)-1].Address != nil) {
		return errors.New("last address of inputs must be nil")
	}
	if len(tr.MultiSigOuts) > ArrayMax {
		return errors.New("length of MultiSigOuts is over 255")
	}
	for n, o := range tr.MultiSigOuts {
		for i, a := range o.Addresses {
			if !checkAdrsPrefix(cfg, a) {
				return fmt.Errorf("incorrect address format in output %d", n)
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
		if o.M > byte(len(o.Addresses)) {
			return fmt.Errorf("M at multisig %d must be under number of address %d",
				n, len(o.Addresses))
		}
		if o.Value > aklib.ADKSupply {
			return fmt.Errorf("value in multisig %d must be under %d adk",
				n, aklib.ADKSupply)
		}
	}
	if len(tr.Parent) == 0 {
		return fmt.Errorf("number of previous tx must be over 0")
	}
	if len(tr.Parent) > ArrayMax {
		return errors.New("length of Previous is over 255")
	}
	for n, i := range tr.Parent {
		if len(i) != 32 {
			return fmt.Errorf("tx hash size at previous tx %d must be 32 bytes", n)
		}
		for j := n + 1; j < len(tr.Parent); j++ {
			if bytes.Equal(i, tr.Parent[j]) {
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
	case TypeRewardFee:
		if tr.HashType&0xfff0 != HashTypeExcludeOutputs {
			return errors.New("hashtype of reward with fee must be 0x1X")
		}
		if n != 1 {
			return errors.New("hashtype must be 0x11")
		}
	case TypeRewardTicket:
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
	if tr.TicketOutput != nil && !checkAdrsPrefix(cfg, tr.TicketOutput) {
		return errors.New("incorrect ticket output format")
	}
	switch typ {
	case TypeRewardTicket:
		if tr.TicketOutput != nil {
			return errors.New("ticket outtput must not be filled for RewardTicket")
		}
		if tr.TicketInput == nil {
			return errors.New("ticket intput must  be filled for RewardTicket")
		}
	case TypeRewardFee:
		if tr.TicketInput != nil || tr.TicketOutput != nil {
			return errors.New("cannot use ticket")
		}
	case TypeNormal:
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

func (tr *Transaction) total(getTX GetTXFunc, cfg *aklib.Config) (uint64, uint64, error) {
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
		a, err := sig.Address(cfg)
		if err != nil {
			return 0, 0, err
		}
		adrs = append(adrs, &addresses{
			adr: a,
		})
	}
	for n, inp := range tr.Inputs {
		inTX, err := getTX(inp.PreviousTX)
		if err != nil {
			return 0, 0, err
		}
		if len(inTX.Outputs) <= int(inp.Index) {
			return 0, 0, fmt.Errorf("invalid input index, should be under  %d", len(inTX.Outputs))
		}
		totalin += inTX.Outputs[inp.Index].Value
		inTXAdr := inTX.Outputs[inp.Index].Address
		if !hasAddress(adrs, inTXAdr) {
			return 0, 0, fmt.Errorf("no signature for input %d", n)
		}
	}
	for n, inp := range tr.MultiSigIns {
		inTX, err := getTX(inp.PreviousTX)
		if err != nil {
			return 0, 0, err
		}
		if len(inTX.MultiSigOuts) <= int(inp.Index) {
			return 0, 0, fmt.Errorf("invalid multisig index, should be under  %d", len(inTX.MultiSigOuts))
		}
		mul := inTX.MultiSigOuts[inp.Index]
		totalin += mul.Value
		exist := 0
		for _, adr := range mul.Addresses {
			if hasAddress(adrs, adr) {
				exist++
			}
		}
		if exist != int(mul.M) {
			return 0, 0, fmt.Errorf("invalid number of valid signatures %d in multisig %d, should be %d", exist, n, mul.M)
		}
	}
	if len(tr.TicketInput) > 0 {
		inTX, err := getTX(tr.TicketInput)
		if err != nil {
			return 0, 0, err
		}
		if !hasAddress(adrs, inTX.TicketOutput) {
			return 0, 0, errors.New("cannot verify the ticket input")
		}
	}
	if hasUunused(adrs) {
		return 0, 0, errors.New("there are(is) unsed signature")
	}
	return totalin, totalout, nil
}

//CheckAll checks the tx, including other txs refered by the tx..
//Genesis block must be saved in the store
func (tr *Transaction) CheckAll(getTX GetTXFunc, cfg *aklib.Config, typ Type) error {
	if err := tr.Check(cfg, typ); err != nil {
		return err
	}
	for _, i := range tr.Parent {
		if _, err := getTX(i); err != nil {
			return err
		}
	}
	tin, tout, err := tr.total(getTX, cfg)
	if err != nil {
		return err
	}
	if tin != tout {
		return fmt.Errorf("total input ADK %v does not equal to one of output %v",
			tin, tout)
	}
	return nil
}
