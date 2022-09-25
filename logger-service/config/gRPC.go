package config

import (
	"fmt"
	"log"
	"logger/controllers"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

const (
	gRpcPort = "50001"
)

func GRPCListen(logServer controllers.LogServer) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &logServer)

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
