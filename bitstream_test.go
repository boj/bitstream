package bitstream

import (
	"testing"
)

func TestBitStream(t *testing.T) {
	var test_byte byte = 32
	var test_int int32 = 95959
	var test_float float32 = 0.0001
	var test_bool bool = true
	var test_string string = "PETER PIPER PICKED A PECK of pickled peppers!"

	bs := NewBitStream()
	bs.WriteByte(test_byte)
	bs.WriteInt(test_int)
	bs.WriteFloat(test_float)
	bs.WriteBool(test_bool)
	bs.WriteString(test_string)
	msg := bs.Write()
	t.Logf("Length: %d - Message: %v\n", bs.Length(), msg)

	bs.Reset()

	err := bs.Read(msg[4:]) // network reader reads first 4 bytes to determine message size
	if err != nil {
		t.Log(err.Error())
	}
	check_byte, err := bs.ReadByte()
	if err != nil {
		t.Log(err.Error())
	}
	check_int, err := bs.ReadInt()
	if err != nil {
		t.Log(err.Error())
	}
	check_float, err := bs.ReadFloat()
	if err != nil {
		t.Log(err.Error())
	}
	check_bool, err := bs.ReadBool()
	if err != nil {
		t.Log(err.Error())
	}
	check_string, err := bs.ReadString()
	if err != nil {
		t.Log(err.Error())
	}

	if check_byte != test_byte {
		t.Errorf("Byte check failed.  Expected %b got %b", test_byte, check_byte)
	}
	if check_int != test_int {
		t.Errorf("Int check failed.  Expected %d got %d", test_int, check_int)
	}
	if check_float != test_float {
		t.Errorf("Float check failed.  Expected %f got %f", test_float, check_float)
	}
	if check_bool != test_bool {
		t.Errorf("Bool check failed.  Expected %t got %t", test_bool, check_bool)
	}
	if check_string != test_string {
		t.Errorf("String check failed.  Expected %s got %s", test_string, check_string)
	}
}

func TestFailBitStream(t *testing.T) {
	bs := NewBitStream()
	bs.WriteString("string")
	msg := bs.Write()

	bs.Reset()
	bs.Read(msg)
	// fail - supposed to read string, not int
	_, err := bs.ReadInt()
	if err == nil {
		t.Errorf("ReadInt against WriteString should throw error")
	}
}
