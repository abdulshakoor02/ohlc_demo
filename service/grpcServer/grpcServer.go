package grpcServer

import (
	"log"
	"sync"

	pb "github.com/abdulshakoor02/ohlc_exinity/ohlc"
)

type Server struct {
	pb.UnimplementedOHLCServiceServer
	OhlcChannel     chan *pb.OHLC
	Clients         map[string]chan *pb.OHLC
	ClientsMu       sync.Mutex
	CurrentOHLC     *pb.OHLC
	CurrentOHLCLock sync.RWMutex
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
	for ohlc := range s.OhlcChannel {
		if err := stream.Send(ohlc); err != nil {
			return err
		}
	}
	return nil
}
