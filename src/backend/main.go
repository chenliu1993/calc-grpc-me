package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/chenliu1993/calc-grpc-me/proto/gen"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	port := 8000
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger :%v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(logger),
	))
	grpc_zap.ReplaceGrpcLogger(logger)
	pb.RegisterCalcServer(server, &CalcService{})
	server.Serve(lis)
}

type CalcService struct{}

func (s *CalcService) Increment(ctx context.Context, req *pb.NumRequest) (*pb.NumResponse, error) {
	req.Val++
	return &pb.NumResponse{Val: req.Val}, nil
}
