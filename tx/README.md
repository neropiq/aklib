 [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/tx?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/tx)

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

	tx:=tx.New()
	if err := tx.Check(MainConfig); err == nil {
	}
	if err := tx.CheckAll(GetTxFunction, MainConfig); err == nil {
	}
	if err := tx.PoW(); err != nil {
	}
	h = tx.Hash()
	b := tx.Body.Pack()
	s := tx.Signatures.Pack()
```

# Expected Time for PoW

The expecting time of suceeding to find a solution of cuckoo is around 20 seconds and 3-sigma
is 60 seconds for general low-end PCs, and this lib is configured to do 2^2=4 times in average 
to meet hash-based PoW now.
And generally high-end servers need half of this time.

So the estimation of expecting PoW time
* for low-end PCs is 80 seconds in average and 3-sigma is 240 seconds.
* for high-end servers is 40 seconds in average and 3-sigma is 120 seconds.
