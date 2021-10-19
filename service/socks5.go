package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"novpn/com"
	"time"
)

// socks5协议解析
func socks5(conn net.Conn) error {
	var s com.NCopy
	s.Init(conn, key, iv)
	buf := make([]byte, 258)
	// 协商加密算法
	if _, err := io.ReadAtLeast(&s, buf[:2], 2); err != nil {
		return err
	}
	if buf[0] != 0x05 {
		return errors.New("Proxy version error")
	}
	nm := int(buf[1])
	if _, err := io.ReadAtLeast(&s, buf[2:2+nm], nm); err != nil {
		return err
	}
	if _, err := (&s).Write([]byte{0x05, 0x00}); err != nil {
		return err
	}
	// 开始解析目标地址
	var host string
	var port uint16
	buf = make([]byte, 258)
	if _, err := io.ReadAtLeast(&s, buf[:5], 5); err != nil {
		return err
	}
	if buf[0] != 0x05 || buf[1] != 0x01 {
		return errors.New("Proxy version error")
	}
	switch buf[3] {
	case 0x01:
		// ipv4
		if _, err := io.ReadAtLeast(&s, buf[5:6+net.IPv4len], net.IPv4len+1); err != nil {
			return err
		}
		host = net.IP(buf[4 : 4+net.IPv4len]).String()
		port = (uint16(buf[4+net.IPv4len]) << 8) | uint16(buf[5+net.IPv4len])
	case 0x03:
		// domain
		dlen := int(buf[4])
		if _, err := io.ReadAtLeast(&s, buf[5:7+dlen], 2+dlen); err != nil {
			return err
		}
		host = string(buf[5 : 5+dlen])
		port = (uint16(buf[5+dlen]) << 8) | uint16(buf[6+dlen])
	case 0x04:
		// ipv6
		if _, err := io.ReadAtLeast(&s, buf[5:6+net.IPv6len], net.IPv6len+1); err != nil {
			return err
		}
		host = net.IP(buf[4 : 4+net.IPv6len]).String()
		port = (uint16(buf[4+net.IPv6len]) << 8) | uint16(buf[5+net.IPv6len])
	}
	// 处理本地映射
	if host == "remotehost" {
		host = "127.0.0.1"
	}
	host = fmt.Sprintf("%v:%v", host, port)
	log.Println("HOST:", host)
	rconn, err := net.DialTimeout("tcp", host, time.Duration(time.Second*10))
	if err != nil {
		return err
	}
	// 回应socks5
	buf = make([]byte, 258)
	tcpAddr := rconn.LocalAddr().(*net.TCPAddr)
	if tcpAddr.Zone == "" {
		if tcpAddr.IP.Equal(tcpAddr.IP.To4()) {
			tcpAddr.Zone = "ip4"
		} else {
			tcpAddr.Zone = "ip6"
		}
	}

	buf[0] = 0x05
	buf[1] = 0x00
	buf[2] = 0x00
	var ip net.IP
	if "ip6" == tcpAddr.Zone {
		ip = tcpAddr.IP.To16()
		buf[3] = 0x04
	} else {
		buf[3] = 0x01
		ip = tcpAddr.IP.To4()
	}
	pindex := 4
	for _, b := range ip {
		buf[pindex] = b
		pindex++
	}
	buf[pindex] = byte((tcpAddr.Port >> 8) & 0xff)
	buf[pindex+1] = byte(tcpAddr.Port & 0xff)
	(&s).Write(buf[:pindex+2])
	// socks5回应完成，交换数据
	go com.WCopy(&s, rconn, "Service conn", &count, &countmu)
	go com.RCopy(rconn, &s, "Service conn", &count, &countmu)
	// go com.NetCopy(conn, rconn, "Service conn", &count, &countmu)
	// go com.NetCopy(rconn, conn, "Service conn", &count, &countmu)
	return nil
}
