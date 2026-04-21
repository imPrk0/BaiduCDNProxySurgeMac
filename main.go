package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// bufPool Buffer pool
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, bufSize)
		return &b
	},
}

func main() {
	port := flag.Int("p", 1080, "SOCKS5 Port (127.0.0.1:<port>)")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen error: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

// handle Socks5 社交的手腕
func handle(client net.Conn) {
	defer client.Close()

	hdr := make([]byte, 2)
	if _, err := io.ReadFull(client, hdr); err != nil {
		return
	}
	if hdr[0] != socks5Ver {
		return
	}
	methods := make([]byte, hdr[1])
	if _, err := io.ReadFull(client, methods); err != nil {
		return
	}
	if _, err := client.Write([]byte{socks5Ver, 0x00}); err != nil {
		return
	}

	req := make([]byte, 4)
	if _, err := io.ReadFull(client, req); err != nil {
		return
	}
	if req[0] != socks5Ver || req[1] != cmdCONN {
		client.Write([]byte{socks5Ver, 0x07, 0x00, atypIPv4, 0, 0, 0, 0, 0, 0})
		return
	}

	host, port, err := readAddr(client, req[3])
	if err != nil {
		client.Write([]byte{socks5Ver, 0x04, 0x00, atypIPv4, 0, 0, 0, 0, 0, 0})
		return
	}
	target := net.JoinHostPort(host, strconv.Itoa(port))

	up, err := dialUpstream(target)
	if err != nil {
		client.Write([]byte{socks5Ver, 0x04, 0x00, atypIPv4, 0, 0, 0, 0, 0, 0})
		return
	}
	defer up.Close()

	client.Write([]byte{socks5Ver, repOK, 0x00, atypIPv4, 0, 0, 0, 0, 0, 0})
	pipe(client, up)
}

// readAddr 获取网络地址
func readAddr(r io.Reader, atype byte) (string, int, error) {
	var host string
	switch atype {
	case atypIPv4:
		b := make([]byte, 4)
		if _, err := io.ReadFull(r, b); err != nil {
			return "", 0, err
		}
		host = net.IP(b).String()
	case atypIPv6:
		b := make([]byte, 16)
		if _, err := io.ReadFull(r, b); err != nil {
			return "", 0, err
		}
		host = net.IP(b).String()
	case atypDOMAIN:
		lenB := make([]byte, 1)
		if _, err := io.ReadFull(r, lenB); err != nil {
			return "", 0, err
		}
		domain := make([]byte, lenB[0])
		if _, err := io.ReadFull(r, domain); err != nil {
			return "", 0, err
		}
		host = string(domain)
	default:
		return "", 0, fmt.Errorf("unknown atype %d", atype)
	}
	portB := make([]byte, 2)
	if _, err := io.ReadFull(r, portB); err != nil {
		return "", 0, err
	}
	return host, int(binary.BigEndian.Uint16(portB)), nil
}

// dialUpstream 上游百度隧道
func dialUpstream(target string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", upstreamAddr, dialTimeout)
	if err != nil {
		return nil, err
	}

	req := fmt.Sprintf(
		"CONNECT %s HTTP/1.1\r\n"+
			"Host: %s\r\n"+
			"Proxy-Connection: Keep-Alive\r\n"+
			"X-T5-Auth: %s\r\n"+
			"User-Agent: %s\r\n\r\n",
		target, connectHostHdr, connectAuthHdr, connectUserAgent,
	)

	conn.SetDeadline(time.Now().Add(dialTimeout))
	if _, err := io.WriteString(conn, req); err != nil {
		conn.Close()
		return nil, err
	}

	br := bufio.NewReader(conn)
	resp, err := http.ReadResponse(br, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		conn.Close()
		return nil, fmt.Errorf("CONNECT %s: %s", target, resp.Status)
	}

	conn.SetDeadline(time.Time{})

	if br.Buffered() > 0 {
		return &prefixConn{Conn: conn, r: br}, nil
	}
	return conn, nil
}

type prefixConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *prefixConn) Read(b []byte) (int, error) {
	return c.r.Read(b)
}

// pipe 双向 Pipe
func pipe(a, b net.Conn) {
	done := make(chan struct{}, 2)
	cp := func(dst, src net.Conn) {
		bufp := bufPool.Get().(*[]byte)
		defer bufPool.Put(bufp)
		io.CopyBuffer(dst, src, *bufp)
		if tc, ok := dst.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
		done <- struct{}{}
	}
	go cp(a, b)
	go cp(b, a)
	<-done
	<-done
}
