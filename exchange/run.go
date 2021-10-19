package exchange

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"novpn/com"
	"novpn/config"
)

// handle 处理新连接
func handle(conn net.Conn) {
	cmd := make([]byte, 100)
	_, err := io.ReadAtLeast(conn, cmd[:1], 1)
	if err != nil {
		conn.Close()
		return
	}
	switch cmd[0] {
	case com.SERVICEC:
		// 服务端控制端
		// key id
		// 16  16
		_, err = io.ReadAtLeast(conn, cmd[1:33], 32)
		if err != nil {
			conn.Close()
			return
		}
		// 校验密码
		if !bytes.Equal(key, cmd[1:17]) {
			log.Println("The password is incorrect.")
			conn.Close()
			return
		}
		// 检查ID是否存在
		for id, svc := range svcs {
			if svc != nil && bytes.Equal(svc.ID, cmd[17:33]) {
				// 已经存在
				// 替换现有的控制端
				oldconn := svc.Conn
				svc.Conn = conn
				oldconn.Close()
				conn.Write([]byte{byte(id >> 8), byte(id & 0xff)})
				log.Println("Service relogin", id)
				return
			}
		}
		// 新连接
		id := nextID()
		if id == 0 {
			// id生成满了
			conn.Close()
			return
		}
		svc := Service{
			Conn:    conn,
			workers: make([]*Worker, 1<<16),
			ID:      cmd[17:33],
		}
		conn.Write([]byte{byte(id >> 8), byte(id & 0xff)})
		svcs[id] = &svc
		log.Println("Service login", id)
	case com.CLIENTC:
		// 新登录客户端
		// key id
		// 16  16
		_, err = io.ReadAtLeast(conn, cmd[1:33], 32)
		if err != nil {
			conn.Close()
			return
		}
		// 校验密码
		if !bytes.Equal(key, cmd[1:17]) {
			log.Println("The password is incorrect.")
			conn.Close()
			return
		}
		// 检查对应服务端是否登录
		for id, svc := range svcs {
			// 存在，检查服务端是否存在
			if svc != nil && bytes.Equal(svc.ID, cmd[17:33]) {
				// 写入会话id
				conn.Write([]byte{byte(id >> 8), byte(id & 0xff)})
			}
		}
		conn.Close()
	case com.SERVICED:
		// 服务端数据端
		// id dataid
		// 2  2
		_, err = io.ReadAtLeast(conn, cmd[1:5], 4)
		if err != nil {
			log.Println("Data read error.", err)
			conn.Close()
			return
		}
		id := (uint16(cmd[1]) << 8) | uint16(cmd[2])
		svc := svcs[id]
		if svc == nil {
			log.Println("Service does not exist.", id)
			conn.Close()
			return
		}
		dataid := (uint16(cmd[3]) << 8) | uint16(cmd[4])
		svc.NewServiceConn(dataid, conn)
	case com.CLIENTD:
		// 客户端数据端
		// id
		// 2
		_, err = io.ReadAtLeast(conn, cmd[1:3], 2)
		if err != nil {
			log.Println("Data read error.", err)
			conn.Close()
			return
		}
		id := (uint16(cmd[1]) << 8) | uint16(cmd[2])
		svc := svcs[id]
		if svc == nil {
			log.Println("Service does not exist.", id)
			conn.Close()
			return
		}
		svc.NewClientConn(conn)
	}
}

// Run 运行
func Run(cfg *config.ExchangeConfig) {
	if cfg == nil {
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cfg.Port))
	if err != nil {
		log.Println("E: Listen error", err)
		return
	}
	defer lis.Close()
	// 准入Key
	key = com.GetMd5(cfg.Key)
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("E: Accept error", err)
			return
		}
		go handle(conn)
	}

}
