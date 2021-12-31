package actorsystem

import "github.com/yuwnloyblog/gmicro/actorsystem/rpc"

type Processor func(ActorRef, interface{})
type NewInput func() interface{}

type Executor struct {
	CurrentCount int
	NewInputObj  NewInput
	Proc         Processor
}

/**
* TODO: Need asynchronous
**/
func (exe Executor) Execute(req *rpc.RpcMessageRequest, decoder ReqDecoder, encoder func(interface{}) []byte, msgSender *MsgSender) {
	var sender ActorRef

	srcHost := req.SrcHost
	srcPort := req.SrcPort
	srcMethod := req.SrcMethod
	srcSession := req.Session

	if IsNoSender(req) {
		sender = NoSender
	} else {
		sender = &DefaultActorRef{
			Host:    srcHost,
			Port:    int(srcPort),
			Method:  srcMethod,
			Session: srcSession,
			Encoder: encoder,
			Sender:  msgSender,
		}
	}

	bytes := req.Data
	input := exe.NewInputObj()
	decoder(bytes, input)
	exe.Proc(sender, input)
}
