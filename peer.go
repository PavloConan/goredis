package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
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
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, val := range v.Array() {
				switch val.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of arguments for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}

					p.msgChan <- Message{
						cmd:  cmd,
						peer: p,
					}

				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of arguments for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}

					p.msgChan <- Message{
						cmd:  cmd,
						peer: p,
					}
				}
			}
		}
	}

	return nil
}
