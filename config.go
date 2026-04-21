package main

import "time"

const (
	// upstreamAddr Baidu CDN 网络代理连接地址
	upstreamAddr = "cloudnproxy.baidu.com:443"

	// connectUserAgent Baidu CDN UserAgent
	connectUserAgent = "okhttp/3.11.0 SP-engine/2.71.0 Dalvik/2.1.0 (Linux; U; Android 9; HMA-AL00 Build/PQ3B.190801.002) baiduboxapp/13.33.0.11 (Baidu; P1 9)"

	// connectHostHdr Baidu Host Header
	connectHostHdr = "ascdn.baidu.com"

	// connectAuthHdr Auth Header
	connectAuthHdr = "1951164069"

	dialTimeout = 10 * time.Second
	ioTimeout   = 60 * time.Second
	bufSize     = 32 * 1024 // 32 KB copy buffer

	socks5Ver  = 0x05
	cmdCONN    = 0x01
	atypIPv4   = 0x01
	atypDOMAIN = 0x03
	atypIPv6   = 0x04
	repOK      = 0x00
)
