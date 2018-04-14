// Copyright (c) 2017 Aidos Developer

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

//Config is settings for various parameters.
type Config struct {
	Easiness       uint32
	TicketEasiness uint32
	PrefixPriv     [][]byte
	PrefixAdrs     [][]byte
	PrefixNkey     [][]byte
	PrefixNode     [][]byte
}

var (
	//MainConfig is a Config for MainNet
	MainConfig = &Config{
		Easiness:       0x0fffffff,
		TicketEasiness: 0x07ffffff,
		PrefixPriv: [][]byte{
			[]byte{0x16, 0x55, 0x35}, //VM1
			[]byte{0x16, 0x55, 0x72}, //VM5
			[]byte{0x16, 0x55, 0xa0}, //VM8
			[]byte{0x16, 0x55, 0xbf}, //VMA
		},
		PrefixAdrs: [][]byte{
			[]byte{0x82, 0xab, 0xc6}, //SM1
			[]byte{0x82, 0xad, 0x58}, //SM5
			[]byte{0x02, 0x40, 0xce}, //SM8
			[]byte{0x02, 0x40, 0xd1}, //SMA
		},
		PrefixNkey: [][]byte{
			[]byte{0x02, 0xc9, 0x4b}, //YM1
			[]byte{0x02, 0xc9, 0x4d}, //YM2
			[]byte{0x02, 0xc9, 0x4f}, //YM3
		},
		PrefixNode: [][]byte{
			[]byte{0x01, 0x2f, 0xae}, //EM1
			[]byte{0x01, 0x2f, 0xb0}, //EM2
			[]byte{0x01, 0x2f, 0xb2}, //EM3
		},
	}
	//TestConfig is a Config for TestNet
	TestConfig = &Config{
		Easiness:       0x7fffffff,
		TicketEasiness: 0x3fffffff,
		PrefixPriv: [][]byte{
			[]byte{0x16, 0x6a, 0x12}, //VT1
			[]byte{0x16, 0x6a, 0x50}, //VT5
			[]byte{0x16, 0x6a, 0x7e}, //VT8
			[]byte{0x16, 0x6a, 0x9d}, //VTA
		},
		PrefixAdrs: [][]byte{
			[]byte{0x02, 0x43, 0x1c}, //ST1
			[]byte{0x02, 0x43, 0x23}, //ST5
			[]byte{0x02, 0x43, 0x28}, //ST8
			[]byte{0x02, 0x43, 0x2c}, //STA
		},
		PrefixNkey: [][]byte{
			[]byte{0x02, 0xcb, 0xa6}, //YT1
			[]byte{0xa2, 0x23, 0xe6}, //YT2
			[]byte{0x02, 0xcb, 0xa9}, //YT3
		},
		PrefixNode: [][]byte{
			[]byte{0x01, 0x32, 0x09}, //ET1
			[]byte{0x45, 0x56, 0x50}, //ET2
			[]byte{0x01, 0x32, 0x0c}, //ET3
		},
	}
)

const (
	//ADKMinUnit is minimum unit of ADK.
	ADKMinUnit float64 = 0.00000001
	//OneADK is for converting 1 ADK to unit in transactions.
	OneADK = uint64(1 / ADKMinUnit)
	//ADKSupply is total supply of ADK.
	ADKSupply = 25 * 1000000 * OneADK
)
