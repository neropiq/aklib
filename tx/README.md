 [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/tx?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/tx)
[![Coverage Status](https://coveralls.io/repos/github/AidosKuneen/aklib/badge.svg?branch=master)](https://coveralls.io/github/AidosKuneen/aklib?branch=master)

# transaction

## Overview

This  library is for cheking and packing transactions, and doing Proof of Work(PoW).

## Requirements

* git
* go 1.9+

are required to compile this.

## Installation

     $ go get github.com/AidosKuneen/aklib


## Example
(This example omits error handlings for simplicity.)

```go
	import "github.com/AidosKuneen/aklib/tx"
	import "github.com/AidosKuneen/aklib"
	import "github.com/AidosKuneen/aklib/address"

	var tx1,tx2,tx3,tx4 tx.Hash
	var adr1,adr2,adr3 string 
	var adr address.Address
	tr:=tx.New(aklib.MainConfig,tx1,tx2)
	tr.AddInput(tx3, 0)
	tr.AddMultisigIn(tx4, 1)
	if err := tr.AddOutput(aklib.DebugConfig, adr1, 111); err != nil {
	}
	if err := tr.AddMultisigOut(aklib.DebugConfig, 2, 331,adr1, adr2, adr3);err!=nil{
	}
	if err := tr.Sign(adr); err != nil {
		t.Error(err)
	}
	if err := tr.PoW(); err != nil {
	}
	if err := tr.Check(MainConfig,tx.TxNormal); err == nil {
	}
	if err := tr.CheckAll(GetTxFunction, aklib.MainConfig,TxNormal); err == nil {
	}
	h = tr.Hash()
```

# Expected Time for PoW

The expecting time of suceeding to find a solution of cuckoo is around 20 seconds and 3-sigma
is 60 seconds for general low-end PCs, and this lib is configured to do 2^2=4 times in average 
to meet hash-based PoW now.
And generally high-end servers need half of this time.

So the estimation of expecting PoW time
* for low-end PCs is 80 seconds in average and 3-sigma is 240 seconds.
* for high-end servers is 40 seconds in average and 3-sigma is 120 seconds.
