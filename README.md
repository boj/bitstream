# BitStream

A simple library for sending byte data across the wire.

## Usage

A message gets constructed like the following:

    [total message length]  - 4 bytes
      [byte][value]         - 1 byte + 1 byte
      [int32][value]        - 1 byte + 4 bytes
      [float32][value]      - 1 byte + 4 bytes
      [bool (int8)][value]  - 1 byte + 1 byte
      [string][size][value] - 1 byte + 4 bytes + N bytes
      ...
    
Message reading and writing is completely user driven.

The network loop handles reading the first 4 bytes to determine the message size, and passes the rest of the data to the BitStream.Read() method.

Example:

    // Sender
    bs := NewBitStream()
    bs.WriteByte(messageType)
    bs.WriteInt(objectId)
    bs.WriteFloat32(position.x)
    bs.WriteFloat32(position.y)
    bs.WriteFloat32(position.z)
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

## Author

Brian 'bojo' Jones <mojobojo@gmail.com>

## License

The MIT License
 
Copyright (c) 2012 Brian Jones
 
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
 
The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.
 
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
