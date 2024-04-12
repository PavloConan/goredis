package main

import (
	"fmt"
	"log/slog"
	"net"
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
	fmt.Println(string(rawMsg))
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
	slog.Error("server error", "err", server.Start())
}
