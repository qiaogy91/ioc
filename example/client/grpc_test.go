package clinet_test

import (
	"context"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

type TokenAuth struct {
	token string
}

func (t *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	meta := map[string]string{"Token": t.token}
	return meta, nil
}

func (t *TokenAuth) RequireTransportSecurity() bool {
	return false
}

func TestCreate(t *testing.T) {
	ins, err := grpc.NewClient(
		"127.0.0.1:18080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),      // 非TLS
		grpc.WithPerRPCCredentials(&TokenAuth{token: "my-token-str"}), // 凭证
	)
	if err != nil {
		t.Logf("connect failed: %s", err)
		return
	}

	client := app01.NewServiceClient(ins)
	req := &app01.CreateUserReq{
		Username: "grpc_user04",
		Password: "redhat",
	}

	u, e := client.Create(context.Background(), req)
	if e != nil {
		t.Fatalf("create blog failed: %v", err)
		return
	}
	t.Logf("create user succeed: %+v", u)
}
