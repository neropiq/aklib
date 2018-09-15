[![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/address?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/address)

# address

## Overview

This  library is for generating addresses, marshaling/unmarshalling of addresses, signing and verifying by addresses for AidosKuneen.

## Example
(This example omits error handlings for simplicity.)

```go
	import "github.com/AidosKuneen/aklib/address"
	
	pwd := []byte("some password")
	seed := address.GenerateSeed32()
	adr1 := address.New(MainConfig, seed)
	seed58 := HDSeed58(MainConfig,seed,pwd,false) //base58 encoded seed

	//For node
	// adr1 := address.NewNode(MainConfig, seed)
	// seed58 := HDSeed58(MainConfig,seed,pwd,true) //base58 encoded seed

	seed581,_,err := HDFrom58(MainConfig,seed58,pwd) //base58 encoded seed
	//seed58 and seed581 should be same

	pk58 := adr1.Address58(MainConfig) //base58 encoded public key
	pk, err := address.ParseAddress58(  MainConfig,pk58)

	msg := []byte("This is a test.")
	sig := adr1.Sign(msg)	
	err := address.Verify(sig,msg)	
```


