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

package rpcclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/AidosKuneen/aklib/arypack"
	"github.com/AidosKuneen/aklib/tx"
	"github.com/AidosKuneen/aknode/msg"
	"github.com/AidosKuneen/aknode/rpc"
)

func (client *RPC) request(method string, params interface{}, out interface{}) error {
	req := &rpc.Request{
		JSONRPC: "1.0",
		ID:      "aklib",
		Method:  method,
	}
	var err error
	req.Params, err = json.Marshal(params)
	if err != nil {
		return err
	}
	return client.do(req, out)
}

//RPC is for calling RPCs.
type RPC struct {
	client   *http.Client
	Endpoint string
	User     string
	Password string
}

//New takes an (optional) endpoint and optional http.Client and returns
// an RPC struct.
func New(endpoint, user, password string, c *http.Client) *RPC {
	if c == nil {
		c = http.DefaultClient
	}
	return &RPC{
		client:   c,
		Endpoint: endpoint,
		User:     user,
		Password: password,
	}
}

func (client *RPC) do(cmd interface{}, out interface{}) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	rd := bytes.NewReader(b)
	req, err := http.NewRequest("POST", client.Endpoint, rd)
	if err != nil {
		return err
	}
	if client.User != "" || client.Password != "" {
		req.SetBasicAuth(client.User, client.Password)
		req.Header.Set("Content-Type", "text/plain")
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var r rpc.Response
	if err := json.Unmarshal(bs, &r); err != nil {
		log.Println(string(bs))
		return err
	}
	if r.Error != nil {
		return errors.New(r.Error.Message)
	}
	return json.Unmarshal(bs, out)
}

//SetTxFee sends a settxfee RPC.
func (client *RPC) SetTxFee(amount float64) (bool, error) {
	var out struct {
		Result bool `json:"result"`
	}
	err := client.request("settxfee", []float64{amount}, &out)
	return out.Result, err
}

//GetNewAddress sends a getnewaddress RPC.
func (client *RPC) GetNewAddress(account string) (string, error) {
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("getnewaddress", []string{account}, &out)
	return out.Result, err
}

//WalletPassphrase sends a walletpassphrase RPC.
func (client *RPC) WalletPassphrase(pp string, sec int) error {
	ary := make([]interface{}, 0, 2)
	ary = append(ary, pp)
	ary = append(ary, sec)
	var out struct {
		Result interface{} `json:"result"`
	}
	return client.request("walletpassphrase", ary, &out)
}

//SendMany sends a sendmany RPC.
func (client *RPC) SendMany(account string, amounts map[string]float64) (string, error) {
	ary := make([]interface{}, 0, 2)
	ary = append(ary, account)
	ary = append(ary, amounts)
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("sendmany", ary, &out)
	return out.Result, err
}

//ListTransactions sends a listtransactions RPC.
func (client *RPC) ListTransactions(account string, count, skip int) ([]*rpc.Transaction, error) {
	ary := make([]interface{}, 0, 2)
	ary = append(ary, account)
	ary = append(ary, count)
	ary = append(ary, skip)
	var out struct {
		Result []*rpc.Transaction `json:"result"`
	}
	err := client.request("listtransactions", ary, &out)
	return out.Result, err
}

//ValidateAddress sends a validateaddress RPC.
func (client *RPC) ValidateAddress(address string) (*rpc.Info, error) {
	var out struct {
		Result *rpc.Info `json:"result"`
	}
	err := client.request("validateaddress", []string{address}, &out)
	return out.Result, err
}

//GetBalance sends a getbalance RPC.
func (client *RPC) GetBalance(account string) (float64, error) {
	var out struct {
		Result float64 `json:"result"`
	}
	err := client.request("getbalance", []string{account}, &out)
	return out.Result, err
}

//GetTransaction sends a gettransactions RPC.
func (client *RPC) GetTransaction(txid string) (*rpc.Gettx, error) {
	var out struct {
		Result *rpc.Gettx `json:"result"`
	}
	err := client.request("gettransaction", []string{txid}, &out)
	return out.Result, err
}

//ListPeer sends a listpeer RPC.
func (client *RPC) ListPeer() ([]msg.Addr, error) {
	var out struct {
		Result []msg.Addr `json:"result"`
	}
	err := client.request("listpeer", nil, &out)
	return out.Result, err
}

//DumpPrivKey sends a dumpprivkey RPC.
func (client *RPC) DumpPrivKey() (string, error) {
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("dumpprivkey", nil, &out)
	return out.Result, err
}

//ListBanned sends a listbanned RPC.
func (client *RPC) ListBanned() ([]*rpc.Bans, error) {
	var out struct {
		Result []*rpc.Bans `json:"result"`
	}
	err := client.request("listbanned", nil, &out)
	return out.Result, err
}

//Stop sends a stop RPC.
func (client *RPC) Stop() error {
	var out struct {
		Result string `json:"result"`
	}
	return client.request("stop", nil, &out)
}

//DumpWallet sends a dumpwallet RPC.
func (client *RPC) DumpWallet(fname string) error {
	var out struct {
		Result struct{} `json:"result"`
	}
	return client.request("dumpwallet", []interface{}{fname}, &out)
}

//ImportWallet sends a importwallet RPC.
func (client *RPC) ImportWallet(fname string) error {
	var out struct {
		Result struct{} `json:"result"`
	}
	return client.request("importwallet", []interface{}{fname}, &out)
}

