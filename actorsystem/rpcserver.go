package actorsystem

import (
	"bot/pkg/gmicro/actorsystem/rpc"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type RpcServer struct {
	Host     string
	Port     int
	receiver *MsgReceiver
}

func NewRpcServer(host string, port int, msgReceiver *MsgReceiver) *RpcServer {
	server := &RpcServer{
		Host:     host,
		Port:     port,
		receiver: msgReceiver,
	}
	return server
}

func (server *RpcServer) Send(ctx context.Context, in *rpc.RpcMessageRequest) (*rpc.RpcMessageResponse, error) {
	fmt.Println("receive : ", in)
	server.receiver.Receive(in)
	return &rpc.RpcMessageResponse{Status: 0}, nil
}

func (server *RpcServer) Start() {
	hostStr := server.Host
	port := server.Port
	portStr := strconv.Itoa(port)

	lis, err := net.Listen("tcp", strings.Join([]string{hostStr, portStr}, ":"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterRpcMessageServer(s, server)
	log.Println("rpc started at (" + hostStr + ":" + portStr + ")!")
	s.Serve(lis)
}
