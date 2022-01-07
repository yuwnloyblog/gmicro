package gmicro

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/yuwnloyblog/gmicro/actorsystem"
	"github.com/yuwnloyblog/gmicro/utils"
	"google.golang.org/protobuf/proto"
)

func TestNewCluster(t *testing.T) {
	cluster := NewCluster("cluster1", "127.0.0.1", 8888, []string{"127.0.0.1:2181"})
	cluster.RegisterActorProcessor("m1", func() proto.Message {
		return &utils.Student{}
	}, MyProcessor, 10)

	stu := &utils.Student{
		Name: "name1",
		Age:  1,
	}

	time.Sleep(3 * time.Second)
	fmt.Println("Begin to Tell")
	cluster.UnicastRouteWithNoSender("m1", stu.Name, stu)

	time.Sleep(10 * time.Second)
}

func MyProcessor(sender actorsystem.ActorRef, input proto.Message) {
	fmt.Println("process has been executed.")
	fmt.Println("type:", reflect.TypeOf(input))
	stu := input.(*utils.Student)
	fmt.Println(stu.Name)

	sender.Tell(stu, actorsystem.NoSender)
}
