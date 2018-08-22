[![Build Status](https://travis-ci.org/AidosKuneen/aklib.svg?branch=master)](https://travis-ci.org/AidosKuneen/aklib)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/AidosKuneen/aklib/LICENSE)
[![GoDoc](https://godoc.org/github.com/AidosKuneen/aklib?status.svg)](https://godoc.org/github.com/AidosKuneen/aklib)
[![Coverage Status](https://coveralls.io/repos/github/AidosKuneen/aklib/badge.svg?branch=master)](https://coveralls.io/github/AidosKuneen/aklib?branch=master)

# aklib 

## Overview

This  library is for managing AidosKuneen features including

* [Address](https://github.com/AidosKuneen/aklib/tree/master/address)
* [Transaction, Proof of Work](https://github.com/AidosKuneen/aklib/tree/master/tx)
 (Proof of Work with [Cuckoo Cycle](https://github.com/AidosKuneen/cuckoo))
* [RPC client](https://github.com/AidosKuneen/aklib/tree/master/rpcclient)

## Requirements

* git
* go 1.9+

are required to compile this.

## Installation

     $ go get github.com/AidosKuneen/aklib


## Contribution
Improvements to the codebase and pull requests are encouraged.


## Dependencies and Licenses

This software includes the work that is distributed in the Apache License 2.0.

```
github.com/AidosKuneen/aklib          MIT License
github.com/AidosKuneen/aknode         MIT License
github.com/AidosKuneen/cuckoo         MIT License
github.com/AidosKuneen/numcpu         MIT License
github.com/AidosKuneen/sha256-simd    Apache License 2.0
github.com/AidosKuneen/xmss           MIT License
github.com/AndreasBriese/bbloom       MIT License, Public Domain
github.com/dgraph-io/badger           Apache License 2.0
github.com/dgryski/go-farm            MIT License
github.com/golang/protobuf/proto      BSD 3-clause "New" or "Revised" License
github.com/mattn/go-shellwords        MIT License
github.com/pkg/errors                 BSD 2-clause "Simplified" License
github.com/vmihailenco/msgpack/codes  BSD 2-clause "Simplified" License
golang.org/x/net                      BSD 3-clause "New" or "Revised" License
golang.org/x/sys/unix                 BSD 3-clause "New" or "Revised" License
Golang Standard Library               BSD 3-clause License
```