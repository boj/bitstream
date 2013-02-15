// Copyright 2013 Brian 'bojo' Jones. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package bitstream is a simple library for sending byte data across the wire.

Bitstream messages are built using the following guidelines for primitive types:

    [total message length]  - 4 bytes
      [byte][value]         - 1 byte + 1 byte
      [int32][value]        - 1 byte + 4 bytes
      [float32][value]      - 1 byte + 4 bytes
      [bool (byte)][value]  - 1 byte + 1 byte
      [string][size][value] - 1 byte + 4 bytes + N bytes
      ...

Message reading and writing is completely user driven.

The user defined network loop handles reading the first 4 bytes to determine the message size.
The following data should then be passed to the instantiated Read() method.

Example:

    // Sender
    bs := NewBitStream()
    bs.WriteByte(messageType)
    bs.WriteInt(objectId)
    bs.WriteFloat(position.x)
    bs.WriteFloat(position.y)
    bs.WriteFloat(position.z)
    msg := bs.Write()

    // Network loop reads first 4 bytes, passes rest of data to MessageHandler()

    // Receiver
    func MessageHandler(msg []byte) error {
      bs := NewBitStream()
      if err := bs.Read(msg); err != nil {
        return err
      } else {  
        if msgType, err := bs.ReadByte(); err != nil {
          return err
        } else {
          switch msgType {
            case ENTITY_MOVED:
              // read id
              id, err := bs.ReadInt()
              if err != nil {
                return err
              }
              e := entity[id]
              // read position
              x, err := bs.ReadFloat()
              if err != nil {
                return err
              }
              y, err := bs.ReadFloat()
              if err != nil {
                return err
              }
              z, err := bs.ReadFloat()
              if err != nil {
                return err
              }
              e.SetPosition(x, y, z)
              break
          }
        }
      }
    }

This library assumes the reciever understands the order the message parts are received in.
*/
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

// NewBitStream creates a new bitstream.
func NewBitStream() *BitStream {
	self := new(BitStream)
	return self
}

// Length returns the length of the byte data contained in the bitstream.
// This does not include the first 4 bytes read by the network loop.
func (self *BitStream) Length() int {
	return self.data.Len()
}

// Reset resets the state of the bitstream to nothing.
func (self *BitStream) Reset() {
	self.data.Reset()
}

// Read takes in a stream of bytes read in from the network, not including
// the initial 4 bytes read by the network loop.  The library can then
// parse and return the results.
func (self *BitStream) Read(b []byte) error {
	if len(b) == 0 {
		return errors.New("BitStream: passed 0 length message to Read()\n")
	}
	self.data.Write(b)
	return nil
}

// ReadByte returns a byte.
func (self *BitStream) ReadByte() (byte, error) {
	var check []byte = self.data.Next(1)
	if check[0] != TYPE_BYTE {
		return 0, errors.New("BitStream: mismatching data type. expected TYPE_BYTE\n")
	}
	buffer := self.data.Next(1)
	return buffer[0], nil
}

// ReadInt returns an int32.
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

// ReadFloat returns a float32.
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

// ReadBool returns a bool.
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

// ReadString returns a string.
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

// WriteByte writes a byte to the bitstream.
func (self *BitStream) WriteByte(b byte) {
	self.data.WriteByte(TYPE_BYTE)
	self.data.WriteByte(b)
}

// WriteInt writes a in32 to the bitstream.
func (self *BitStream) WriteInt(i int32) {
	self.data.WriteByte(TYPE_INT)
	tmp := new(bytes.Buffer)
	binary.Write(tmp, binary.LittleEndian, i)
	self.data.Write(tmp.Bytes())
}

// WriteFloat writes a float32 to the bitstream.
func (self *BitStream) WriteFloat(f float32) {
	self.data.WriteByte(TYPE_FLOAT)
	tmp := new(bytes.Buffer)
	binary.Write(tmp, binary.LittleEndian, f)
	self.data.Write(tmp.Bytes())
}

// WriteBool writes a bool to the bitstream.
func (self *BitStream) WriteBool(b bool) {
	self.data.WriteByte(TYPE_BOOL)
	if b {
		self.data.WriteByte(1)
	} else {
		self.data.WriteByte(0)
	}
}

// WriteString writes a string to the bitstream.
func (self *BitStream) WriteString(s string) {
	self.data.WriteByte(TYPE_STRING)
	size := new(bytes.Buffer)
	binary.Write(size, binary.LittleEndian, int32(len(s)))
	self.data.Write(size.Bytes())
	self.data.WriteString(s)
}

// Write returns the bitstream data as a byte array.
func (self *BitStream) Write() []byte {
	// write message size
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, int32(len(self.data.Bytes())))
	// write data
	b.Write(self.data.Bytes())
	return b.Bytes()
}
