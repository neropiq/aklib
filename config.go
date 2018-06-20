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

const (
	//ADK is for converting 1 ADK to unit in transactions.
	ADK = 100000000
	//ADKSupply is total supply of ADK.
	ADKSupply = 25 * 1000000 * ADK
)

//SRV is a param for SRVLookup
type SRV struct {
	Service string
	Name    string
}

//Config is settings for various parameters.
type Config struct {
	Name           string
	Easiness       uint32
	TicketEasiness uint32
	PrefixPriv     [][]byte
	PrefixAdrs     [][]byte
	PrefixNkey     [][]byte
	PrefixNode     [][]byte

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
		Name:           "mainnet",
		Easiness:       0x3fffffff,
		TicketEasiness: 0x03ffffff,
		PrefixPriv: [][]byte{
			[]byte{0x06, 0xa2, 0x5e}, //"VM1" in base5
			[]byte{0x06, 0xa2, 0x71}, //"VM5" in base5
			[]byte{0x06, 0xa2, 0x7f}, //"VM8" in base5
			[]byte{0x06, 0xa2, 0x87}, //"VMA" in base5
		},
		PrefixAdrs: [][]byte{
			[]byte{0x26, 0xd1, 0x41}, //"SM1" in base5
			[]byte{0x26, 0xd1, 0xb8}, //"SM5" in base5
			[]byte{0x26, 0xd2, 0x12}, //"SM8" in base5
			[]byte{0x26, 0xd2, 0x4d}, //"SMA" in base5
		},
		PrefixNkey: [][]byte{
			[]byte{0x30, 0x01, 0xbf}, //YM1
			[]byte{0x30, 0x01, 0xdd}, //YM2
			[]byte{0x30, 0x01, 0xfb}, //YM3
		},
		PrefixNode: [][]byte{
			[]byte{0x14, 0x70, 0x45}, //EM1
			[]byte{0x14, 0x70, 0x63}, //EM2
			[]byte{0x14, 0x70, 0x81}, //EM3
		},
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
		Name:           "testnet",
		Easiness:       0x7fffffff,
		TicketEasiness: 0x3fffffff,
		PrefixPriv: [][]byte{
			[]byte{0x06, 0xa8, 0x91}, //"VT1" in base5
			[]byte{0x06, 0xa8, 0xa3}, //"VT5" in base5
			[]byte{0x06, 0xa8, 0xb0}, //"VT8" in base5
			[]byte{0x06, 0xa8, 0xba}, //"VTA" in base5
		},
		PrefixAdrs: [][]byte{
			[]byte{0x26, 0xf9, 0xd0}, //"ST1" in base5
			[]byte{0x26, 0xfa, 0x48}, //"ST5" in base5
			[]byte{0x26, 0xfa, 0xa1}, //"ST8" in base5
			[]byte{0x26, 0xfa, 0xdd}, //"STA" in base5
		},
		PrefixNkey: [][]byte{
			[]byte{0x30, 0x2a, 0x4e}, //YT1
			[]byte{0x30, 0x2a, 0x6c}, //YT2
			[]byte{0x30, 0x2a, 0x8a}, //YT3
		},
		PrefixNode: [][]byte{
			[]byte{0x14, 0x98, 0xd4}, //ET1
			[]byte{0x14, 0x98, 0xf2}, //ET2
			[]byte{0x14, 0x99, 0x10}, //ET3
		},
		DefaultPort:         14370,
		DefaultRPCPort:      14371,
		DefaultExplorerPort: 8080,
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
		Name:           "debug",
		Easiness:       0xffffffff,
		TicketEasiness: 0x7fffffff,
		PrefixPriv: [][]byte{
			[]byte{0x06, 0xa8, 0x91}, //"VT1" in base5
			[]byte{0x06, 0xa8, 0xa3}, //"VT5" in base5
			[]byte{0x06, 0xa8, 0xb0}, //"VT8" in base5
			[]byte{0x06, 0xa8, 0xba}, //"VTA" in base5
		},
		PrefixAdrs: [][]byte{
			[]byte{0x26, 0xf9, 0xd0}, //"ST1" in base5
			[]byte{0x26, 0xfa, 0x48}, //"ST5" in base5
			[]byte{0x26, 0xfa, 0xa1}, //"ST8" in base5
			[]byte{0x26, 0xfa, 0xdd}, //"STA" in base5
		},
		PrefixNkey: [][]byte{
			[]byte{0x30, 0x2a, 0x4e}, //YT1
			[]byte{0x30, 0x2a, 0x6c}, //YT2
			[]byte{0x30, 0x2a, 0x8a}, //YT3
		},
		PrefixNode: [][]byte{
			[]byte{0x14, 0x98, 0xd4}, //ET1
			[]byte{0x14, 0x98, 0xf2}, //ET2
			[]byte{0x14, 0x99, 0x10}, //ET3
		},
		DefaultPort:         14370,
		DefaultRPCPort:      14371,
		DefaultExplorerPort: 8080,
		DNS: []SRV{
			SRV{
				Service: "seeds",
				Name:    "aidoskuneen.com",
			},
		},
		MessageMagic: 0xD9CBA322,
		Genesis: map[string]uint64{
			"": ADKSupply,
		},
	}
)

//Configs is a list of nets.
var Configs = []*Config{MainConfig, TestConfig}
