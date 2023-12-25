package HengBus

import (
	"fmt"
	"net/rpc"
)

type Publisher struct {
	serverAddr string
	serverPath string
}

func NewPublisher(server_addr, server_path string) *Publisher {
	publisher := new(Publisher)
	publisher.serverAddr = server_addr
	publisher.serverPath = server_path
	return publisher
}

func (publisher *Publisher) Publish(topic string, args ...interface{}) error {
	fmt.Println("[Publisher]topic:", topic, "args:", args)
	rpcClient, err := rpc.DialHTTPPath("tcp", publisher.serverAddr, publisher.serverPath)
	defer rpcClient.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("dialing: %v", err))
	}
	reply := new(bool)
	handleArgs := &HandleArgs{args, topic}
	fmt.Printf("%#v\n", handleArgs)
	err = rpcClient.Call(PublishService, handleArgs, reply)
	if err != nil {
		fmt.Println(fmt.Errorf("Publish error: %v", err))
	}
	if !(*reply) {
		return fmt.Errorf("fail to publish")
	}
	return nil
}
