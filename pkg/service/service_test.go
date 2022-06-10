package service

import (
	"context"
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
			name: "Zero test",
			args: args{
				ctx: ctx,
				request: &pb.UploadFileRequest{
					User:     "test",
					Token:    token,
					Filename: "",
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
}
