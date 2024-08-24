package main

import (
	"net"

	"github.com/abdulshakoor02/ohlc_exinity/config"
	"github.com/abdulshakoor02/ohlc_exinity/database/dbAdapter"
	"github.com/abdulshakoor02/ohlc_exinity/database/migration"
	"github.com/abdulshakoor02/ohlc_exinity/logger"
	pb "github.com/abdulshakoor02/ohlc_exinity/ohlc"
	"github.com/abdulshakoor02/ohlc_exinity/service/aggregateData"
	"github.com/abdulshakoor02/ohlc_exinity/service/grpcServer"
	"google.golang.org/grpc"
)

func main() {
	log := logger.Logger
	config.LoadEnv()
	dbAdapter.DbConnect()
	migration.MigrateDb()

	ohlcServer := &grpcServer.Server{
		OhlcChannel: make(chan *pb.OHLC),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start the server")
	}
	log.Info().Msgf("server listening on port %v\n", lis.Addr())

	s := grpc.NewServer()

	pb.RegisterOHLCServiceServer(s, ohlcServer)
	go aggregateData.AggregateData(true, ohlcServer.OhlcChannel)
	go ohlcServer.ProcessOhlcData()

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
