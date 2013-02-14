package bitstream

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	TYPE_BYTE   byte = 1 // byte
	TYPE_INT    byte = 2 // int32
	TYPE_FLOAT  byte = 3 // float32
	TYPE_BOOL   byte = 4 // byte
	TYPE_STRING byte = 5 // len(string) + string
)

type BitStream struct {
	data bytes.Buffer
}

func NewBitStream() *BitStream {
	self := new(BitStream)
	return self
}

func (self *BitStream) Length() int {
	return self.data.Len()
}

func (self *BitStream) Reset() {
	self.data.Reset()
}

func (self *BitStream) Read(b []byte) error {
	if len(b) == 0 {
		return errors.New("BitStream: passed 0 length message to Read()\n")
	}
	self.data.Write(b)
	return nil
}

func (self *BitStream) ReadByte() (byte, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_BYTE {
		return 0, errors.New("BitStream: mismatching data type. expected TYPE_BYTE\n")
	}
	buffer := self.data.Next(1)
	return buffer[0], nil
}

func (self *BitStream) ReadInt() (int32, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_INT {
		return 0, errors.New("BitStream: mismatching data type. expected TYPE_INT\n")
	}
	var ret int32
	buf := bytes.NewBuffer(self.data.Next(4))
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret, nil
}

func (self *BitStream) ReadFloat() (float32, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_FLOAT {
		return 0, errors.New("BitStream: mismatching data type. expected TYPE_FLOAT\n")
	}
	var ret float32
	buf := bytes.NewBuffer(self.data.Next(4))
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret, nil

}

func (self *BitStream) ReadBool() (bool, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_BOOL {
		return false, errors.New("BitStream: mismatching data type. expected TYPE_BOOL\n")
	}
	buffer := self.data.Next(1)
	if buffer[0] == byte(1) {
		return true, nil
	}
	return false, nil
}

func (self *BitStream) ReadString() (string, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_STRING {
		return "", errors.New("BitStream: mismatching data type. expected TYPE_STRING\n")
	}
	var size int32
	buf := bytes.NewBuffer(self.data.Next(4))
	binary.Read(buf, binary.LittleEndian, &size)
	ret := self.data.Next(int(size))
	return string(ret), nil
}

func (self *BitStream) WriteByte(b byte) {
	self.data.WriteByte(TYPE_BYTE)
	self.data.WriteByte(b)
}

func (self *BitStream) WriteInt(i int32) {
	self.data.WriteByte(TYPE_INT)
	tmp := new(bytes.Buffer)
	binary.Write(tmp, binary.LittleEndian, i)
	self.data.Write(tmp.Bytes())
}

func (self *BitStream) WriteFloat(f float32) {
	self.data.WriteByte(TYPE_FLOAT)
	tmp := new(bytes.Buffer)
	binary.Write(tmp, binary.LittleEndian, f)
	self.data.Write(tmp.Bytes())
}

func (self *BitStream) WriteBool(b bool) {
	self.data.WriteByte(TYPE_BOOL)
	if b {
		self.data.WriteByte(1)
	} else {
		self.data.WriteByte(0)
	}
}

func (self *BitStream) WriteString(s string) {
	self.data.WriteByte(TYPE_STRING)
	size := new(bytes.Buffer)
	binary.Write(size, binary.LittleEndian, int32(len(s)))
	self.data.Write(size.Bytes())
	self.data.WriteString(s)
}

func (self *BitStream) Write() []byte {
	// write message size
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, int32(len(self.data.Bytes())))
	// write data
	b.Write(self.data.Bytes())
	return b.Bytes()
}
