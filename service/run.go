package service

import (
	"io"
	"log"
	"net"
	"novpn/com"
	"novpn/config"
)

func initSession(cfg *config.ServiceConfig) (net.Conn, error) {
	log.Println("Open server...")
	ctrl, err := net.Dial("tcp", cfg.Exchange)
	if err != nil {
		log.Println("Open exchange error", err)
		return ctrl, err
	}
	// 认证
	buf := make([]byte, 33)
	buf[0] = com.SERVICEC
	copy(buf[1:], com.GetMd5(cfg.Key))
	copy(buf[17:33], com.GetMd5(cfg.ID))
	_, err = ctrl.Write(buf[:33])
	if err != nil {
		log.Println("Init service is error", err)
		return ctrl, err
	}
	_, err = io.ReadAtLeast(ctrl, OpenD[1:3], 2)
	if err != nil {
		log.Println("Init service is error", err)
		return ctrl, err
	}
	OpenD[0] = com.SERVICED
	log.Println("Server is running")
	return ctrl, nil
}

func openNewConn(id []byte, cfg *config.ServiceConfig) {
	sconn, err := net.Dial("tcp", cfg.Exchange)
	if err != nil {
		log.Println("Open exchange error", err)
		return
	}
	_, err = sconn.Write(append(OpenD, id...))
	if err != nil {
		log.Println("Open exchange error", err)
		sconn.Close()
		return
	}
	// 开始解析socks5协议
	if err = socks5(sconn); err != nil {
		sconn.Close()
	}
}

// Run 运行
func Run(cfg *config.ServiceConfig) {
	if cfg == nil {
		return
	}
	key, iv = com.GetKeyIv(cfg.Password)
	ctrl, err := initSession(cfg)
	if err != nil {
		return
	}
	buf := make([]byte, 1)
	for {
		_, err := ctrl.Read(buf)
		if err != nil {
			log.Println("Read control is break.")
			ctrl, err = initSession(cfg)
			if err != nil {
				return
			}
			continue
		}
		switch buf[0] {
		case com.OPEN:
			// 打开新连接
			id := make([]byte, 2)
			_, err = io.ReadAtLeast(ctrl, id, 2)
			if id[0] == 0 && id[1] == 0 {
				continue
			}
			if err != nil {
				log.Println("Read error", err)
				return
			}
			go openNewConn(id, cfg)
		}
	}
}
