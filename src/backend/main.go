package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

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

	mux := runtime.NewServeMux()
	if err := pb.RegisterCalcHandlerFromEndpoint(context.Background(), mux, fmt.Sprintf("localhost:%d", port), []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		logger.Fatal("http reg extablished failed ", zap.Error(err))
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", rest),
		Handler: moveToBody(mux),
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

func moveToBody(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Has("val") {
			val, err := strconv.Atoi(r.URL.Query().Get("val"))
			if err != nil {
				log.Fatal(err)
			}
			log.Println(val)
			req, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewReader([]byte(fmt.Sprintf("{\"val\":\"%d\"}", val))))
			if err != nil {
				log.Fatal(err)
			}
			f.ServeHTTP(w, req)
		} else {
			f.ServeHTTP(w, r)
		}
	})
}
