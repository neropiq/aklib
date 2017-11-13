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
	PrefixPriv []byte
	PrefixAdrs []byte
}

var (
	//MainConfig is a Config for MainNet
	MainConfig = &Config{
		Difficulty: 2,
		PrefixPriv: []byte{0xbf, 0x9d}, //"VM" in base58
		PrefixAdrs: []byte{0xab, 0x55}, //"SM" in base58
	}
	//TestConfig is a Config for TestNet
	TestConfig = &Config{
		Difficulty: 1,
		PrefixPriv: []byte{0xc0, 0x50}, //"VT" in base58
		PrefixAdrs: []byte{0xac, 0x8},  //"ST" in base58
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
