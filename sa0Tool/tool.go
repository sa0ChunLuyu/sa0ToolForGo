package sa0Tool

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/sa0ChunLuyu/sa0ToolForGo/sa0Server"
	"io/ioutil"
	"net/http"
	"strings"
)

func Data_(key string) (bool, string) {
	keys, ok := sa0Server.Sa0R.URL.Query()[key]
	if !ok || len(keys) < 1 {
		return false, ""
	}
	return true, string(keys[0])
}

func Get_(url string) (bool, string) {
	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, ""
	}
	return true, string(body)
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

func View_(path string, data map[string]string) string {
	s, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	view := string(s)
	for key, value := range data {
		view = strings.Replace(view, "<<."+key+">>", value, -1)
	}
	return view
}
