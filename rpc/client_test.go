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

package rpc_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AidosKuneen/aklib"
	"github.com/AidosKuneen/aklib/address"
	"github.com/AidosKuneen/aklib/arypack"
	"github.com/AidosKuneen/aklib/db"
	rpcc "github.com/AidosKuneen/aklib/rpc"
	"github.com/AidosKuneen/aklib/tx"
	"github.com/AidosKuneen/aknode/akconsensus"
	"github.com/AidosKuneen/aknode/imesh"
	"github.com/AidosKuneen/aknode/imesh/leaves"
	"github.com/AidosKuneen/aknode/node"
	"github.com/AidosKuneen/aknode/rpc"
	"github.com/AidosKuneen/aknode/setting"
	"github.com/AidosKuneen/consensus"
	"github.com/dgraph-io/badger"
)

var s setting.Setting
var a, b *address.Address
var genesis tx.Hash
var l, l2 net.Listener
var tdir string
var walletpwd = "hoe"
var cnotify chan []tx.Hash

func setup(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	var err error
	tdir, err = ioutil.TempDir("", "gotest")
	if err != nil {
		t.Fatal(err)
	}
	tdir = filepath.Join(tdir, "tmp.dat")
	var err2 error
	if err = os.RemoveAll("./test_db"); err != nil {
		log.Println(err)
	}
	s.DB, err2 = db.Open("./test_db")
	if err2 != nil {
		panic(err2)
	}

	s.Config = aklib.DebugConfig
	s.MaxConnections = 1
	s.Bind = "127.0.0.1"
	s.Port = 9624
	s.MyHostPort = ":9624"
	seed := address.GenerateSeed32()
	a, err2 = address.New(s.Config, seed)
	if err2 != nil {
		t.Error(err2)
	}
	s.Config.Genesis = map[string]uint64{
		a.Address58(s.Config): aklib.ADKSupply,
	}
	seed = address.GenerateSeed32()
	b, err2 = address.New(s.Config, seed)
	if err2 != nil {
		t.Error(err2)
	}
	s.RPCUser = "user"
	s.RPCPassword = "user"
	s.RPCTxTag = "test"
	s.RPCBind = "127.0.0.1"
	s.RPCPort = s.Config.DefaultRPCPort
	s.WalletNotify = "echo %s"
	leaves.Init(&s)
	if err = imesh.Init(&s); err != nil {
		t.Error(err)
	}
	gs := leaves.Get(1)
	if len(gs) != 1 {
		t.Error("invalid genesis")
	}
	genesis = gs[0]

	l, err = node.Start(&s, true)
	if err != nil {
		t.Error(err)
	}
	if err = rpc.Init(&s); err != nil {
		t.Error(err)
	}
	if err = rpc.New(&s, []byte(walletpwd)); err != nil {
		t.Error(err)
	}
	l2 = rpc.Run(&s)
	rpc.GoNotify(&s, node.RegisterTxNotifier,
		func(ch chan []tx.Hash) {
			cnotify = ch
		})
}

func teardown(t *testing.T) {
	if err := os.RemoveAll("./test_db"); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(tdir); err != nil {
		t.Log(err)
	}
	if err := l.Close(); err != nil {
		t.Log(err)
	}
	if err := l2.Close(); err != nil {
		t.Log(err)
	}
}

var (
	ledger *consensus.Ledger
)

