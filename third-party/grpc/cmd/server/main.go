package main

import (
	"context"
	pb "grpcpractice/protocol"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	// gRPC 提供的 default implementation, 只實現需要實現的方法
	// 避免未來不會因為缺少實作新方法導致編譯錯誤, 也增加一定程度的簡潔性
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	kaep := keepalive.EnforcementPolicy{
		MinTime:             1 * time.Minute, // keepalive ping 的最小間隔
		PermitWithoutStream: true,            // 允許客戶端不發送數據
	}

	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Minute, // 達到閒置時間上限時, server 主動關閉連線
		MaxConnectionAge:      30 * time.Minute, // 連線時間達到閥值時, server 主動關閉連線
		MaxConnectionAgeGrace: 5 * time.Minute,  // client 被標記成閒置時的寬限期
		Time:                  5 * time.Minute,  // server 主動發送心跳檢查的頻率
		Timeout:               1 * time.Minute,  // heartbeat timeout
	}

	s := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.UnaryInterceptor(serverInterceptor),
	)

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Interceptor 範例
//   - logging
//   - authentication
//   - execution time
func serverInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	log.Printf("Request method: %s", info.FullMethod)
	log.Printf("Request details: %v", req)

	md, ok := metadata.FromIncomingContext(ctx)

	if ok {
		tokens := md["token"]
		if len(tokens) > 0 && tokens[0] == "valid-token" {
			resp, err := handler(ctx, req)
			log.Printf("Request processed in %s", time.Since(start))
			return resp, err
		}
	}

	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
}
