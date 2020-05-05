### Documentation

        conn,err:=gozk.Connect(zkList)
        if err ==nil {
            gozk.Register(conn, zkPath, serviceName, servicePort, serviceType)
            gozk.Discover(conn, zkPath)
        }else {
            log.Printf("zk connect error")
        }
        
