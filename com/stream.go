package com

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"net"
	"sync"
)

// NStreamCrypt AES CTR加密算法
type NStreamCrypt struct {
	rstream cipher.Stream
	wstream cipher.Stream
}

// Init 初始化
func (my *NStreamCrypt) Init(key, iv []byte) {
	//指定加密、解密算法为AES，返回一个AES的Block接口对象
	rblock, rerr := aes.NewCipher(key)
	if rerr != nil {
		panic(rerr)
	}
	wblock, werr := aes.NewCipher(key)
	if werr != nil {
		panic(werr)
	}
	my.rstream = cipher.NewCTR(rblock, iv)
	my.wstream = cipher.NewCTR(wblock, iv)
}

// RCrypt 读时解密
func (my *NStreamCrypt) RCrypt(text []byte) {
	my.rstream.XORKeyStream(text, text)
}

// WCrypt 写时加密
func (my *NStreamCrypt) WCrypt(text []byte) {
	my.wstream.XORKeyStream(text, text)
}

// NCopy 封装后的新连接结构体
type NCopy struct {
	conn  net.Conn
	crypt *NStreamCrypt
}

// Init 初始化
func (my *NCopy) Init(conn net.Conn, key, iv []byte) {
	var c NStreamCrypt
	c.Init(key, iv)
	my.crypt = &c
	my.conn = conn
}

// Write 写入流时加密
func (my *NCopy) Write(p []byte) (n int, err error) {
	my.crypt.WCrypt(p)
	return my.conn.Write(p)
}

// Read 从流里面读时解密
func (my *NCopy) Read(p []byte) (n int, err error) {
	n, err = my.conn.Read(p)
	if n > 0 {
		my.crypt.RCrypt(p[:n])
	}
	return n, err
}

// Close 关闭流
func (my *NCopy) Close() error {
	return my.conn.Close()
}

// WCopy 写的一端加密，读不加密
func WCopy(dst *NCopy, src net.Conn, msg string, count *uint16, mu *sync.Mutex) {
	defer func() {
		src.Close()
		dst.Close()
		(*mu).Lock()
		(*count)--
		log.Println(msg, "-", *count)
		(*mu).Unlock()
	}()
	(*mu).Lock()
	(*count)++
	log.Println(msg, "+", *count)
	(*mu).Unlock()
	buf := make([]byte, 10240)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			(*dst).Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// RCopy 读的一端解密，写不加密
func RCopy(dst net.Conn, src *NCopy, msg string, count *uint16, mu *sync.Mutex) {
	defer func() {
		src.Close()
		dst.Close()
		(*mu).Lock()
		(*count)--
		log.Println(msg, "-", *count)
		(*mu).Unlock()
	}()
	(*mu).Lock()
	(*count)++
	log.Println(msg, "+", *count)
	(*mu).Unlock()
	buf := make([]byte, 10240)
	for {
		// (*src.conn).SetReadDeadline(time.Now().Add(time.Second * 10))
		n, err := (*src).Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// NetCopy 流复制处理
func NetCopy(dst, src net.Conn, msg string, count *uint16, mu *sync.Mutex) {
	defer func() {
		src.Close()
		dst.Close()
		(*mu).Lock()
		(*count)--
		log.Println(msg, "-", *count)
		(*mu).Unlock()
	}()
	(*mu).Lock()
	(*count)++
	log.Println(msg, "+", *count)
	(*mu).Unlock()
	buf := make([]byte, 10240)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			dst.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}
