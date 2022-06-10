package service

import (
	"context"
	"encoding/json"
	"github.com/xpoh/tages_test/pkg/filestorage"
	"github.com/xpoh/tages_test/pkg/login"
	pb "github.com/xpoh/tages_test/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Server struct {
	gs *grpc.Server
	pb.UnimplementedServiceServer
	lg              login.ServiceLogin
	storage         filestorage.ImMemoryLocalStorage
	UpDownCounter   map[string]uint32
	ViewListCounter map[string]uint32
	Mux1            sync.Mutex
	Mux2            sync.Mutex
}

type maxUpDownError struct{}

func (m maxUpDownError) Error() string {
	return "Max count exceed"
}

func (s *Server) UploadFile(ctx context.Context, request *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	user := request.User
	token := request.Token

	if !s.lg.Auth(user, token) {
		return nil, login.NotAuthError{}
	}

	s.Mux1.Lock()
	cnt := s.UpDownCounter[user]
	if cnt >= 10 {
		s.Mux1.Unlock()
		return nil, maxUpDownError{}
	}
	s.UpDownCounter[user]++
	s.Mux1.Unlock()

	if err := s.storage.PutFile(request.User, request.Filename, request.Data); err != nil {
		return nil, err
	}
	return &pb.UploadFileResponse{Result: "Ok"}, nil
}

func (s *Server) DownloadFile(ctx context.Context, request *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	user := request.User
	token := request.Token

	if !s.lg.Auth(user, token) {
		return nil, login.NotAuthError{}
	}

	s.Mux1.Lock()
	cnt := s.UpDownCounter[user]
	if cnt >= 10 {
		s.Mux1.Unlock()
		return nil, maxUpDownError{}
	}
	s.UpDownCounter[user]++
	s.Mux1.Unlock()

	data, err := s.storage.GetFile(request.User, request.Filename)
	if err != nil {
		return nil, err
	}
	return &pb.DownloadFileResponse{File: data}, nil
}

func (s *Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.lg.GetToken(request.User, request.Pass)

	if err != nil {
		return nil, err
	}

	r := pb.LoginResponse{Token: token}
	return &r, nil
}

func (s *Server) GetFilesList(ctx context.Context, request *pb.GetFilesListRequest) (*pb.GetFilesListResponse, error) {
	user := request.User
	token := request.Token
	if !s.lg.Auth(user, token) {
		return nil, login.NotAuthError{}
	}

	s.Mux2.Lock()
	cnt := s.ViewListCounter[user]
	if cnt >= 100 {
		s.Mux2.Unlock()
		return nil, maxUpDownError{}
	}
	s.ViewListCounter[user]++
	s.Mux2.Unlock()

	list, _ := s.storage.GetFileList(request.User)
	marshal, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	r := pb.GetFilesListResponse{Files: string(marshal)}
	return &r, nil
}

func (s *Server) MustEmbedUnimplementedServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Run() {
	s.ViewListCounter = make(map[string]uint32)
	s.UpDownCounter = make(map[string]uint32)

	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc Server
	s.gs = grpc.NewServer()

	pb.RegisterServiceServer(s.gs, s)

	log.Println("start Server")
	// and start...
	if err := s.gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
