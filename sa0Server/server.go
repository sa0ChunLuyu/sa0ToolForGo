package sa0Server

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	Name   string
	Path   string
	Method []string
	Origin []string
	Func   func(string, func(string), string)
}

var sa0Router []Router
var Sa0Path string
var Sa0Index string

func Config_(path_ string) {
	Sa0Path = path_
}

func Server(router_ []Router, path_ string, index_ string) {
	sa0Router = router_
	_, port := GetConfig_("server", "port")
	if port == "" {
		port = "2333"
	}
	Sa0Index = index_
	http.HandleFunc(path_, build_)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(Sa0Path+"/assets"))))
	fmt.Println("Server on http://127.0.0.1:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func RouterInfo(key string) [][]string {
	_, method := GetConfig_(key+"_router_info", "method")
	_, origin := GetConfig_(key+"_router_info", "origin")
	return [][]string{
		strings.Split(method, ","),
		strings.Split(origin, ","),
	}
}

func GetConfig_(group_ string, key_ string) (bool, string) {
	cfg, err := goconfig.LoadConfigFile(Sa0Path + "/config.ini")
	if err != nil {
		return false, ""
	}
	value_, err := cfg.GetValue(group_, key_)
	return true, value_
}

var Sa0R *http.Request
var Sa0W http.ResponseWriter

func build_(w http.ResponseWriter, r *http.Request) {
	Sa0R = r
	Sa0W = w
	_, routeInfo := Data_("routeInfo");
	routeInfoArray := strings.Split(routeInfo, "/")
	controllerName := routeInfoArray[0]
	controllerState := false
	for routerIndex, routerItem := range sa0Router {
		if routerItem.Path == controllerName {
			controllerState = true
			methodState := false
			for _, methodItem := range sa0Router[routerIndex].Method {
				if methodItem == r.Method || methodItem == "*" {
					methodState = true
					break
				}
			}
			if methodState {
				Sa0W.Header().Set("Access-Control-Allow-Origin", "*")
				Sa0W.Header().Add("Access-Control-Allow-Headers", "*")
				Sa0W.Header().Set("content-type", "*")
				functionName := routeInfoArray[1]

				sa0Router[routerIndex].Func(functionName, Print_, Error_(map[string]string{
					"message": "FUNCTION NOT FOUND",
				}))
			} else {
				Print_(Error_(map[string]string{
					"message": "REQUEST METHOD NOT ALLOW",
				}))
			}
			break
		}
	}
	if controllerState == false {
		if Sa0Index != "NOINDEX" {
			Print_(Jump_(map[string]string{
				"indexPath": Sa0Index,
			}))
		} else {
			Print_(Error_(map[string]string{
				"message": "CONTROLLER NOT FOUND",
			}))
		}
	}
}

func Data_(key string) (bool, string) {
	keys, ok := Sa0R.URL.Query()[key]
	if !ok || len(keys) < 1 {
		return false, ""
	}
	return true, string(keys[0])
}

func Print_(content string) {
	_, _ = fmt.Fprintf(Sa0W, content+"\n")
}

func Jump_(data map[string]string) string {
	view := Jump
	for key, value := range data {
		view = strings.Replace(view, "<<."+key+">>", value, -1)
	}
	return view
}

func Error_(data map[string]string) string {
	view := Error
	for key, value := range data {
		view = strings.Replace(view, "<<."+key+">>", value, -1)
	}
	return view
}
