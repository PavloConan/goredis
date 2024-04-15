package main

import (
	"context"
	"fmt"
	"goredis/client"
	"log"
	"log/slog"
	"net"
	"time"
)

const defaultListenAddr = ":6969"

type Config struct {
	listenAddr string
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	peersChan chan *Peer
	quitChan  chan struct{}
	msgChan   chan []byte

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
		msgChan:   make(chan []byte),
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

	slog.Info("server running", "listenAddr", s.listenAddr)

	return s.listen()
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgChan:
			if err := s.handleRawMessage(rawMsg); err != nil {
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

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))

	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		s.kv.Get(v.key)
	}

	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	s.peersChan <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	client := client.NewClient("localhost:6969")

	for i := 0; i < 10; i++ {
		go func() {
			if err := client.Set(context.Background(), fmt.Sprintf("foo_%d", i), "bar"); err != nil {
				log.Fatal(err)
			}
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(server.kv.data)
	select {}
}
