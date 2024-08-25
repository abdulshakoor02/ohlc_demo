package grpcServer

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/abdulshakoor02/ohlc_exinity/logger"
	pb "github.com/abdulshakoor02/ohlc_exinity/ohlc"
)

var log = logger.Logger

type Server struct {
	pb.UnimplementedOHLCServiceServer
	OhlcChannel      chan *pb.OHLC
	CurrentOHLCs     map[string]*pb.OHLC
	Clients          map[string]chan *pb.OHLC
	ClientsMu        sync.Mutex
	CurrentOHLC      *pb.OHLC
	CurrentOHLCLock  sync.RWMutex
	ClientCount      int
	ClientCountMutex sync.Mutex
}

func (s *Server) ProcessOhlcData() {
	for ohlc := range s.OhlcChannel {
		s.CurrentOHLCLock.Lock()
		s.CurrentOHLC = ohlc
		s.CurrentOHLCLock.Unlock()

		s.BroadcastToClients(ohlc)
	}
}

func (s *Server) BroadcastToClients(ohlc *pb.OHLC) {
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()

	for _, ch := range s.Clients {
		select {
		case ch <- ohlc:
		default:
			log.Println("Dropped OHLC data due to slow consumer")
		}
	}
}

func (s *Server) StreamOHLCData(
	req *pb.OHLCrequest,
	stream pb.OHLCService_StreamOHLCDataServer,
) error {
	clientID := fmt.Sprintf("%d", time.Now().UnixNano())
	log.Info().Msgf("Requested trading pair for client id %v: %v\n", clientID, req.TradePair)
	clientChannel := make(chan *pb.OHLC, 10)
	s.AddClient(clientID, clientChannel)
	defer s.RemoveClient(clientID)

	for {
		select {
		case ohlc, ok := <-clientChannel:
			if !ok {
				return nil
			}
			data := ohlc
			if data.TradePair == strings.ToUpper(req.TradePair) {
				if err := stream.Send(data); err != nil {
					log.Println("Error sending OHLC data:", err)
					return err
				}
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *Server) AddClient(clientID string, ch chan *pb.OHLC) {
	log.Info().Msgf("adding client %v\n", clientID)
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()
	s.Clients[clientID] = ch
	s.IncrementClientCount()
	log.Info().Msgf("added client %v\n", clientID)
}

func (s *Server) RemoveClient(clientID string) {
	log.Info().Msgf("removing client %v\n", clientID)
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()
	log.Info().Msgf("removed client %v\n", clientID)
	if ch, ok := s.Clients[clientID]; ok {
		close(ch)
		delete(s.Clients, clientID)
		s.DecrementClientCount()
	}
}

func (s *Server) IncrementClientCount() {
	s.ClientCountMutex.Lock()
	defer s.ClientCountMutex.Unlock()
	s.ClientCount++
	log.Info().Msgf("No of clients connected: %v", s.ClientCount)
	// if s.ClientCount == 1 {
	// Start WebSocket stream when the first client connects
	// s.StartWebSocket <- true
	// }
}

func (s *Server) DecrementClientCount() {
	s.ClientCountMutex.Lock()
	defer s.ClientCountMutex.Unlock()
	s.ClientCount--
	log.Info().Msgf("No of clients connected: %v", s.ClientCount)
	// if s.ClientCount == 0 {
	// Stop WebSocket stream when the last client disconnects
	// s.stopWebSocket <- true
	// }
}
