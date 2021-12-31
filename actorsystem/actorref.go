package actorsystem

import "github.com/yuwnloyblog/gmicro/actorsystem/rpc"

type ActorRef interface {
	Tell(interface{}, ActorRef)
	GetMethod() string
	GetHost() string
	GetPort() int
}

type DefaultActorRef struct {
	Host    string
	Port    int
	Method  string
	Session []byte
	Encoder func(interface{}) []byte
	Sender  *MsgSender
}

type DeadLetterActorRef struct {
	DefaultActorRef
}

func (ref *DeadLetterActorRef) Tell(message interface{}, sender ActorRef) {
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

func NewActorRef(host string, port int, method string, session []byte, encoder func(interface{}) []byte, sender *MsgSender) ActorRef {
	ref := &DefaultActorRef{
		Host:    host,
		Port:    port,
		Method:  method,
		Session: session,
		Encoder: encoder,
		Sender:  sender,
	}
	return ref
}

func (ref *DefaultActorRef) Tell(message interface{}, sender ActorRef) {
	if message != nil {
		bytes := ref.Encoder(message)

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
