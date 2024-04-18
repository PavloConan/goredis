package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultListenAddr = ":6969"

type Config struct {
	listenAddr string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	peersChan chan *Peer
	quitChan  chan struct{}
	msgChan   chan Message

	kv *KeyValStore
}

func NewServer(cfg Config) *Server {
	if len(cfg.listenAddr) == 0 {
		cfg.listenAddr = defaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		peersChan: make(chan *Peer),
		quitChan:  make(chan struct{}),
		msgChan:   make(chan Message),
		kv:        NewKeyValStore(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)

	if err != nil {
		return nil
	}
	s.ln = ln

	go s.loop()

	slog.Info("goredis server running", "listenAddr", s.listenAddr)

	return s.listen()
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case peer := <-s.peersChan:
			s.peers[peer] = true
		case <-s.quitChan:
			return
		}
	}
}

func (s *Server) listen() error {
	for {
		conn, err := s.ln.Accept()

		if err != nil {
			slog.Error("server error", "err", err)
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}

		_, err := msg.peer.Send(val)

		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	}

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	s.peersChan <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of goredis server")
	server := NewServer(Config{
		listenAddr: *listenAddr,
	})

	log.Fatal(server.Start())
}
