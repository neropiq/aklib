[![Build Status](https://travis-ci.org/AidosKuneen/aklib.svg?branch=master)](https://travis-ci.org/AidosKuneen/aklib)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/aklib/LICENSE)

aklib : [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib)

address : [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/address?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/address)


tx : [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/tx?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/tx)

# aklib 

## Overview

This  library is for managing AidosKuneen features including

* Address
* Transaction  
* Proof of Work (in go and C)
* API call (not yet)
* Transfer (not yet)

## Requirements

* git
* go 1.9+
* yasm

are required to compile this.

## Installation

     $ go get github.com/AidosKuneen/aklib


## Example
(This example omits error handlings for simplicity.)

### address

```go
	import "github.com/AidosKuneen/aklib/address"
	
	seed := GenerateSeed()
	a := New(10, seed, MainConfig)
	s58 := a.Seed58()
	aa, err := NewFrom58(10, s58, MainConfig)

	pk58 := a.PK58()
	pk, err := FromPK58(pk58, MainConfig)

	msg := []byte("This is a test for XMSS.")
	sig := aa.Sign(msg)	
	if !Verify(sig, msg, pk){
		...
	}

	b, err := aa.MarshalJSON()

	c := Address{}
	err := c.UnmarshalJSON(b)
```
### transaction
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


# Contribution
Improvements to the codebase and pull requests are encouraged.


