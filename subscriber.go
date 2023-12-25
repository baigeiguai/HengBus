package HengBus

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

const (
	// PublishService - Client service method
	HandleService = "SubscriberService.HandleEvent"
)

type HandleArgs struct {
	Args  []interface{}
	Topic string
}

type Subscriber struct {
	eventBus Bus
	address  string
	path     string
	serverPath string
	serverAddr string
	service  *SubscriberService
}

func NewSubscriber(address, path,server_add,server_path string, eventBus Bus) *Subscriber {
	subscriber := new(Subscriber)
	subscriber.address = address
	subscriber.path = path
	subscriber.serverAddr = server_add
	subscriber.serverPath = server_path
	subscriber.eventBus = eventBus
	subscriber.service = &SubscriberService{subscriber, &sync.WaitGroup{}, false}
	return subscriber
}

func (subscriber *Subscriber) Start() error {
	var err error
	service := subscriber.service
	if !service.started {
		server := rpc.NewServer()
		server.Register(service)
		server.HandleHTTP(subscriber.path, "/debug"+subscriber.path)
		l, err := net.Listen("tcp", subscriber.address)
		if err == nil {
			service.wg.Add(1)
			service.started = true
			go http.Serve(l, nil)
		}
	} else {
		err = errors.New("Client service already started")
	}
	return err
}

func (subscriber *Subscriber) doSubscribe(topic string, fn interface{}) {
	fmt.Println("[Publisher]topic:", topic, "function:", fn)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Server not found -", r)
		}
	}()
	rpcClient, err := rpc.DialHTTPPath("tcp", subscriber.serverAddr, subscriber.serverPath)
	defer rpcClient.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("dialing: %v", err))
	}
	args := &SubscribeArg{subscriber.address, subscriber.path, HandleService, topic}
	reply := new(bool)
	err = rpcClient.Call(RegisterService, args, reply)
	if err != nil {
		fmt.Println(fmt.Errorf("Register error: %v", err))
	}
	if *reply {
		subscriber.eventBus.Subscribe(topic, fn)
	}
}

func (subscriber *Subscriber) Subscribe(topic string, fn interface{}) {
	subscriber.doSubscribe(topic, fn)
}

func (subscriber *Subscriber)Stop(){
	service:=subscriber.service
	if service.started{
		service.wg.Done()
		service.started =false
	}
}

type SubscriberService struct {
	subscriber *Subscriber
	wg         *sync.WaitGroup
	started    bool
}

func (service *SubscriberService) HandleEvent(arg *HandleArgs, reply *bool) error {
	service.subscriber.eventBus.Publish(arg.Topic, arg.Args...)
	*reply = true
	return nil
}
