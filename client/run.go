package client

import (
	"fmt"
	"io"
	"log"
	"net"
	"novpn/com"
	"novpn/config"
	"sync"
)

var (
	// 会话id
	opencmd = make([]byte, 3)
	// 全局连接计数器
	count   uint16 = 0
	countmu sync.Mutex
	key, iv []byte
)

func initSession(cfg *config.ClientConfig) {
	conn, err := net.Dial("tcp", cfg.Exchange)
	if err != nil {
		log.Println("Open exchange error", err)
		return
	}
	defer conn.Close()
	// 认证
	buf := make([]byte, 33)
	buf[0] = com.CLIENTC
	copy(buf[1:], com.GetMd5(cfg.Key))
	copy(buf[17:], com.GetMd5(cfg.ID))
	_, err = conn.Write(buf)
	if err != nil {
		log.Println("Init service is error", err)
		return
	}
	opencmd[0] = com.CLIENTD
	_, err = io.ReadAtLeast(conn, opencmd[1:3], 2)
	if err != nil {
		log.Println("Error opening client.", err, opencmd)
		return
	}
}

// 处理新连接
func handleConn(conn net.Conn, cfg *config.ClientConfig) {
	econn, err := net.Dial("tcp", cfg.Exchange)
	if err != nil {
		log.Println("Open exchange error", err)
		conn.Close()
		return
	}
	_, err = econn.Write(opencmd)
	if err != nil {
		log.Println("Open exchange error", err)
		conn.Close()
		econn.Close()
		return
	}
	var s com.NCopy
	s.Init(econn, key, iv)
	go com.RCopy(conn, &s, "Client conn", &count, &countmu)
	go com.WCopy(&s, conn, "Client conn", &count, &countmu)
	// go com.NetCopy(conn, econn, "Client conn", &count, &countmu)
	// go com.NetCopy(econn, conn, "Client conn", &count, &countmu)
}

// Run 运行
func Run(cfg *config.ClientConfig) {
	if cfg == nil {
		return
	}
	// 初始化加密
	key, iv = com.GetKeyIv(cfg.Password)
	// 初始化sessionid
	initSession(cfg)
	// 监听本地连接
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", cfg.Port))
	if err != nil {
		log.Println("Client listen error", err)
		return
	}
	defer lis.Close()
	log.Println("Client is running")
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Local error", err)
			return
		}
		go handleConn(conn, cfg)
	}
}
