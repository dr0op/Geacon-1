package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	TRUE = true
	Buffer bytes.Buffer
)

func SplitX(Str string, Sep ...rune) []string {
	return strings.FieldsFunc(Str, func(r rune) bool {
		for _, v := range Sep {
			if r == v { return true }
		}
		return false
	})
}

func IsLocal(Addr string) bool {
	return strings.Contains(Addr, "0.0.0.0") || strings.Contains(Addr, "127.0.0.1")
}

func RandInt(Min, Max int) int {
	rand.Seed(time.Now().UnixNano())
	if Max > Min {
		return Min + rand.Intn(Max-Min)
	}
	if Max < Min {
		return Max + rand.Intn(Min-Max)
	}
	return Max
}

func RandByte(Min, Max int) []byte {
	Len := RandInt(Min, Max)
	Buf := make([]byte, Len)
	rand.Read(Buf)
	return Buf
}

func ReadFrom(r io.Reader, Len int) []byte {
	if Len <= 0 { return nil }
	Buf := make([]byte, Len)
	_, err := io.ReadFull(r, Buf)
	if err != nil { return nil }
	return Buf
}

func ToLength(Str string) int {
	v, err := strconv.Atoi(Str)
	if err == nil { return v }
	return len(Str)
}

func JoinBytes(Bytes ...[]byte) []byte {
	return bytes.Join(Bytes, []byte(""))
}

func IntAsByte(Len, Num int) []byte {
	Buf := make([]byte, Len)
	switch Len {
	case 1:
		Buf[0] = byte(Num)
	case 2:
		binary.LittleEndian.PutUint16(Buf, uint16(Num))
	case 4:
		binary.LittleEndian.PutUint32(Buf, uint32(Num))
	}
	return Buf
}

func IntToByte(Len, Num int) []byte {
	Buf := make([]byte, Len)
	switch Len {
	case 1:
		Buf[0] = byte(Num)
	case 2:
		binary.BigEndian.PutUint16(Buf, uint16(Num))
	case 4:
		binary.BigEndian.PutUint32(Buf, uint32(Num))
	}
	return Buf
}

func ByteAsInt(Bytes []byte) int {
	switch len(Bytes) {
	case 1:
		return int(Bytes[0])
	case 2:
		return int(binary.LittleEndian.Uint16(Bytes))
	case 4:
		return int(binary.LittleEndian.Uint32(Bytes))
	}
	return 0
}

func ByteToInt(Bytes []byte) int {
	switch len(Bytes) {
	case 1:
		return int(Bytes[0])
	case 2:
		return int(binary.BigEndian.Uint16(Bytes))
	case 4:
		return int(binary.BigEndian.Uint32(Bytes))
	}
	return 0
}

func BigToLittle(Bytes []byte) []byte {
	return IntAsByte(len(Bytes), ByteToInt(Bytes))
}

func LittleToBig(Bytes []byte) []byte {
	return IntToByte(len(Bytes), ByteAsInt(Bytes))
}