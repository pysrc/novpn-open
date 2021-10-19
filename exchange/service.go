package exchange

import (
	"log"
	"net"
	"novpn/com"
	"sync"
	"time"
)

// Worker 工作对
type Worker struct {
	Sconn net.Conn // 服务端连接
	Cconn net.Conn // 客户端连接
	Ctime int64    // 客户端连接时时间戳
}

// Service 服务端
type Service struct {
	Conn     net.Conn   // 服务端控制连接
	ID       []byte     // 服务-客户端ID
	workers  []*Worker  // 工作负载
	workermu sync.Mutex // 工作负载锁
	connid   uint16     // 服务端数据连接ID当前值
	lastLife int64      // 最后活跃时间,Unix
}

// NewClientConn 客户端新进来一个连接
func (s *Service) NewClientConn(conn net.Conn) uint16 {
	// 获取新连接id
	s.workermu.Lock()
	for {
		s.connid++
		if s.connid == 0 {
			s.connid = 1
		}
		wk := s.workers[s.connid]
		if wk == nil {
			break
		}
		// 等待连接超时，移除空余位置
		if time.Now().Unix()-wk.Ctime > com.TimeoutWite {
			if wk.Cconn != nil {
				wk.Cconn.Close()
			}
			if wk.Sconn != nil {
				wk.Sconn.Close()
			}
			s.workers[s.connid] = nil
			break
		}
	}
	cid := s.connid
	s.workers[cid] = &Worker{
		Cconn: conn,
		Ctime: time.Now().Unix(),
	}
	s.workermu.Unlock()
	_, err := s.Conn.Write([]byte{com.OPEN, byte(cid >> 8), byte(cid & 0xff)})
	log.Println("New id", cid)
	if err != nil {
		log.Println("Send open new service conn is error.", err)
		return 0
	}
	return cid
}

// NewServiceConn 服务端新进来连接
func (s *Service) NewServiceConn(id uint16, conn net.Conn) {
	// 获取新连接id
	log.Println("Recv id", id)
	w := s.workers[id]
	if w == nil {
		log.Println("Connection has been deleted.", id)
		conn.Close()
		return
	}
	w.Sconn = conn
	go com.NetCopy(w.Cconn, w.Sconn, "Exchange conn", &count, &countmu)
	go com.NetCopy(w.Sconn, w.Cconn, "Exchange conn", &count, &countmu)
	s.workers[id] = nil
}
