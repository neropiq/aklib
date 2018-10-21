// Copyright (c) 2018 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package aklib

import (
	"github.com/dgraph-io/badger"
)

const (
	//ADK is for converting 1 ADK to unit in transactions.
	ADK = 100000000
	//ADKSupply is total supply of ADK.
	ADKSupply uint64 = 25 * 1000000 * ADK
)

//SRV is a param for SRVLookup
type SRV struct {
	Service string
	Name    string
}

//DBConfig is a set of config for db.
type DBConfig struct {
	DB     *badger.DB `json:"-"`
	Config *Config    `json:"-"`
}

//Config is settings for various parameters.
type Config struct {
	Name           string
	Easiness       uint32
	TicketEasiness uint32
	PrefixPriv     []byte
	PrefixAdrs     []byte
	PrefixMsig     []byte
	PrefixNkey     []byte
	PrefixNode     []byte

	DefaultPort         uint16
	DefaultRPCPort      uint16
	DefaultExplorerPort uint16

	DNS          []SRV
	MessageMagic uint32
	Genesis      map[string]uint64
}

var (
	//MainConfig is a Config for MainNet
	MainConfig = &Config{
		Name:                "mainnet",
		Easiness:            0x3fffffff,
		TicketEasiness:      0x03ffffff,
		PrefixPriv:          []byte{0x1d, 0x49}, //VM
		PrefixAdrs:          []byte{0xab, 0x55}, //SM
		PrefixMsig:          []byte{0x67, 0xbc}, //GM
		PrefixNkey:          []byte{0x20, 0x62}, //YM
		PrefixNode:          []byte{0x5a, 0x37}, //EM
		DefaultPort:         14270,
		DefaultRPCPort:      14271,
		DefaultExplorerPort: 8080,
		DNS: []SRV{
			SRV{
				Service: "",
				Name:    "",
			},
		},
		MessageMagic: 0xD3AB9E77,
		Genesis: map[string]uint64{
			"": ADKSupply,
		},
	}
	//TestConfig is a Config for TestNet
	TestConfig = &Config{
		Name:                "testnet",
		Easiness:            0x7fffffff,
		TicketEasiness:      0x3fffffff,
		PrefixPriv:          []byte{0x1d, 0x64}, //VT
		PrefixAdrs:          []byte{0xac, 0x08}, //ST
		PrefixMsig:          []byte{0x68, 0x6f}, //GT
		PrefixNkey:          []byte{0x20, 0x7d}, //YT
		PrefixNode:          []byte{0x5a, 0xea}, //ET
		DefaultPort:         14370,
		DefaultRPCPort:      14371,
		DefaultExplorerPort: 8081,
		DNS: []SRV{
			SRV{
				Service: "",
				Name:    "",
			},
		},
		MessageMagic: 0xD9CBA322,
		Genesis: map[string]uint64{
			"": ADKSupply,
		},
	}
	//DebugConfig is a Config for debug
	DebugConfig = &Config{
		Name:                "debug",
		Easiness:            0xffffffff,
		TicketEasiness:      0x7fffffff,
		PrefixPriv:          []byte{0x1d, 0x24}, //VD
		PrefixAdrs:          []byte{0xaa, 0x67}, //SD
		PrefixMsig:          []byte{0x66, 0xcd}, //GD
		PrefixNkey:          []byte{0x20, 0x3e}, //YD
		PrefixNode:          []byte{0x59, 0x48}, //ED
		DefaultPort:         14370,
		DefaultRPCPort:      14371,
		DefaultExplorerPort: 8081,
		DNS: []SRV{
			SRV{
				Service: "testnetseeds",
				Name:    "aidoskuneen.com",
			},
		},
		MessageMagic: 0xD9CBA322,
		Genesis: map[string]uint64{
			"": ADKSupply, //
		},
	}
)

//Configs is a list of nets.
var Configs = []*Config{MainConfig, TestConfig, DebugConfig}
