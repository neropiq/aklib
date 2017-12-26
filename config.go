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
	Difficulty byte
	PrefixPriv [][]byte
	PrefixAdrs [][]byte
}

var (
	//MainConfig is a Config for MainNet
	MainConfig = &Config{
		Difficulty: 4,
		PrefixPriv: [][]byte{
			[]byte{0x06, 0xa2, 0x71}, //"VM5" in base5
			[]byte{0x06, 0xa2, 0x7f}, //"VM8" in base5
			[]byte{0x06, 0xa2, 0x87}, //"VMA" in base5
		},
		PrefixAdrs: [][]byte{
			[]byte{0x26, 0xd1, 0xb8}, //"SM5" in base5
			[]byte{0x26, 0xd2, 0x12}, //"SM8" in base5
			[]byte{0x26, 0xd2, 0x4d}, //"SMA" in base5
		},
	}
	//TestConfig is a Config for TestNet
	TestConfig = &Config{
		Difficulty: 1,
		PrefixPriv: [][]byte{
			[]byte{0x06, 0xa8, 0xa3}, //"VT5" in base5
			[]byte{0x06, 0xa8, 0xb0}, //"VT8" in base5
			[]byte{0x06, 0xa8, 0xba}, //"VTA" in base5
		},
		PrefixAdrs: [][]byte{
			[]byte{0x26, 0xfa, 0x48}, //"ST5" in base5
			[]byte{0x26, 0xfa, 0xa1}, //"ST8" in base5
			[]byte{0x26, 0xfa, 0xdd}, //"STA" in base5
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
