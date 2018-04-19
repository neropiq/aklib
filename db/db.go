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
	"os"

	"github.com/dgraph-io/badger"
)

//Headers for DB key.
const (
	HeaderTxBody byte = iota
	HeaderTxSig
	HeaderBalance
	HeaderNodeIP
	HeaderNodeSK
	HeaderStatements
	HeaderWalletSK
)

const dbDir = "./db"

//New open badger db.
func New(dir string) (*badger.DB, error) {
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0755); err != nil {
			return nil, err
		}
	}
	opts := badger.DefaultOptions
	opts.SyncWrites = false
	opts.Dir = dir
	opts.ValueDir = dir
	return badger.Open(opts)
}

//Copy copies db to toDir.
func Copy(db *badger.DB, todir string) {
	db2, err := New(todir)
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
	db2.Close()
}
