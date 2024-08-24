package main

import (
	"log"
	"net"
	"sync"

	"github.com/abdulshakoor02/ohlc_exinity/config"
	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/database/migration"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
	pb "github.com/abdulshakoor02/ohlc_exinity/ohlc"
	"github.com/abdulshakoor02/ohlc_exinity/service/aggregateData"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedOHLCServiceServer
	ohlcChannel     chan *pb.OHLC
	clients         map[string]chan *pb.OHLC
	clientsMu       sync.Mutex
	currentOHLC     *pb.OHLC
	currentOHLCLock sync.RWMutex
}

func (s *server) processOhlcData() {
	for ohlc := range s.ohlcChannel {
		s.currentOHLCLock.Lock()
		s.currentOHLC = ohlc
		s.currentOHLCLock.Unlock()

		s.broadcastToClients(ohlc)
	}
}

func (s *server) broadcastToClients(ohlc *pb.OHLC) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for _, ch := range s.clients {
		select {
		case ch <- ohlc:
		default:
			log.Println("Dropped OHLC data due to slow consumer")
		}
	}
}

func (s *server) StreamOHLCData(
	req *pb.OHLCrequest,
	stream pb.OHLCService_StreamOHLCDataServer,
) error {
	for ohlc := range s.ohlcChannel {
		if err := stream.Send(ohlc); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	log := logger.Logger
	config.LoadEnv()
	dbAdapter.DbConnect()
	migration.MigrateDb()

	ohlcServer := &server{
		ohlcChannel: make(chan *pb.OHLC),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start the server")
	}
	log.Info().Msgf("server listening on port %v\n", lis.Addr())

	s := grpc.NewServer()

	pb.RegisterOHLCServiceServer(s, ohlcServer)
	go aggregateData.AggregateData(true, ohlcServer.ohlcChannel)
	// go aggregateData.AggregateData("ethusdt", true, ohlcServer.ohlcChannel)
	go ohlcServer.processOhlcData()
	// go aggregateData.AggregateData("pepeusdt", true)
	// aggregateData.AggregateData("btcusdt", true)

	err = s.Serve(lis)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start grpc server")
	}

	log.Info().Msg("grpc server started")

	// conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("failed to start grpc server")
	// }
	// defer conn.Close()
	//
	// client := pb.NewOHLCServiceClient(conn)
	// stream, err := client.StreamOHLCData(context.Background(), &pb.Empty{})
	// if err != nil {
	// 	log.Fatal().Err(err).Msgf("failed to start grpc server")
	// }
	//
	// for {
	// 	ohlc, err := stream.Recv()
	// 	if err != nil {
	// 		log.Fatal().Err(err).Msgf("failed to start grpc server")
	// 	}
	// 	log.Printf(
	// 		"Received OHLC: OpenTime: %s, CloseTime: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %.4f",
	// 		ohlc.OpenTime,
	// 		ohlc.CloseTime,
	// 		ohlc.Open,
	// 		ohlc.High,
	// 		ohlc.Low,
	// 		ohlc.Close,
	// 	)
	// }
}
