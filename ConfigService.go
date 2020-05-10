package gozk

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
)

func ConfigData(conn *zk.Conn,zkpath string)(string) {

	data,  _, err := conn.Get(zkpath)
	if err != nil {
		log.Println(err)
	}
	return string(data)
}

func ConfigWatch(conn *zk.Conn,zkpath string,callback func(interface{})) {
	for {

		data, _, get_ch, _ := conn.GetW(zkpath)
		callback(string(data))
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