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

package db

import (
	"context"
	"os"
	"time"

	"github.com/AidosKuneen/aklib/arypack"
	"github.com/dgraph-io/badger"
)

//Header is a header type of db.
type Header byte

//Headers for DB key.
const (
	HeaderTxInfo Header = iota + 1
	HeaderTxSig
	HeaderTxNo
	HeaderNodeIP
	HeaderTxRewardFee
	HeaderTxRewardTicket
	HeaderLeaves
	HeaderUnresolvedTx
	HeaderUnresolvedInfo
	HeaderBrokenTx
	HeaderAddressToTx
	HeaderMultisigAddress
	HeaderWallet
	HeaderWalletAddress
	HeaderWalletHistory
	HeaderWalletConfig

	HeaderLedger
	HeaderLastLedger
)

//Open open  or make a badger db.
func Open(dir string) (*badger.DB, error) {
	if _, err := os.Stat(dir); err != nil {
		if err2 := os.Mkdir(dir, 0755); err2 != nil {
			return nil, err2
		}
	}
	opts := badger.DefaultOptions
	opts.SyncWrites = false
	opts.Dir = dir
	opts.ValueDir = dir
	return badger.Open(opts)
}

//GoGC runs gc for badger DB.
func GoGC(ctx context.Context, db *badger.DB) {
	go func() {
		ctx2, cancel2 := context.WithCancel(ctx)
		defer cancel2()
		for {
			select {
			case <-ctx2.Done():
				return
			case <-time.After(5 * time.Minute):
			again:
				if err := db.RunValueLogGC(0.7); err == nil {
					goto again
				}
			}
		}
	}()
}

//Copy copies db to toDir.
func Copy(db *badger.DB, todir string) {
	db2, err := Open(todir)
	if err != nil {
		panic(err)
	}
	err = db2.Update(func(txn2 *badger.Txn) error {
		return db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchSize = 10
			it := txn.NewIterator(opts)
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				v, err2 := item.ValueCopy(nil)
				if err2 != nil {
					return err2
				}
				if err2 := txn2.Set(k, v); err2 != nil {
					return err2
				}
			}
			return nil
		})
	})
	if err != nil {
		panic(err)
	}
	if err := db2.Close(); err != nil {
		panic(err)
	}
}

//Key return h+okey
func Key(okey []byte, h Header) []byte {
	return append([]byte{byte(h)}, okey...)
}

//Get get a data from DB.
func Get(txn *badger.Txn, key []byte, dat interface{}, header Header) error {
	item, err := txn.Get(append([]byte{byte(header)}, key...))
	if err != nil {
		return err
	}
	val, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}
	return arypack.Unmarshal(val, dat)
}

//Put puts a dat into db.
func Put(txn *badger.Txn, key []byte, dat interface{}, header Header) error {
	return txn.Set(append([]byte{byte(header)}, key...), arypack.Marshal(dat))
}

//Del deletes a dat from db.
func Del(txn *badger.Txn, key []byte, header Header) error {
	return txn.Delete(append([]byte{byte(header)}, key...))
}
