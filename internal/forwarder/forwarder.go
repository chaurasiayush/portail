package forwarder

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/chaurasiayush/portail/internal/config"
)

func StartTCPForward(ctx context.Context, wg *sync.WaitGroup, rule config.ForwardRule) {
	defer wg.Done()

	listener, err := net.Listen("tcp", rule.Listen)
	if err != nil {
		log.Printf("[TCP] Listen error on %s: %v", rule.Listen, err)
		return
	}
	log.Printf("[TCP] Forwarding %s → %s", rule.Listen, rule.Forward)

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("[TCP] Accept error: %v", err)
				continue
			}
		}
		go handleTCPConnection(conn, rule)
	}
}

func handleTCPConnection(src net.Conn, rule config.ForwardRule) {
	defer src.Close()

	var dst net.Conn
	var err error

	if rule.TLS != nil && rule.TLS.Enabled {
		tlsCfg := &tls.Config{InsecureSkipVerify: rule.TLS.SkipVerify}
		dst, err = tls.Dial("tcp", rule.Forward, tlsCfg)
	} else {
		dst, err = net.Dial("tcp", rule.Forward)
	}

	if err != nil {
		log.Printf("[TCP] Dial error: %v", err)
		return
	}
	defer dst.Close()

	go io.Copy(dst, src)
	io.Copy(src, dst)
}

func StartUDPForward(ctx context.Context, wg *sync.WaitGroup, rule config.ForwardRule) {
	defer wg.Done()

	laddr, err := net.ResolveUDPAddr("udp", rule.Listen)
	if err != nil {
		log.Printf("[UDP] Resolve error: %v", err)
		return
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Printf("[UDP] Listen error: %v", err)
		return
	}
	defer conn.Close()

	raddr, err := net.ResolveUDPAddr("udp", rule.Forward)
	if err != nil {
		log.Printf("[UDP] Resolve forward error: %v", err)
		return
	}

	log.Printf("[UDP] Forwarding %s → %s", rule.Listen, rule.Forward)

	buf := make([]byte, 65535)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("[UDP] Read error: %v", err)
				continue
			}

			go func(data []byte, clientAddr *net.UDPAddr) {
				dstConn, err := net.DialUDP("udp", nil, raddr)
				if err != nil {
					log.Printf("[UDP] Dial error: %v", err)
					return
				}
				defer dstConn.Close()

				dstConn.Write(data)

				dstConn.SetReadDeadline(time.Now().Add(2 * time.Second))
				resp := make([]byte, 65535)
				n, _, err := dstConn.ReadFrom(resp)
				if err == nil {
					conn.WriteToUDP(resp[:n], clientAddr)
				}
			}(append([]byte(nil), buf[:n]...), clientAddr)
		}
	}
}
