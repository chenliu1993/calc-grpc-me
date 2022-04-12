package main

import (
	"fmt"
	"net"
	"net/http"

	pb "github.com/chenliu1993/calc-grpc-me/proto"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	port := 8000
	rest := 8001
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("failed to create logger :", zap.Error(err))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_zap.UnaryServerInterceptor(logger),
	))
	grpc_zap.ReplaceGrpcLogger(logger)
	pb.RegisterCalcServer(server, &CalcService{})
	go func() {
		if err := server.Serve(lis); err != nil {
			logger.Fatal("server exists", zap.Error(err))
		}
	}()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		logger.Fatal("http conn established failed ", zap.Error(err))
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterCalcHandler(context.Background(), mux, conn); err != nil {
		logger.Fatal("http reg extablished failed ", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", rest),
		Handler: mux,
	}

	logger.Fatal("http listen failed", zap.String("msg", httpServer.ListenAndServe().Error()))
}

type CalcService struct{}

func (s *CalcService) Increment(ctx context.Context, req *pb.NumRequest) (*pb.NumResponse, error) {
	req.Val++
	return &pb.NumResponse{Val: req.Val}, nil
}

func (s *CalcService) Work(ctx context.Context, req *pb.WorkRequest) (*pb.WorkResponse, error) {
	return &pb.WorkResponse{
		Reply: "This is boring",
	}, nil
}
