package actorsystem

import (
	context "context"
	"fmt"
	"log"

	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"

	grpc "google.golang.org/grpc"
)

type RpcClient struct {
	Address     string
	conn        *grpc.ClientConn
	msgClient   rpc.RpcMessageClient
	isConnected bool
}

func NewRpcClient(address string) *RpcClient {
	client := &RpcClient{
		Address:     address,
		isConnected: false,
	}

	return client
}

func (client *RpcClient) connect() {
	tmpConn, err := grpc.Dial(client.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		client.isConnected = true
	}
	client.conn = tmpConn
	client.msgClient = rpc.NewRpcMessageClient(tmpConn)
}

func (client *RpcClient) DisConnect() {
	client.conn.Close()
}

func (client *RpcClient) Send(req *rpc.RpcMessageRequest) {
	if !client.isConnected {
		client.connect()
	}
	resp, err := client.msgClient.Send(context.Background(), req)
	if err != nil {
		fmt.Println("resp:", resp, "err:", err)
	}
}
