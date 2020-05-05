package main

import (
	"fmt"
	"gozk/gozk"
	"log"
	"net/http"
)

func main() {
	zkList 		:= []string{"127.0.0.1:2181"}
	zkPath 		:= "services"
	serviceName	:= "otrade-go"
	servicePort	:= 8897
	serviceType	:= "DYNAMIC"
	conn,err:=gozk.Connect(zkList)
	if err ==nil {
		gozk.Register(conn, zkPath, serviceName, servicePort, serviceType)
		gozk.Discover(conn, zkPath)
	}else {
		log.Printf("zk connect error")
	}

	http.HandleFunc("/api/serverlist", ServerListAPIHandler)
	http.HandleFunc("/api/getserver", SeverInfoAPIHandler)
	http.ListenAndServe("0.0.0.0:8897", nil)
}

func ServerListAPIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, gozk.Serverlist())
}
func SeverInfoAPIHandler(w http.ResponseWriter, r *http.Request) {

	serviceName, ok := r.URL.Query()["name"]
	if ok{
		serverInfo:=gozk.GetServerInfo(serviceName[0])
		fmt.Fprintln(w, serverInfo)

	}else {
		fmt.Fprintln(w, "name empty")
	}


}