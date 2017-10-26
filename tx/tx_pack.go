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

package tx

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//byte2Varint converts VarInt of Bitcoin to uint64.
func byte2Varint(dat *bytes.Buffer) (uint64, error) {
	var err error
	var bb byte
	if bb, err = dat.ReadByte(); err != nil {
		return 0, err
	}
	switch bb {
	case 0xfd:
		bs := make([]byte, 2)
		if _, err := dat.Read(bs); err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint16(bs)
		return uint64(v), nil
	case 0xfe:
		bs := make([]byte, 4)
		if _, err := dat.Read(bs); err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint32(bs)
		return uint64(v), nil
	case 0xff:
		bs := make([]byte, 8)
		if _, err := dat.Read(bs); err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint64(bs)
		return v, nil

	default:
		return uint64(bb), nil
	}
}

//int2Varint converts  uint64 to VarInt of Bitcoin.
func int2Varint(dat uint64) []byte {
	var b []byte
	switch {
	case dat < uint64(0xfd):
		b = make([]byte, 1)
		b[0] = byte(dat & 0xff)
	case dat <= uint64(0xffff):
		b = make([]byte, 3)
		b[0] = 0xfd
		binary.BigEndian.PutUint16(b[1:], uint16(dat))
	case dat <= uint64(0xffffffff):
		b = make([]byte, 5)
		b[0] = 0xfe
		binary.BigEndian.PutUint32(b[1:], uint32(dat))
	default:
		b = make([]byte, 9)
		b[0] = 0xff
		binary.BigEndian.PutUint64(b[1:], dat)
	}
	return b
}

func appendVarint(b []byte, v interface{}) []byte {
	var val uint64
	switch l := v.(type) {
	case int:
		val = uint64(l)
	case uint16:
		val = uint64(l)
	case uint32:
		val = uint64(l)
	case uint64:
		val = uint64(l)
	default:
		panic("v must be integer")
	}
	t := int2Varint(val)
	return append(b, t...)
}

//Pack return bytes slice of a transaction body.
func (bd *Body) Pack() []byte {
	b := make([]byte, 0, TransactionMax)

	b = append(b, bd.Type...)
	b = append(b, bd.Nonce...)
	b = appendVarint(b, bd.Time)
	b = appendVarint(b, len(bd.Message))
	b = append(b, bd.Message...)
	b = appendVarint(b, len(bd.Inputs))
	for _, inp := range bd.Inputs {
		b = append(b, inp.PreviousTX...)
		b = append(b, inp.Index)
	}
	b = appendVarint(b, len(bd.Outputs))
	for _, out := range bd.Outputs {
		b = append(b, out.Address...)
		b = appendVarint(b, out.Value)
	}
	b = appendVarint(b, len(bd.MultiSigs))
	for _, mul := range bd.MultiSigs {
		b = append(b, mul.N)
		b = appendVarint(b, len(mul.Addresses))
		for _, a := range mul.Addresses {
			b = append(b, a...)
		}
		b = appendVarint(b, mul.Value)
	}
	b = appendVarint(b, len(bd.Previous))
	for _, a := range bd.Previous {
		b = append(b, a...)
	}
	b = append(b, bd.Difficulty)
	b = appendVarint(b, bd.LockTime)

	//for avoiding extra space of heap
	bb := make([]byte, len(b))
	copy(bb, b)
	return bb
}

//Pack return bytes slice of a transaction signature..
func (sig Signatures) Pack() []byte {
	b := make([]byte, 0, TransactionMax)
	b = appendVarint(b, len(sig))
	for _, s := range sig {
		b = appendVarint(b, len(s))
		b = append(b, s...)
	}

	//for avoiding extra space of heap
	bb := make([]byte, len(b))
	copy(bb, b)
	return bb
}

//Pack returns byte slice of tx.
func (tx *Transaction) Pack() []byte {
	b1 := tx.Body.Pack()
	b2 := tx.Signatures.Pack()
	b := make([]byte, len(b1)+len(b2))
	copy(b, b1)
	copy(b[len(b1):], b2)
	return b
}

//UnpackSignature return a transaction signature from byte slice.
func UnpackSignature(dat []byte) (Signatures, error) {
	buf := bytes.NewBuffer(dat)
	n, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	ary := make(Signatures, n)
	for i := range ary {
		n, err := byte2Varint(buf)
		if err != nil {
			return nil, err
		}
		ary[i] = make([]byte, n)
		if _, err := buf.Read(ary[i]); err != nil {
			return nil, err
		}
	}
	if buf.Len() != 0 {
		return nil, errors.New("invalid data length")
	}
	return ary, nil
}

func unpackByteSlice2(buf *bytes.Buffer) ([][]byte, error) {
	n, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	ary := make([][]byte, n)
	for i := range ary {
		ary[i] = make([]byte, 32)
		if _, err := buf.Read(ary[i]); err != nil {
			return nil, err
		}
	}
	return ary, nil
}
func unpackByte32(buf *bytes.Buffer) []byte {
	ary := make([]byte, 32)
	if _, err := buf.Read(ary); err != nil {
		panic(err)
	}
	return ary
}
func unpackByteSlice(buf *bytes.Buffer) ([]byte, error) {
	n, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	ary := make([]byte, n)
	if _, err := buf.Read(ary); err != nil {
		return nil, err
	}
	return ary, nil
}

//UnpackBody return a transaction bodyfrom byte slice.
func UnpackBody(dat []byte) (*Body, error) {
	var err error
	if len(dat) < 36 {
		return nil, errors.New("dat is too short")
	}
	bd := &Body{}
	bd.Type = make([]byte, 4)
	copy(bd.Type, dat)
	bd.Nonce = make([]byte, 32)
	copy(bd.Nonce, dat[4:])
	buf := bytes.NewBuffer(dat[4+32:])
	tim, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	bd.Time = uint32(tim)
	bd.Message, err = unpackByteSlice(buf)
	if err != nil {
		return nil, err
	}

	n, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	bd.Inputs = make([]*Input, n)
	for i := range bd.Inputs {
		in := &Input{
			PreviousTX: unpackByte32(buf),
		}
		in.Index, err = buf.ReadByte()
		if err != nil {
			return nil, err
		}
		bd.Inputs[i] = in
	}

	n, err = byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	bd.Outputs = make([]*Output, n)
	for i := range bd.Outputs {
		out := &Output{
			Address: unpackByte32(buf),
		}
		out.Value, err = byte2Varint(buf)
		if err != nil {
			return nil, err
		}
		bd.Outputs[i] = out
	}

	n, err = byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	bd.MultiSigs = make([]*MultiSig, n)
	for i := range bd.MultiSigs {
		mul := &MultiSig{}
		mul.N, err = buf.ReadByte()
		if err != nil {
			return nil, err
		}
		mul.Addresses, err = unpackByteSlice2(buf)
		if err != nil {
			return nil, err
		}
		mul.Value, err = byte2Varint(buf)
		if err != nil {
			return nil, err
		}
		bd.MultiSigs[i] = mul
	}
	bd.Previous, err = unpackByteSlice2(buf)
	if err != nil {
		return nil, err
	}
	bd.Difficulty, err = buf.ReadByte()
	if err != nil {
		return nil, err
	}
	lt, err := byte2Varint(buf)
	if err != nil {
		return nil, err
	}
	bd.LockTime = uint32(lt)
	if buf.Len() != 0 {
		return nil, errors.New("invalid data length")
	}
	return bd, nil
}
