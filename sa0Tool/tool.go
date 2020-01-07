package sa0Tool

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/sa0ChunLuyu/sa0ToolForGo/sa0Server"
)

func Data_(key string) (bool, string) {
	keys, ok := sa0Server.Sa0R.URL.Query()[key]
	if !ok || len(keys) < 1 {
		return false, ""
	}
	return true, string(keys[0])
}

func Print_(content string) {
	_, _ = fmt.Fprintf(sa0Server.Sa0W, content+"\n")
}

func GetConfig_(group_ string, key_ string) (bool, string) {
	cfg, err := goconfig.LoadConfigFile(sa0Server.Sa0Path + "/config.ini")
	if err != nil {
		return false, ""
	}
	value_, err := cfg.GetValue(group_, key_)
	return true, value_
}
