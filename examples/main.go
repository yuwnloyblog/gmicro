package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yuwnloyblog/gmicro"
	"github.com/yuwnloyblog/gmicro/actorsystem"
	"github.com/yuwnloyblog/gmicro/utils"
	"google.golang.org/protobuf/proto"
)

func ReflectObj(value reflect.Value) interface{} {
	obj := reflect.New(value.Type()).Elem().Interface()
	return obj
}

func main() {
	fmt.Println("start")
	stu := utils.Student{
		Name: "abc",
		Age:  1,
	}

	intv := ReflectObj(reflect.ValueOf(stu))
	// bytes, _ := proto.Marshal(stu)

	// objValue := reflect.ValueOf(utils.Student{})

	// obj := reflect.New(objValue.Type()).Elem().Interface()
	// proto.Unmarshal(bytes, obj.(proto.Message))

	fmt.Println(intv.(proto.Message))
	// typ := reflect.ValueOf(utils.Student{}).Type()
	// fmt.Println(typ)

	// aa := reflect.New(typ).Elem().Interface()
	// fmt.Println(aa.(utils.Student).Age)

	// t := reflect.TypeOf(aa)
	// fmt.Println(t)
}

func ActorSystemTest() {
	actorSystem := actorsystem.NewActorSystemNoRpc("MyActorSystem")
	actorSystem.RegisterActorProcessor("m1", func() proto.Message {
		return &utils.Student{}
	}, MyProcessor, 10)

	actor := actorSystem.LocalActorOf("m1")
	actor.Tell(&utils.Student{
		Name: "name2",
		Age:  1,
	}, actorsystem.NoSender)
	time.Sleep(5 * time.Second)
}

func MyProcessor(sender actorsystem.ActorRef, input proto.Message) {
	fmt.Println("process has been executed.")
	fmt.Println("type:", reflect.TypeOf(input))
	stu := input.(*utils.Student)
	fmt.Println(stu.Name)

	sender.Tell(stu, actorsystem.NoSender)
}

func ClusterTest() {
	cluster := gmicro.NewCluster("cluster1", "127.0.0.1", 8888, []string{"127.0.0.1:2181"})
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
