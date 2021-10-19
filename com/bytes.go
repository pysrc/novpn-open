package com

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMd5 获取key的md5
func GetMd5(key string) []byte {
	d5 := md5.New()
	d5.Write([]byte(key))
	return d5.Sum(nil)
}

// GetHexString 获取16进制字符串
func GetHexString(b []byte) string {
	return hex.EncodeToString(b)
}

// GetKeyIv 通过提供的加密字符串通过md5计算出key iv
func GetKeyIv(passwd string) (key []byte, iv []byte) {
	var pb = []byte(passwd)
	var split = len(passwd) / 2
	var d5 = md5.New()
	d5.Write(pb[:split])
	key = d5.Sum(nil)
	d5 = md5.New()
	d5.Write(pb[split:])
	iv = d5.Sum(nil)
	return key, iv
}
