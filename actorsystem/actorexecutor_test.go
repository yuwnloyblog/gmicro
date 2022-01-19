package actorsystem

import (
	"fmt"
	"time"

	"github.com/yuwnloyblog/gmicro/utils"
	"google.golang.org/protobuf/proto"
)

func generateObj() proto.Message {
	stu := &utils.Student{
		Name: "name1",
		Age:  1,
	}
	return stu
}

func process(sender ActorRef, msg proto.Message) {
	fmt.Println("execute")
	time.Sleep(3 * time.Second)
}

// func TestNewExecutor(t *testing.T) {
// 	executor := NewExecutor(3, generateObj, process)

// 	bytes, _ := utils.JsonMarshal(generateObj())

// 	req := &rpc.RpcMessageRequest{
// 		SrcMethod: "srcMethod",
// 		SrcHost:   "127.0.0.1",
// 		SrcPort:   8888,
// 		Data:      bytes,
// 	}

// 	for i := 0; i < 100; i++ {
// 		fmt.Println("begin num:", i)
// 		executor.Execute(req, nil)
// 	}
// }
