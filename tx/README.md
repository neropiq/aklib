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


## Performance of PoW

By using a (cheap) linux PC:

```
* Compiler: go version go1.9.1 linux/amd64
* Kernel: Linux 4.14.5-1-ARCH #1 SMP PREEMPT Sun Dec 10 14:50:30 UTC 2017 x86_64 GNU/Linux
* CPU:  Celeron(R) CPU G1840 @ 2.80GHz 
* Memory: 8 GB
```

It takes about 5 minutes.

```
BenchmarkPoW-2   	       1	304205459533 ns/op
--- BENCH: BenchmarkPoW-2
PASS
ok  	github.com/AidosKuneen/aklib/tx	304.214s
```

By using DIGNO M KYL22(Android Smartphone):

```
* Compiler: go version go1.9.1 linux/arm
* OS: 	Android 4.2.2
* CPU:	Qualcomm Snapdragon 800 MSM8974 2.2GHz (quad core)
* Memory: 2 GB
```

It takes about 12 minutes.

```
BenchmarkPoW 	       1	725203681289 ns/op
```

By using a cloud server:

* Compiler: go version go1.8.1 linux/amd64
* Kernel: Linux 4.8.0-58-generic #63~16.04.1-Ubuntu SMP Mon Jun 26 18:08:51 UTC 2017 x86_64 x86_64 x86_64 GNU/Linux
* CPU:  CAMD Ryzen 7 1700X Eight-Core Processor @ 2.20GHz (8 cores)
* Memory: 64 GB


It takes about 4 minutes.

```
BenchmarkPoW-16    	       1	246550764640 ns/op
```