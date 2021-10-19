package exchange

import (
	"log"
	"novpn/com"
	"sync"
	"time"
)

var (
	// 准入key
	key []byte
	// 服务端列表
	// svcmap = make(map[uint16]*Service, 10)
	svcs = make([]*Service, 1<<16)
	// 全局连接计数器
	count   uint16 = 0
	countmu sync.Mutex
)

func init() {
	go func() {
		for {
			time.Sleep(com.Heartbeat * time.Second)
			for id, svc := range svcs {
				if svc != nil {
					_, err := svc.Conn.Read([]byte{})
					now := time.Now().Unix()
					if err == nil {
						svc.lastLife = now
					}
					if now-svc.lastLife > com.TimeoutWite {
						// 连接超时
						svc.Conn.Close()
						svcs[id] = nil
						log.Println("Delete service", id)
					}
				}
			}
		}
	}()
}

// 获取服务端新id
func nextID() uint16 {
	for i, v := range svcs {
		if i == 0 {
			continue
		}
		if v == nil {
			return uint16(i)
		}
	}
	return 0
}