func confirmAll(t *testing.T, confirm bool) {
	var txs []tx.Hash
	ledger = &consensus.Ledger{
		ParentID:  consensus.GenesisID,
		Seq:       1,
		CloseTime: time.Now(),
	}
	id := ledger.ID()
	t.Log("ledger id", hex.EncodeToString(id[:]))
	if err := akconsensus.PutLedger(&s, ledger); err != nil {
		t.Fatal(err)
	}
	err := s.DB.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte{byte(db.HeaderTxInfo)}); it.ValidForPrefix([]byte{byte(db.HeaderTxInfo)}); it.Next() {
			var dat []byte
			err2 := it.Item().Value(func(d []byte) error {
				dat = make([]byte, len(d))
				copy(dat, d)
				return nil
			})
			if err2 != nil {
				return err2
			}
			var ti imesh.TxInfo
			if err := arypack.Unmarshal(dat, &ti); err != nil {
				return err
			}
			if confirm && ti.StatNo != imesh.StatusGenesis {
				ti.StatNo = imesh.StatNo(id)
			}
			if !confirm && ti.StatNo != imesh.StatusGenesis {
				ti.StatNo = imesh.StatusPending
			}
			h := it.Item().Key()[1:]
			if err := db.Put(txn, h, &ti, db.HeaderTxInfo); err != nil {
				return err
			}
			if !bytes.Equal(h, genesis) {
				txs = append(txs, h)
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	select {
	case cnotify <- txs:
	case <-time.After(time.Second):
		t.Fatal("failed to notify")
	}
}
func TestRPCClient2(t *testing.T) {
	setup(t)
	defer teardown(t)
	ipport := fmt.Sprintf("http://%s:%d", s.RPCBind, s.RPCPort)
	cl := rpcc.New(ipport, s.RPCUser, s.RPCPassword, nil)

	tr := tx.New(s.Config, genesis)
	tr.AddInput(genesis, 0)
	if err := tr.AddMultisigOut(s.Config, 1, aklib.ADKSupply, a.Address58(s.Config), b.Address58(s.Config)); err != nil {
		t.Error(err)
	}
	if err := tr.Sign(a); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
		t.Error(err)
	}
	_, err := cl.SendRawTX(tr, tx.TypeNormal)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(6 * time.Second)
	mout := tr.MultiSigOuts[0]

	msig, err := cl.GetMultisigInfo(mout.Address(s.Config))
	if err != nil {
		t.Error(err)
	}
	if msig.M != 1 {
		t.Error("invalid msig")
	}
	if len(msig.Addresses) != 2 {
		t.Error("invalid msig")
	}
	if bytes.Equal(msig.Addresses[0], a.Address(s.Config)) && bytes.Equal(msig.Addresses[1], b.Address(s.Config)) ||
		bytes.Equal(msig.Addresses[1], a.Address(s.Config)) && bytes.Equal(msig.Addresses[0], b.Address(s.Config)) {
	} else {
		t.Error("invalid msig")
	}
}
func TestRPCClient(t *testing.T) {
	setup(t)
	defer teardown(t)
	ipport := fmt.Sprintf("http://%s:%d", s.RPCBind, s.RPCPort)
	cl := rpcc.New(ipport, "hoehoe", s.RPCPassword, nil)
	_, err := cl.SetTxFee(0.01)
	if err == nil {
		t.Error("should be error")
	}
	time.Sleep(3 * time.Second)
	cl = rpcc.New(ipport, s.RPCUser, s.RPCPassword, nil)
	_, err = cl.SetTxFee(0.01)
	if err != nil {
		t.Error(err)
	}

	adr, err := cl.GetNewAddress("")
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(adr, "AKADR") {
		t.Error("invalid address")
	}
	t.Log("wallet address", adr)

	tr := tx.New(s.Config, genesis)
	tr.AddInput(genesis, 0)
	if err = tr.AddOutput(s.Config, adr, aklib.ADK); err != nil {
		t.Error(err)
	}
	if err = tr.AddOutput(s.Config, a.Address58(s.Config), aklib.ADKSupply-aklib.ADK); err != nil {
		t.Error(err)
	}
	if err = tr.Sign(a); err != nil {
		t.Error(err)
	}
	if err = tr.PoW(); err != nil {
		t.Error(err)
	}
	id, err2 := cl.SendRawTX(tr, tx.TypeNormal)
	if err2 != nil {
		t.Error(err2)
	}
	if id != tr.Hash().String() {
		t.Error("invalid txid")
	}
	time.Sleep(2 * time.Second)
	confirmAll(t, true)
	time.Sleep(2 * time.Second)
	bal, err2 := cl.GetBalance("")
	if err2 != nil {
		t.Error(err2)
	}
	if bal != 1 {
		t.Fatal("invlid balance", bal)
	}

	err = cl.WalletPassphrase(walletpwd, 1000)
	if err != nil {
		t.Error(err)
	}

	id, err = cl.SendMany("", map[string]float64{
		b.Address58(s.Config): 0.1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 64 {
		t.Error("invalid txid")
	}
	time.Sleep(2 * time.Second)
	confirmAll(t, true)
	time.Sleep(2 * time.Second)
	inf, err2 := cl.ValidateAddress(adr)
	if err2 != nil {
		t.Error(err2)
	}
	if *inf.Account != "" {
		t.Error("invalid account")
	}
	if inf.Address != adr {
		t.Error("invalid address")
	}
	if !inf.IsMine {
		t.Error("invalid ismine")
	}
	if !inf.IsValid {
		t.Error("invalid isvalid")
	}
	bal, err = cl.GetBalance("")
	if err != nil {
		t.Error(err)
	}
	if bal != 1-0.1 {
		t.Fatal("invlid balance", bal)
	}
	gettx, err2 := cl.GetTransaction(id)
	if err2 != nil {
		t.Error(err2)
	}
	if gettx.Amount != -0.1 {
		t.Error("invalid amount")
	}
	if gettx.Confirmations != 100000 {
		t.Error("invalid confirmation")
	}
	if time.Now().Unix()-gettx.Time > 10*60 || time.Now().Unix() < gettx.Time {
		t.Error("invalid time")
	}
	if time.Now().Unix()-gettx.TimeReceived > 10*60 || time.Now().Unix() < gettx.TimeReceived {
		t.Error("invalid recvtime")
	}
	if gettx.Txid != id {
		t.Error("invalid id")
	}
	ps, err := cl.ListPeer()
	if err != nil {
		t.Error(err)
	}
	if len(ps) != 0 {
		t.Error("invalid listpeer")
	}
	pk, err := cl.DumpPrivKey()
	if err != nil {
		t.Error(err)
	}
	if !strings.HasPrefix(pk, "AKPRIV") {
		t.Error("invalid privkey")
	}
	bs, err := cl.ListBanned()
	if err != nil {
		t.Error(err)
	}
	if len(bs) != 0 {
		t.Error("invalid ListBanned")
	}
	if err = cl.DumpWallet(tdir); err != nil {
		t.Error(err)
	}
	if err = cl.ImportWallet(tdir); err != nil {
		t.Error(err)
	}
	ni, err := cl.GetNodeinfo()
	if err != nil {
		t.Error(err)
	}
	if *ni.Balance != uint64((1-0.1)*float64(aklib.ADK)) {
		t.Error("invalid balance", *ni.Balance)
	}
	if ni.Connections != 0 {
		t.Error("invalid connection")
	}
	if ni.Leaves != 1 {
		t.Error("invalid #leaves")
	}
	if ni.TxNo != 3 {
		t.Error("invalid txno", ni.TxNo)
	}
	ls, err := cl.GetLeaves()
	if err != nil {
		t.Error(err)
	}
	if len(ls) != 1 {
		t.Error("invliad getkeaves")
	}
	if ls[0] != id {
		t.Error("invalid leave")
	}
	hs, err := cl.GetLastHistory(adr)
	if err != nil {
		t.Error(err)
	}
	if len(hs) != 1 {
		t.Error("invalid history", len(hs))
	}
	if hs[0].Type != tx.TypeIn {
		t.Error("invalid history")
	}
	if hs[0].Index != 0 {
		t.Error("invalid index")
	}
	if hs[0].Hash != id {
		t.Error("invalid hash", hs[0].Hash)
	}
	tr2, err := cl.GetRawTx(id)
	if err != nil {
		t.Error(err)
	}
	if tr2.Hash().String() != id {
		t.Error("invalid getrawtx")
	}

	mfee := tx.NewMinableFee(s.Config, genesis)
	mfee.AddInput(tr.Hash(), 1)
	mfee.AddOutput(s.Config, a.Address58(s.Config), aklib.ADKSupply-100000000-10)
	mfee.AddOutput(s.Config, "", 10)
	if err = mfee.Sign(a); err != nil {
		t.Error(err)
	}
	_, err = cl.SendRawTX(mfee, tx.TypeRewardFee)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	fe, err := cl.GetMinableFeeTx(10.0 / aklib.ADK)
	if err != nil {
		t.Error(err)
	}
	if fe.Hash().String() != mfee.Hash().String() {
		t.Error("invalid feetx")
	}

	ti, err := tx.IssueTicket(s.Config, a.Address(s.Config), genesis)
	if err != nil {
		t.Error(err)
	}
	_, err = cl.SendRawTX(ti, tx.TypeNormal)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	confirmAll(t, true)
	time.Sleep(2 * time.Second)

	mtic := tx.NewMinableTicket(s.Config, ti.Hash(), genesis)
	mtic.AddInput(genesis, 0)
	if err = mtic.AddOutput(s.Config, a.Address58(s.Config), aklib.ADKSupply); err != nil {
		t.Error(err)
	}
	if err = mtic.Sign(a); err != nil {
		t.Error(err)
	}
	_, err = cl.SendRawTX(mtic, tx.TypeRewardTicket)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	tic, err := cl.GetMinableTicketTx()
	if err != nil {
		t.Error(err)
	}
	if tic.Hash().String() != mtic.Hash().String() {
		t.Error("invalid tickettx")
	}

	stats, err := cl.GetTxsStatus(id)
	if err != nil {
		t.Error(err)
	}
	if len(stats) != 1 {
		t.Error("invalid length")
	}
	if !stats[0].IsAccepted() {
		t.Error("invalid tx stat")
	}

	sfid, err := cl.SendFrom("", b.Address58(s.Config), 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sfid) != 64 {
		t.Error("invalid id")
	}
	time.Sleep(2 * time.Second)
	confirmAll(t, true)
	time.Sleep(2 * time.Second)

	id, err = cl.SendToAddress(b.Address58(s.Config), 0.1)
	if err != nil {
		t.Error(err)
	}
	if len(id) != 64 {
		t.Error("invalid id")
	}
	time.Sleep(2 * time.Second)
	confirmAll(t, true)
	time.Sleep(2 * time.Second)
	acs, err := cl.ListAddressGroupings()
	if err != nil {
		t.Error(err)
	}
	if len(acs) != 4 {
		t.Error("invalid length of listaddressgroupings", len(acs))
	}
	nadr := 0
	namt := 0
	for _, ac := range acs {
		if !strings.HasPrefix(ac.Address, "AKADRS") {
			t.Error("invalid account")
		}
		t.Log(ac.Account, ac.Address, ac.Amount)
		if ac.Address == adr {
			nadr++
			if *ac.Account != "" {
				t.Error("invalid account")
			}
			if ac.Amount != 0 {
				t.Error("invalid amount")
			}
		} else {
			if ac.Account != nil {
				t.Error("invalid account")
			}
			if ac.Amount == 1-0.3 {
				namt++
			} else {
				if ac.Amount != 0 {
					t.Error("invalid amount")
				}
			}
		}
	}
	if nadr != 1 {
		t.Error("invalid addess")
	}
	if namt != 1 {
		t.Error("invalid amount", namt)
	}
	ac, err := cl.ListAccounts()
	if err != nil {
		t.Error(err)
	}
	if len(ac) != 1 {
		t.Error("invalid listaccounts")
	}
	a, ok := ac[""]
	if !ok {
		t.Error("invalid account")
	}
	if a != 1-0.3 {
		t.Error("invliad amount")
	}
	acc, err := cl.GetAccount(adr)
	if err != nil {
		t.Error(err)
	}
	if acc != "" {
		t.Error("invlaid getaccount")
	}
	trs, err := cl.ListTransactions("", 1, 3)
	if err != nil {
		t.Error(err)
	}
	if len(trs) != 1 {
		t.Error("invalid listtx", len(trs))
	}
	t.Log(trs[0])
	if *trs[0].Account != "" {
		t.Error("invalid account")
	}
	if !strings.HasPrefix(trs[0].Address, "AKADR") {
		t.Error("invalid address")
	}
	if trs[0].Amount != 0.8 {
		t.Error("invalid amount", trs[0].Amount)
	}
	if trs[0].Confirmations != 100000 {
		t.Error("invalid confirmation")
	}
	if trs[0].Txid != sfid {
		t.Error("invalid txid", trs[0].Txid)
	}
	if trs[0].Vout != 1 {
		t.Error("invalid vout", trs[0].Vout)
	}
	if trs[0].Category != "receive" {
		t.Error("invalid category")
	}

	if err = cl.WalletLock(); err != nil {
		t.Error(err)
	}

	lid := ledger.ID()
	l, err := cl.GetLedger(hex.EncodeToString(lid[:]))
	if err != nil {
		t.Error(err)
	}
	if l.ID != hex.EncodeToString(lid[:]) || l.Seq != 1 {
		t.Error("invalid ledger")
	}
}
