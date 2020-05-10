### Example

        zkList      := []string{"127.0.0.1:2181"}
        zkPath      := "services"
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
        
        //config center
        configs:=gozk.ConfigData(conn,serviceConfigPath)
        //config watch
        gozk.ConfigWatch(conn,serviceConfigPath,callback)

     
### Register Server

    gozk.Register(conn, zkPath, serviceName, servicePort, serviceType)
    
### Discover Server

    gozk.Discover(conn, zkPath)
    
### Get Server random IP Port 

    gozk.GetServerInfo(serviceName)

### Config Center

    //config center
    configs:=gozk.ConfigData(conn,serviceConfigPath)
    //config watch
    gozk.ConfigWatch(conn,serviceConfigPath,callback)