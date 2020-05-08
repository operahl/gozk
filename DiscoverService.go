package gozk

import (
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"math/rand"
	"strings"
	"sync"
)
type DiscoverService struct {
	ServerList   map[string][]ServiceInfo
	conn *zk.Conn
	sync.RWMutex
}

var discoverService DiscoverService

func (this *DiscoverService)ZkChildrenWatch(path string)  {
	for {
		v, _, get_ch, err := this.conn.ChildrenW(path)
		if err != nil {
			log.Println(err)
		}

		this.UpdateServices(path,v)
		for {
			select {
			case ch_event := <-get_ch:
				{
					if ch_event.Type == zk.EventNodeCreated {
						log.Printf("has new node[%d] create\n", ch_event.Path)
					} else if ch_event.Type == zk.EventNodeDeleted {
						log.Printf("has node[%s] detete\n", ch_event.Path)
					} else if ch_event.Type == zk.EventNodeDataChanged {
						log.Printf("node change%+v\n", ch_event.Path)
					} else if ch_event.Type == zk.EventNodeChildrenChanged {
						log.Printf("children change%+v\n", ch_event.Path)
					}

				}
			}
			break
		}
	}
}

func (this *DiscoverService)UpdateServices(path string,v []string) {

	this.Lock()
	defer this.Unlock()
	pathArr :=strings.Split(path,"/")
	serviceName := "";
	if len(pathArr)>0{
		serviceName = pathArr[len(pathArr)-1]
	}else{
		return
	}
	var services []ServiceInfo
	for _,value :=range v{
		servicePath := path+"/"+value
		log.Printf("value of path[%s]=[%s].\n", servicePath, value)
		data, _, err := this.conn.Get(servicePath)
		if err != nil {
			log.Print(err)
		}
		ser:=ServiceInfo{}

		json.Unmarshal(data,&ser)
		log.Print(ser)
		services=append(services, ser)
	}
	if services!=nil{
		this.ServerList[serviceName] = services
	}
}


func Discover(conn *zk.Conn,zkpath string)  {
	go discoverService.Subscribe(conn,zkpath)
}

func Serverlist()(map[string][]ServiceInfo)  {

	discoverService.RLock()
	defer discoverService.RUnlock()
	return discoverService.ServerList
}

func GetServerInfo(serviceName string)(ServiceInfo){
	discoverService.RLock()
	defer discoverService.RUnlock()
	if _,ok:=discoverService.ServerList[serviceName];ok{

		serviceCount :=len(discoverService.ServerList[serviceName])
		if(serviceCount == 1){
			return discoverService.ServerList[serviceName][0]
		}else{
			key :=rand.Intn(serviceCount)
			return discoverService.ServerList[serviceName][key]
		}
	}

	return ServiceInfo{}
}

func (w *DiscoverService)Subscribe(conn *zk.Conn,zkpath string) {

	w.conn = conn
	children, _, _ := w.conn.Children("/"+zkpath)
	var wg sync.WaitGroup
	wg.Add(len(children))
	serviceInfo :=ServiceInfo{}
	serviceInfoMap :=make(map[string][]ServiceInfo,0)
	for _, path := range children {
		serviceInfoMap[path] = append(serviceInfoMap[path], serviceInfo)
	}
	w.ServerList = make(map[string][]ServiceInfo,0)

	for _, path := range children {
		path = "/"+zkpath + "/" + path
		v, _, _, err := w.conn.ChildrenW(path)
		if err != nil {
			log.Println(err)
		}
		w.UpdateServices(path,v)
	}
	log.Print(w.ServerList)
 	for _, path := range children {
		path = "/"+zkpath + "/" + path
		go func(path string) {
			defer wg.Done()
			log.Print("Zookeeper Watcher Starting, ", path)
			w.ZkChildrenWatch(path)
		}(path)
	}
	wg.Wait()
}