package service

import (
	"context"
	pb "github.com/xpoh/tages_test/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestServer_Login(t *testing.T) {
	var s Server
	go s.Run()
	time.Sleep(1 * time.Second)

	s.lg.AddUser("test", "test")

	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with cacher %v", err)
	}

	type args struct {
		ctx     context.Context
		request *pb.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.LoginResponse
		wantErr bool
	}{
		{
			name: "Zero test",
			args: args{
				ctx: context.Background(),
				request: &pb.LoginRequest{
					User: "test",
					Pass: "test",
				},
			},
			want:    &pb.LoginResponse{Token: "dGVzdHRlc3TaOaPuXmtLDTJVv--VYBiQr9gHCQ=="},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := pb.NewServiceClient(conn)
			got, err := client.Login(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Token, tt.want.Token) {
				t.Errorf("Login() got = %v, want %v", got.Token, tt.want.Token)
			}
		})
	}
}

//func TestServer_DownloadFile(t *testing.T) {
//
//	type fields struct {
//		UnimplementedServiceServer proto.UnimplementedServiceServer
//	}
//	type args struct {
//		request *pb.DownloadFileRequest
//		server  proto.Service_DownloadFileServer
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := Server{
//				UnimplementedServiceServer: tt.fields.UnimplementedServiceServer,
//			}
//			if err := s.DownloadFile(tt.args.request, tt.args.server); (err != nil) != tt.wantErr {
//				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestServer_GetFilesList(t *testing.T) {
//	type fields struct {
//		UnimplementedServiceServer proto.UnimplementedServiceServer
//	}
//	type args struct {
//		ctx     context.Context
//		request *pb.GetFilesListRequest
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *pb.GetFilesListResponse
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := Server{
//				UnimplementedServiceServer: tt.fields.UnimplementedServiceServer,
//			}
//			got, err := s.GetFilesList(tt.args.ctx, tt.args.request)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetFilesList() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetFilesList() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestServer_UploadFile(t *testing.T) {
//	type fields struct {
//		UnimplementedServiceServer proto.UnimplementedServiceServer
//	}
//	type args struct {
//		server proto.Service_UploadFileServer
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := Server{
//				UnimplementedServiceServer: tt.fields.UnimplementedServiceServer,
//			}
//			if err := s.UploadFile(tt.args.server); (err != nil) != tt.wantErr {
//				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
