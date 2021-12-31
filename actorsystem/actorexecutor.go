package actorsystem

import (
	"github.com/yuwnloyblog/gmicro/actorsystem/rpc"
	"google.golang.org/protobuf/proto"
)

type Executor struct {
	CurrentCount int
	NewInputObj  NewInput
	Proc         Processor
}

/**
* TODO: Need asynchronous
**/
func (exe Executor) Execute(req *rpc.RpcMessageRequest, msgSender *MsgSender) {
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
			Sender:  msgSender,
		}
	}

	bytes := req.Data
	input := exe.NewInputObj()
	proto.Unmarshal(bytes, input)
	exe.Proc(sender, input)
}
