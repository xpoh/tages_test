package service

import (
	"context"
	"encoding/json"
	pb "github.com/xpoh/tages_test/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"os"
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
	s.gs.Stop()
}

func TestServer_UploadFile(t *testing.T) {
	var s Server
	go s.Run()
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	s.lg.AddUser("test", "test")

	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with service %v", err)
	}
	client := pb.NewServiceClient(conn)
	login, err := client.Login(ctx, &pb.LoginRequest{
		User: "test",
		Pass: "test",
	})
	if err != nil {
		log.Fatalf("can not login service %v", err)
	}
	token := login.Token

	testdata, err := os.ReadFile("../../test/add-to-desktop.png")
	if err != nil {
		t.Fatalf("Error read test data.")
	}
	type args struct {
		ctx     context.Context
		request *pb.UploadFileRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.UploadFileResponse
		wantErr bool
	}{
		{
			name: "Upload one file test",
			args: args{
				ctx: ctx,
				request: &pb.UploadFileRequest{
					User:     "test",
					Token:    token,
					Filename: "add-to-desktop.png",
					Data:     testdata,
				},
			},
			want:    &pb.UploadFileResponse{Result: "Ok"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.UploadFile(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
	s.gs.Stop()
}

func TestServer_DownloadFile(t *testing.T) {
	var s Server

	go s.Run()
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	s.lg.AddUser("test", "test")

	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with service %v", err)
	}
	client := pb.NewServiceClient(conn)
	login, err := client.Login(ctx, &pb.LoginRequest{
		User: "test",
		Pass: "test",
	})
	if err != nil {
		log.Fatalf("can not login service %v", err)
	}
	token := login.Token

	testdata, err := os.ReadFile("../../test/add-to-desktop.png")
	if err != nil {
		t.Fatalf("Error read test data.")
	}

	_, err = client.UploadFile(ctx, &pb.UploadFileRequest{
		User:     "test",
		Token:    token,
		Filename: "add-to-desktop.png",
		Data:     testdata,
	})
	if err != nil {
		t.Errorf("Error upload test file %v", err)
	}

	type args struct {
		ctx     context.Context
		request *pb.DownloadFileRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DownloadFileResponse
		wantErr bool
	}{
		{
			name: "Download one file test",
			args: args{
				ctx: ctx,
				request: &pb.DownloadFileRequest{
					User:     "test",
					Token:    token,
					Filename: "add-to-desktop.png",
				},
			},
			want:    &pb.DownloadFileResponse{File: testdata},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.DownloadFile(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.File, tt.want.File) {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
	s.gs.Stop()
}

func TestServer_GetFilesList(t *testing.T) {
	var s Server

	go s.Run()
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	s.lg.AddUser("test", "test")

	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with service %v", err)
	}
	client := pb.NewServiceClient(conn)
	login, err := client.Login(ctx, &pb.LoginRequest{
		User: "test",
		Pass: "test",
	})
	if err != nil {
		log.Fatalf("can not login service %v", err)
	}
	token := login.Token

	testdata, err := os.ReadFile("../../test/add-to-desktop.png")
	if err != nil {
		t.Fatalf("Error read test data.")
	}

	_, err = client.UploadFile(ctx, &pb.UploadFileRequest{
		User:     "test",
		Token:    token,
		Filename: "add-to-desktop.png",
		Data:     testdata,
	})
	if err != nil {
		t.Errorf("Error upload test file %v", err)
	}

	list, _ := s.storage.GetFileList("test")
	marshal, err := json.Marshal(list)
	if err != nil {
		t.Errorf("Fatal marshal json %v", err)
	}

	type args struct {
		ctx     context.Context
		request *pb.GetFilesListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetFilesListResponse
		wantErr bool
	}{
		{
			name: "GetFilesList test",
			args: args{
				ctx: ctx,
				request: &pb.GetFilesListRequest{
					User:  "test",
					Token: token,
				},
			},
			want:    &pb.GetFilesListResponse{Files: string(marshal)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := client.GetFilesList(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Files, tt.want.Files) {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
