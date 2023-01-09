package im

import (
	"bufio"
	"encoding/binary"
	"strings"
)

type OpCode int

const (
	OpCodeUndefined = OpCode(iota)
	OpCodeMin

	OpCodeSubscribeChannel
	OpCodeUnsubscribeChannel
	OpCodeSendMessage
	OpCodeCancelMessage
	OpCodeReadMessage
	OpCodeFetchHistoryMessages

	OpCodeMax
)

type Conn struct {
	r bufio.Reader
	w bufio.Writer
}

func (c Conn) readOpcode() OpCode {
	num := c.readUint8()
	if num <= uint8(OpCodeMin) || num >= uint8(OpCodeMax) {
		return OpCodeUndefined
	}
	return OpCode(num)
}

func (c Conn) readUint8() uint8 {
	b, e := c.r.ReadByte()
	if e != nil {
		panic(e)
	}
	return b
}

func (c Conn) readUint16() uint16 {
	var buf [2]byte
	for i := 0; i < 2; i++ {
		b, e := c.r.ReadByte()
		if e != nil {
			panic(e)
		}
		buf[i] = b
	}
	return binary.BigEndian.Uint16(buf[:])
}

func (c Conn) readTxt() string {
	remains := int(c.readUint16())

	sb := strings.Builder{}
	sb.Grow(remains)

	var buf [256]byte

	for remains > 0 {
		var bv []byte
		if remains < 256 {
			bv = buf[:remains]
		} else {
			bv = buf[:]
		}

		l, e := c.r.Read(bv)
		if e != nil || l < 1 {
			panic(e)
		}

		sb.Write(bv[:l])
		remains -= l
	}
	return sb.String()
}

func (c Conn) readChannelId() string {
	return c.readTxt()
}
