package main

import (
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan Message
}

func NewPeer(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])
		p.msgChan <- Message{
			data: msgBuf,
			peer: p,
		}
	}
}
