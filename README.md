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
* Proof of Work with [Cuckoo Cycle](https://github.com/AidosKuneen/cuckoo)
* API call (not yet)
* Transfer (not yet)

## Requirements

* git
* go 1.9+

are required to compile this.

## Installation

     $ go get github.com/AidosKuneen/aklib


## Example
(This example omits error handlings for simplicity.)

### address

```go
	import "github.com/AidosKuneen/aklib/address"
	
	seed := GenerateSeed()
	adr1 := New(address.Height10, seed, MainConfig)
	seed58 := adr1.Seed58() //base58 encoded seed
	adr2, err := NewFrom58(seed58, MainConfig)
	//adr1 and adr2 should be same

	pk58 := adr1.PK58() //base58 encoded public key
	pk, err := FromPK58(pk58, MainConfig)

	msg := []byte("This is a test for XMSS.")
	sig := adr1.Sign(msg)	
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


