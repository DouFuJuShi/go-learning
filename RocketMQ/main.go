package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"time"
)

// 积分服务使用
type OtherListener struct {
}

func (o *OtherListener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地业务逻辑入库")
	time.Sleep(5 * time.Second)
	fmt.Println("本地业务逻辑入库成功")
	return primitive.UnknowState
}
func (o *OtherListener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	// 注意，rocketmq回调时，不会调用到这里，只会调用到上一层的 Listener 里的 CheckLocalTransaction 方法
	return primitive.CommitMessageState
}

// 库存服务使用
type Listener struct{}

func (l *Listener) ExecuteLocalTransaction(message *primitive.Message) primitive.LocalTransactionState {
	// 调用积分的half
	o := OtherListener{}
	p, err := rocketmq.NewTransactionProducer(
		&o,
		producer.WithNameServer([]string{"127.0.0.1:9876"}),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	// 注意，这里不能使用 p.Shutdown()
	//这会导致所有都与borker的连接关闭，引发很多问题
	// 外层的不是同一个producer也会被影响到
	msg := &primitive.Message{
		Topic: "transTopicJifen",
		Body:  []byte("Jifen!"),
	}
	// 1. SendMessageInTransaction 阻塞，然后执行 ExecuteLocalTransaction
	// 2. ExecuteLocalTransaction 执行结束后 SendMessageInTransaction 解除阻塞
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("积分发送成功,data:%s,msgID:=%s\n", res.String(), res.MsgID)
	return primitive.UnknowState
}

// 嵌套事务的所有rocketmq回调都会调用到这里的 CheckLocalTransaction
func (l *Listener) CheckLocalTransaction(ext *primitive.MessageExt) primitive.LocalTransactionState {
	// 根据ext自行区分topic，处理业务逻辑
	fmt.Printf("库存，收到Rocketmq主动请求信息，msg:%+v\n", *ext)
	return primitive.CommitMessageState
}

func main() {
	l := Listener{}
	p, err := rocketmq.NewTransactionProducer(
		&l,
		producer.WithNameServer([]string{"127.0.0.1:9876"}),
	)
	if err != nil {
		panic(err)
	}
	err = p.Start()
	defer func() {
		err = p.Shutdown()
		if err != nil {
			fmt.Printf("shutdown producer error: %s", err.Error())
		}
	}()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	msg := &primitive.Message{
		Topic: "transTopic",
		Body:  []byte("kucun!"),
	}
	// 1. SendMessageInTransaction 阻塞，然后执行 ExecuteLocalTransaction
	// 2. ExecuteLocalTransaction 执行结束后 SendMessageInTransaction 解除阻塞
	res, err := p.SendMessageInTransaction(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("库存发送成功,data:%s,msgID:=%s\n", res.String(), res.MsgID)
	<-(chan interface{})(nil)
}
