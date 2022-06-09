package service

import (
	"context"
	pb "github.com/xpoh/tages_test/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	pb.UnimplementedServiceServer
}

func (s Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	r := pb.LoginResponse{Token: request.User + request.Pass}
	return &r, nil
}

func (s Server) GetFilesList(ctx context.Context, request *pb.GetFilesListRequest) (*pb.GetFilesListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) UploadFile(server pb.Service_UploadFileServer) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) DownloadFile(request *pb.DownloadFileRequest, server pb.Service_DownloadFileServer) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) mustEmbedUnimplementedServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Run() {
	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc Server
	gs := grpc.NewServer()

	pb.RegisterServiceServer(gs, s)

	log.Println("start Server")
	// and start...
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
