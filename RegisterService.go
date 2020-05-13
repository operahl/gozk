package gozk

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/satori/go.uuid"
	"log"
	"net"
	"time"
)
type ServiceInfo struct {
	Id		string `json:"id"`
	Name	string  `json:"name"`
	Address	string  `json:"address"`
	Port	int  `json:"port"`
	ServiceType	string  `json:"serviceType"`
}

type RegisterService struct {
	root    string
	conn      *zk.Conn
}

var registerService RegisterService;

func Connect(zkServers []string) (*zk.Conn,error) {
	conn, _, err := zk.Connect(zkServers, 10*time.Second)
	return conn,err
}

func Register(conn *zk.Conn,zkNode string,serviceName string,servicePort int,ServiceType string) {
	go registerService.RegisterZK(conn, zkNode, serviceName, servicePort, ServiceType)
}

func (s *RegisterService)checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (s *RegisterService)RegisterZK(conn *zk.Conn,zkNode string,serviceName string,servicePort int,ServiceType string) {
	var err error
	s.conn = conn

	exist, _, _, err := conn.ExistsW("/"+zkNode)
	if !exist {
		if _, err := conn.Create("/"+zkNode, []byte("init"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			fmt.Println("Create returned error: %v", err)
		}

	}
	serviceNode := "/"+zkNode+"/"+serviceName
	exist, _, _, err = conn.ExistsW(serviceNode)
	if !exist {
		if _, err := conn.Create(serviceNode, []byte("init"), 0, zk.WorldACL(zk.PermAll)); err != nil {
			fmt.Println("Create returned error: %v", err)
		}
	}
	u1 := uuid.Must(uuid.NewV4(),err)

	var serviceInfo	ServiceInfo

	serviceInfo.Id=u1.String()
	serviceInfo.Name=serviceName
	serviceInfo.Port = servicePort
	serviceInfo.Address=s.getLocalServerIp()
	serviceInfo.ServiceType = ServiceType

	s.root =  fmt.Sprintf("/%s/%s/%s",zkNode,serviceName,u1.String())
	jsonStr,_:=json.Marshal(serviceInfo)
	err = s.RegistServer(string(jsonStr))
	if err != nil {
		fmt.Printf(" regist node error: %s ", err)
	}
}


func (s *RegisterService)RegistServer( str string) (err error) {
	if err := s.conn.Delete(s.root, -1); err != nil {
		//fmt.Println("Delete returned error: %+v", err)
	}
	if _, err := s.conn.Create(s.root, []byte(str), zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
		log.Printf("Create returned error: %v", err)
	}

	return err
}

func (s *RegisterService)GetServerList() (list []string, err error) {
	list, _, err = s.conn.Children(s.root)
	return
}

func (s *RegisterService)ensureRoot() error {
	exists, _, _ := s.conn.Exists(s.root)
	if !exists {
		_, err := s.conn.Create(s.root, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}
func (s *RegisterService)getLocalServerIp() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}
	ip :=""
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
					}
				}
			}
		}
	}
	return ip
}

func (s *RegisterService)mirror(path string) (chan []string, chan error) {
	snapshots := make(chan []string)
	errors := make(chan error)

	go func() {
		for {
			snapshot, _, events, err := s.conn.ChildrenW(path)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
				errors <- evt.Err
				return
			}
		}
	}()

	return snapshots, errors

}

func (s *RegisterService)Watch(node string){
	snapshots, errors := s.mirror( node)
	go func() {
		for {
			select {
			case snapshot := <-snapshots:
				fmt.Printf("%+v\n", snapshot)
			case err := <-errors:
				fmt.Printf("%+v\n", err)
			}
		}
	}()
}

