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

package rpc

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/AidosKuneen/aklib/tx"
	"github.com/AidosKuneen/consensus"
)

//Request is for parsing request from client.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

//Response is for respoding to clinete in jsonrpc.
type Response struct {
	Result interface{} `json:"result"`
	Error  *Err        `json:"error"`
	ID     interface{} `json:"id"`
}

//Err represents error struct for response.
type Err struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

//Bans is a struct for listbanned RPC.
type Bans struct {
	Address string `json:"address"`
	Created int64  `json:"ban_created"`
	Until   int64  `json:"banned_until"`
	Reason  string `json:"ban_reason"`
}

//NodeInfo is a struct for getnodeinfo RPC.
type NodeInfo struct {
	Version         string  `json:"version"`
	ProtocolVersion int     `json:"protocolversion"`
	WalletVersion   int     `json:"walletversion"`
	Balance         *uint64 `json:"balance,omitempty"`
	Connections     int     `json:"connections"`
	Proxy           string  `json:"proxy"`
	Testnet         byte    `json:"testnet"`
	KeyPoolSize     int     `json:"keypoolsize"`
	Leaves          int     `json:"leaves"`
	Time            int64   `json:"time"`
	TxNo            uint64  `json:"txno"`
	LatestLedger    string  `json:"latest_ledger"`
}

//InoutHash is a struct for getlasthistory RPC.
type InoutHash struct {
	Hash  string           `json:"hash"`
	Type  tx.InOutHashType `json:"type"`
	Index byte             `json:"index"`
}

//TxStatus represents a confirmation status of a tx.
type TxStatus struct {
	Hash        string `json:"hash"`
	Exists      bool   `json:"exists"`
	IsRejected  bool   `json:"is_rejected"`
	IsConfirmed bool   `json:"is_confirmed"`
	LedgerID    string `json:"ledger_id"`
}

//IsAccepted returns true if the tx is confirmed and accepted.
func (s *TxStatus) IsAccepted() bool {
	return s.Exists && !s.IsRejected && s.IsConfirmed
}

//Ledger represents a ledger info.
type Ledger struct {
	ID                  string        `json:"id"`
	ParentID            string        `json:"parent_id"`
	Seq                 consensus.Seq `json:"sequence_no"`
	Txs                 string        `json:"transaction_id"`
	CloseTimeResolution time.Duration `json:"closetime_resolution"`
	CloseTime           time.Time     `json:"closetime"`
	ParentCloseTime     time.Time     `json:"parent_closetime"`
	CloseTimeAgree      bool          `json:"closetime_agree"`
}

//NewLedger converts the leger.
func NewLedger(tr *consensus.Ledger) *Ledger {
	id := tr.ID()
	var h consensus.TxID
	for h = range tr.Txs {
	}
	l := &Ledger{
		ID:                  hex.EncodeToString(id[:]),
		ParentID:            hex.EncodeToString(tr.ParentID[:]),
		Seq:                 tr.Seq,
		CloseTimeResolution: tr.CloseTimeResolution,
		CloseTime:           tr.CloseTime,
		ParentCloseTime:     tr.ParentCloseTime,
		CloseTimeAgree:      tr.CloseTimeAgree,
	}
	if len(tr.Txs) > 0 {
		l.Txs = hex.EncodeToString(h[:])
	}
	return l
}

//Info is a struct for validateaddress RPC.
type Info struct {
	IsValid      bool    `json:"isvalid"`
	Address      string  `json:"address"`
	ScriptPubKey string  `json:"scriptPubkey"`
	IsMine       bool    `json:"ismine"`
	IsWatchOnly  *bool   `json:"iswatchonly,omitempty"`
	IsScript     *bool   `json:"isscript,omitempty"`
	Pubkey       *string `json:"pubkey,omitempty"`
	IsCompressed *bool   `json:"iscompressed,omitempty"`
	Account      *string `json:"account,omitempty"`
}

//Details is a struct for gettransaction RPC.
type Details struct {
	Account   string  `json:"account"`
	Address   string  `json:"address"`
	Category  string  `json:"category"`
	Amount    float64 `json:"amount"`
	Vout      int64   `json:"vout"`
	Fee       float64 `json:"fee"`
	Abandoned *bool   `json:"abandoned,omitempty"`
}

//Gettx is a struct for gettransaction RPC.
type Gettx struct {
	Amount            float64    `json:"amount"`
	Fee               float64    `json:"fee"`
	Confirmations     int        `json:"confirmations"`
	Blockhash         *string    `json:"blockhash,omitempty"`
	Blockindex        *int64     `json:"blockindex,omitempty"`
	Blocktime         *int64     `json:"blocktime,omitempty"`
	Txid              string     `json:"txid"`
	WalletConflicts   []string   `json:"walletconflicts"`
	Time              int64      `json:"time"`
	TimeReceived      int64      `json:"timereceived"`
	BIP125Replaceable string     `json:"bip125-replaceable"`
	Details           []*Details `json:"details"`
	Hex               string     `json:"hex"`
}

//Transaction is a struct for listtransactions RPC.
type Transaction struct {
	Account  *string `json:"account"`
	Address  string  `json:"address"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	// Label             string      `json:"label"`
	Vout          int64   `json:"vout"`
	Fee           float64 `json:"fee"`
	Confirmations int     `json:"confirmations"`
	Trusted       *bool   `json:"trusted,omitempty"`
	// Generated         bool        `json:"generated"`
	Blockhash       *string  `json:"blockhash,omitempty"`
	Blockindex      *int64   `json:"blockindex,omitempty"`
	Blocktime       *int64   `json:"blocktime,omitempty"`
	Txid            string   `json:"txid"`
	Walletconflicts []string `json:"walletconflicts"`
	Time            int64    `json:"time"`
	TimeReceived    int64    `json:"timereceived"`
	// Comment           string      `json:"string"`
	// To                string `json:"to"`
	// Otheraccount      string `json:"otheraccount"`
	BIP125Replaceable string `json:"bip125-replaceable"`
	Abandoned         *bool  `json:"abandoned,omitempty"`
}

//ToDetail converts dt to Details struct.
func (dt *Transaction) ToDetail() (*Details, error) {
	if dt.Account == nil {
		return nil, errors.New("Account is nil")
	}
	return &Details{
		Account:   *dt.Account,
		Address:   dt.Address,
		Category:  dt.Category,
		Amount:    dt.Amount,
		Vout:      dt.Vout,
		Abandoned: dt.Abandoned,
	}, nil
}