//SendRawTX sends a sendrawtx RPC.
func (client *RPC) SendRawTX(tx *tx.Transaction, typ tx.Type) (string, error) {
	ttx := arypack.Marshal(tx)
	var out struct {
		Result string `json:"result"`
	}
	return out.Result, client.request("sendrawtx", []interface{}{ttx, typ}, &out)
}

//GetNodeinfo sends a getnodeinfo RPC.
func (client *RPC) GetNodeinfo() (*rpc.NodeInfo, error) {
	var out struct {
		Result *rpc.NodeInfo `json:"result"`
	}
	err := client.request("getnodeinfo", []interface{}{}, &out)
	return out.Result, err
}

//GetLeaves sends a getleaves RPC.
func (client *RPC) GetLeaves() ([]string, error) {
	var out struct {
		Result []string `json:"result"`
	}
	err := client.request("getleaves", []interface{}{}, &out)
	return out.Result, err
}

//GetHistory sends a getlasthistory RPC.
func (client *RPC) GetHistory(adr string, all bool) ([]*rpc.InoutHash, error) {
	var out struct {
		Result []*rpc.InoutHash `json:"result"`
	}
	err := client.request("gethistory", []interface{}{adr, all}, &out)
	return out.Result, err
}

//GetRawTx sends a getrawtx RPC.
func (client *RPC) GetRawTx(txid string) (*tx.Transaction, error) {
	var out struct {
		Result []byte `json:"result"`
	}
	err := client.request("getrawtx", []interface{}{txid}, &out)
	if err != nil {
		return nil, err
	}
	var tr tx.Transaction
	err = arypack.Unmarshal(out.Result, &tr)
	return &tr, err
}

func (client *RPC) getMinableTx(arg interface{}) (*tx.Transaction, error) {
	var out struct {
		Result []byte `json:"result"`
	}
	err := client.request("getminabletx", []interface{}{arg}, &out)
	if err != nil {
		return nil, err
	}
	var tr tx.Transaction
	err = arypack.Unmarshal(out.Result, &tr)
	return &tr, err
}

//GetMinableTicketTx sends a getminabletx RPC for tiecket reward.
func (client *RPC) GetMinableTicketTx() (*tx.Transaction, error) {
	return client.getMinableTx("ticket")
}

//GetMinableFeeTx sends a getminabletx RPC for fee reward.
func (client *RPC) GetMinableFeeTx(fee float64) (*tx.Transaction, error) {
	return client.getMinableTx(fee)
}

//GetTxsStatus sends a gettxsstatus RPC.
func (client *RPC) GetTxsStatus(txid ...string) ([]int, error) {
	var out struct {
		Result []int `json:"result"`
	}
	arg := make([]interface{}, len(txid))
	for i := range txid {
		arg[i] = txid[i]
	}
	err := client.request("gettxsstatus", arg, &out)
	if err != nil {
		return nil, err
	}
	return out.Result, err
}

//WalletLock sends a walletlock RPC.
func (client *RPC) WalletLock() error {
	var out struct {
		Result struct{} `json:"result"`
	}
	return client.request("walletlock", []interface{}{}, &out)
}

//SendFrom sends a sendfrom RPC.
func (client *RPC) SendFrom(account, address string, amount float64) (string, error) {
	ary := make([]interface{}, 0, 3)
	ary = append(ary, account)
	ary = append(ary, address)
	ary = append(ary, amount)
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("sendfrom", ary, &out)
	return out.Result, err
}

//SendToAddress sends a sendtoaddress RPC.
func (client *RPC) SendToAddress(address string, amount float64) (string, error) {
	ary := make([]interface{}, 0, 2)
	ary = append(ary, address)
	ary = append(ary, amount)
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("sendtoaddress", ary, &out)
	return out.Result, err
}

//AddressGroup is a result of ListAddressGroupings.
type AddressGroup struct {
	Address string
	Amount  float64
	Account *string
}

//ListAddressGroupings sends a listaddressgroupings RPC.
func (client *RPC) ListAddressGroupings() ([]*AddressGroup, error) {
	var out struct {
		Result [][][]interface{} `json:"result"`
	}
	err := client.request("listaddressgroupings", nil, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Result) != 1 {
		return nil, errors.New("invalid length of result")
	}
	r := make([]*AddressGroup, len(out.Result[0]))
	for i, o := range out.Result[0] {
		if len(o) > 3 || len(o) < 2 {
			return nil, errors.New("invalid length of result")
		}
		a := &AddressGroup{}
		var ok bool
		a.Address, ok = o[0].(string)
		if !ok {
			return nil, errors.New("#1 of result must be string")
		}
		a.Amount, ok = o[1].(float64)
		if !ok {
			return nil, errors.New("#2 of result must be float64")
		}
		if len(o) == 3 {
			acc, ok := o[2].(string)
			if !ok {
				return nil, errors.New("#3 of result must be string")
			}
			a.Account = &acc
		}
		r[i] = a
	}
	return r, nil
}

//ListAccounts sends a listaccounts RPC.
func (client *RPC) ListAccounts() (map[string]float64, error) {
	var out struct {
		Result map[string]float64 `json:"result"`
	}
	err := client.request("listaccounts", nil, &out)
	return out.Result, err
}

//GetAccount sends a getaccount RPC.
func (client *RPC) GetAccount(account string) (string, error) {
	var out struct {
		Result string `json:"result"`
	}
	err := client.request("getaccount", []interface{}{account}, &out)
	return out.Result, err
}
