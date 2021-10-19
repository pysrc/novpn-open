package service

import "sync"

var (
	// OpenD 打开连接命令
	OpenD = make([]byte, 3)
	// 全局连接计数器
	count   uint16 = 0
	countmu sync.Mutex
	key, iv []byte
)
