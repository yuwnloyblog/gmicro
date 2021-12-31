package actorsystem

import (
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
	"google.golang.org/protobuf/proto"
)

type ActorRef interface {
	Tell(proto.Message, ActorRef)
	GetMethod() string
	GetHost() string
	GetPort() int
}

type DefaultActorRef struct {
	Host    string
	Port    int
	Method  string
	Session []byte
	Sender  *MsgSender
}

type DeadLetterActorRef struct {
	DefaultActorRef
}

func (ref *DeadLetterActorRef) Tell(message proto.Message, sender ActorRef) {
	//do nothing
}
func (ref *DeadLetterActorRef) GetMethod() string {
	return ref.Method
}
func (ref *DeadLetterActorRef) GetHost() string {
	return ref.Host
}
func (ref *DeadLetterActorRef) GetPort() int {
	return ref.Port
}

var NoSender *DeadLetterActorRef

func init() {
	NoSender = &DeadLetterActorRef{}
	NoSender.Host = "0.0.0.0"
	NoSender.Port = 0
}

func IsNoSender(req *rpc.RpcMessageRequest) bool {
	if req != nil {
		srcHost := req.SrcHost
		srcPort := req.SrcPort
		if srcHost == "0.0.0.0" && srcPort == 0 {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func NewActorRef(host string, port int, method string, session []byte, sender *MsgSender) ActorRef {
	ref := &DefaultActorRef{
		Host:    host,
		Port:    port,
		Method:  method,
		Session: session,
		Sender:  sender,
	}
	return ref
}

func (ref *DefaultActorRef) Tell(message proto.Message, sender ActorRef) {
	if message != nil {
		bytes, _ := proto.Marshal(message)

		rpcReq := &rpc.RpcMessageRequest{
			Session:   ref.Session,
			TarMethod: ref.Method,
			TarHost:   ref.Host,
			TarPort:   int32(ref.Port),

			SrcMethod: sender.GetMethod(),
			SrcHost:   sender.GetHost(),
			SrcPort:   int32(sender.GetPort()),

			Data: bytes,
		}
		ref.Sender.Send(rpcReq)
	}
}
func (ref *DefaultActorRef) GetMethod() string {
	return ref.Method
}
func (ref *DefaultActorRef) GetHost() string {
	return ref.Host
}
func (ref *DefaultActorRef) GetPort() int {
	return ref.Port
}
