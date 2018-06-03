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
	"log"
	"os"
	"time"

	"github.com/AidosKuneen/aklib/arypack"
	"github.com/dgraph-io/badger"
)

//Headers for DB key.
const (
	HeaderTxBody byte = iota + 1
	HeaderTxSig
	HeaderBalance
	HeaderNodeIP
	HeaderNodeSK
	HeaderWalletSK
	HeaderFeeTx
	HeaderTicketTx
	HeaderLeaves
	HeaderUnresolvedTx
	HeaderUnresolvedInfo
	HeaderBrokenTx
)

const dbDir = "./db"

//Open open  or make a badger db.
func Open(dir string) (*badger.DB, error) {
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}
	opts := badger.DefaultOptions
	opts.SyncWrites = false
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			time.Sleep(time.Hour)
			if err := db.RunValueLogGC(0.5); err != nil {
				log.Println(err)
			}
		}
	}()
	return db, nil
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
				v, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				if err := txn2.Set(k, v); err != nil {
					return err
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
func Key(okey []byte, h byte) []byte {
	return append([]byte{h}, okey...)
}

//Get get a data from DB.
func Get(txn *badger.Txn, key []byte, dat interface{}, header byte) error {
	item, err := txn.Get(append([]byte{header}, key...))
	if err != nil {
		return err
	}
	val, err := item.ValueCopy(nil)
	if err != badger.ErrConflict && err != nil {
		return err
	}
	return arypack.Unmarshal(val, dat)
}

//Put puts a dat into db.
func Put(txn *badger.Txn, key []byte, dat interface{}, header byte) error {
	err := txn.Set(append([]byte{header}, key...), arypack.Marshal(dat))
	if err == badger.ErrConflict {
		return nil
	}
	return err
}

//Del deletes a dat from db.
func Del(txn *badger.Txn, key []byte, header byte) error {
	return txn.Delete(append([]byte{header}, key...))
}
