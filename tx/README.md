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

