[![Build Status](https://travis-ci.org/AidosKuneen/aklib.svg?branch=master)](https://travis-ci.org/AidosKuneen/aklib)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/aklib/LICENSE)

aklib : [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib)

address : [![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/address?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/address)


# aklib 

## Overview

This  library is for managing AidosKuneen features including

* Address
* Transaction  (not yet)
* Proof of Work (in C) (not yet)
* API call (not yet)
* Transfer (not yet)

## Requirements

This requires

* git
* go 1.9+


## Installation

     $ go get github.com/AidosKuneen/aklib


## Example
(This example omits error handlings for simplicity.)

```go
	import "github.com/AidosKuneen/aklib"
	
	seed := GenerateSeed()
	a := New(10, seed, net)
	s58 := a.Seed58()
	aa, err := NewFrom58(10, s58, net)

	pk58 := a.PK58()
	pk, err := FromPK58(pk58, net)

	msg := []byte("This is a test for XMSS.")
	sig := aa.Sign(msg)	
	if !Verify(sig, msg, pk){
		...
	}

	b, err := aa.MarshalJSON()

	c := Address{}
	err := c.UnmarshalJSON(b)
```


# Contribution
Improvements to the codebase and pull requests are encouraged.


