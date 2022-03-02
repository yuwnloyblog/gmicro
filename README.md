# gmicro
一个使用zookeeper作为服务注册和发现，基于Actor模型的微服务框架。

#示例

```go

type MyActor struct {
	actorsystem.UntypedActor
}

func (act *MyActor) OnReceive(input proto.Message) {
	fmt.Println("process has been executed.")
	fmt.Println("type:", reflect.TypeOf(input))
	stu := input.(*utils.Student)
	fmt.Println(stu.Name)
	time.Sleep(3 * time.Second)
}
func (act *MyActor) CreateInputObj() proto.Message {
	return &utils.Student{}
}

func main() {
	cluster := gmicro.NewCluster("clusterName", "node1", "127.0.0.1", 9999, []string{"127.0.0.1:2181"}) //初始化一个集群对象，依赖ZK做注册发现

  //为当前节点注册一个处理任务的Actor，并指定并发处理数量为64， 其他节点可通过方法名"method1" 来调用到他
	cluster.RegisterActor("method1", func() actorsystem.IUntypedActor {
		return &MyActor{}
	}, 64)

	cluster.Startup()   //启动

  //在集群中，找到提供m1这个方法的所有节点，并以target_id来计算哈希值，以一致哈希的方式找到处理这个请求的节点，将数据发送给它
	cluster.UnicastRoute("m1", "target_id", &utils.Student{
		Name: "name1",
		Age:  12,
	}, actorsystem.NoSender)
}

```
