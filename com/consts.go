package com

// Exchange接收指令
const (
	_ uint8 = iota
	// SERVICEC 服务端连接（控制连接）
	SERVICEC
	// SERVICED 服务端连接（数据端连接）
	SERVICED
	// CLIENTC 客户端连接（控制）
	CLIENTC
	// CLIENTD 客户端连接（数据）
	CLIENTD
	// OPEN 向服务端申请打开新连接
	OPEN
	// ERROR 错误
	ERROR
	// SUCCESS 成功
	SUCCESS
)

const (
	// BytePackLimit 数据包大小
	BytePackLimit = (1 << 16) - 10
	// TimeoutWite 连接超时时间180秒
	TimeoutWite = 180
	// Heartbeat 心跳时间
	Heartbeat = 30
)
