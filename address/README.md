[![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib/address?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib/address)

# address

## Overview

This  library is for generating addresses, marshaling/unmarshalling of addresses, signing and verifying by addresses for AidosKuneen.

## Example
(This example omits error handlings for simplicity.)

```go
	import "github.com/AidosKuneen/aklib/address"
	
	pwd := []byte("some password")
	seed := GenerateSeed(pwd)
	adr1 := New(address.Height10, seed, MainConfig)
	seed58 := adr1.Seed58() //base58 encoded seed
	adr2, err := NewFrom58(seed58, pwd, MainConfig)
	//adr1 and adr2 should be same

	pk58 := adr1.PK58() //base58 encoded public key
	pk, err := FromPK58(pk58, pwd, MainConfig)

	msg := []byte("This is a test for XMSS.")
	sig := adr1.Sign(msg)	
	if !Verify(sig, msg, pk){
		...
	}

	b, err := aa.MarshalJSON()

	c := Address{}
	err := c.UnmarshalJSON(b)
```

