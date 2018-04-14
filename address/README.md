[![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/address?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/address)

# address

## Overview

This  library is for generating addresses, marshaling/unmarshalling of addresses, signing and verifying by addresses for AidosKuneen.

## Example
(This example omits error handlings for simplicity.)

```go
	import "github.com/AidosKuneen/aklib/address"
	import "github.com/AidosKuneen/xmss"
	
	pwd := []byte("some password")
	seed := address.GenerateSeed()
	adr1 := address.New(address.Height10, seed, MainConfig)
	seed58 := adr1.Seed58(pwd) //base58 encoded seed
	adr2, err := address.NewFrom58(seed58, pwd, MainConfig)
	//adr1 and adr2 should be same

	pk58 := adr1.PK58() //base58 encoded public key
	pk, err := address.FromPK58(pk58,  MainConfig)

	msg := []byte("This is a test for XMSS.")
	sig := adr1.Sign(msg)	
	if !xmss.Verify(sig, msg, pk){
		...
	}

	b, err := aa.MarshalJSON()

	c := address.Address{}
	err := json.Unmarshal(b,&c)
```

```go
	import "github.com/AidosKuneen/aklib/address/node"
	
	node1 := node.New(address.Height40, seed, MainConfig)
	seed58 := node1.Seed58() 
	adr2, err := node.NewFrom58(seed58, MainConfig)
	pk58 := node1.PK58() //base58 encoded public key
	pk, err :=node.FromPK58(pk58,  MainConfig)
	sig := adr1.Sign(msg)	
	if !node.Verify(sig, msg, pk){
		...
	}

	b, err := node1.MarshalJSON()

	c := node.Address{}
	err := json.Unmarshal(b,&c)

```

