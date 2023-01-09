package im

import (
	"bufio"
	"encoding/binary"
	"strings"
)

type OpCode int

const (
	opCodeUndefined = OpCode(iota)
	opCodeMin

	OpCodeSubscribeChannel
	OpCodeUnsubscribeChannel
	OpCodeSendMessage
	OpCodeCancelMessage
	OpCodeReadMessage
	OpCodeFetchHistoryMessages

	opCodeMax
)

type Conn struct {
	r bufio.Reader
	w bufio.Writer
}

func (c Conn) ReadOpcode() OpCode {
	num := c.readUint8()
	if num <= uint8(opCodeMin) || num >= uint8(opCodeMax) {
		return opCodeUndefined
	}
	return OpCode(num)
}

func (c Conn) ReadMessageType() MessageType {
	num := c.readUint8()
	if num <= uint8(messageTypeMin) || num >= uint8(messageTypeMax) {
		return messageTypeUndefined
	}
	return MessageType(num)
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

func (c Conn) readUint64() uint64 {
	var buf [8]byte
	for i := 0; i < 8; i++ {
		b, e := c.r.ReadByte()
		if e != nil {
			panic(e)
		}
		buf[i] = b
	}
	return binary.BigEndian.Uint64(buf[:])
}

func (c Conn) readTxt() string {
	remains := int(c.readUint16())
	if remains < 1 {
		return ""
	}

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

func (c Conn) readBytes(dist *[]byte) {
	remains := int(c.readUint16())
	if remains < 0 {
		*dist = make([]byte, 0, 0)
		return
	}

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

		*dist = append(*dist, bv[:l]...)
		remains -= l
	}
}

func (c Conn) ReadChannelId() string {
	return c.readTxt()
}

func (c Conn) readMap() map[string]string {
	kvCount := int(c.readUint16())
	if kvCount < 1 {
		return nil
	}
	m := make(map[string]string, kvCount)
	for i := 0; i < int(c.readUint16()); i++ {
		m[c.readTxt()] = c.readTxt()
	}
	return m
}

func (c Conn) ReadMsg() *Message {
	var msg Message
	stop := false
	for !stop {
		switch c.readUint8() {
		case 0:
			{
				msg.From = c.readUint64()
				continue
			}
		case 1:
			{
				msg.Unix = c.readUint64()
				continue
			}
		case 2:
			{
				msg.Until = c.readUint64()
				continue
			}
		case 3:
			{
				msg.Type = c.ReadMessageType()
				if msg.Type == messageTypeUndefined {
					return nil
				}
				continue
			}
		case 4:
			{
				c.readBytes(&msg.Content)
				continue
			}
		case 5:
			{
				msg.Ext = c.readMap()
				continue
			}
		default:
			{
				stop = true
				continue
			}
		}
	}
	if msg.From == 0 || msg.Unix == 0 || msg.Type == messageTypeMin {
		return nil
	}
	return &msg
}

func (c Conn) WriteMsg() {

}
