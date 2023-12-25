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
	RegisterService = "ServerService.Register"
	PublishService  = "ServerService.PushEvent"
	UnRegisterService = "ServerService.UnRegister"
)

type SubscribeArg struct {
	ClientAddr    string
	ClientPath    string
	ServiceMethod string
	Topic         string
	// SubscribeType SubscribeType
}

type SubscriberHandler struct {
	subArg SubscribeArg
	sync.Mutex
}

type Server struct {
	eventBus    Bus
	address     string
	path        string
	subscribers map[string][]*SubscriberHandler
	service     *ServerService
	lock        sync.Mutex
}

func (server *Server) EventBus() Bus {
	return server.eventBus
}

func NewServer(address, path string) *Server {
	server := new(Server)
	server.eventBus = New()
	server.address = address
	server.path = path
	server.subscribers = make(map[string][]*SubscriberHandler)
	server.service = &ServerService{server, &sync.WaitGroup{}, false}
	return server
}

func (server *Server) Start() error {
	var err error
	service := server.service
	if !service.started {
		rpcServer := rpc.NewServer()
		rpcServer.Register(service)
		rpcServer.HandleHTTP(server.path, "/debug"+server.path)
		l, e := net.Listen("tcp", server.address)
		if e != nil {
			err = e
			fmt.Println(fmt.Errorf("listen error: %v", e))
		}
		service.started = true
		service.wg.Add(1)
		http.Serve(l, nil)
	} else {
		err = errors.New("Server bus already started")
	}
	return err
}

func (server *Server) HasClientSubscribed(arg *SubscribeArg) bool {
	if topicSubscribers, ok := server.subscribers[arg.Topic]; ok {
		for _, topicSubscriber := range topicSubscribers {
			if topicSubscriber.subArg == *arg {
				return true
			}
		}
	}
	return false
}

func (server *Server) rpcCallback(subscribeArg *SubscribeArg) func(args ...interface{}) {
	return func(args ...interface{}) {
		client, connErr := rpc.DialHTTPPath("tcp", subscribeArg.ClientAddr, subscribeArg.ClientPath)
		defer client.Close()
		if connErr != nil {
			fmt.Println(fmt.Errorf("dialing: %v", connErr))
		}
		clientArg := new(HandleArgs)
		clientArg.Topic = subscribeArg.Topic
		clientArg.Args = args
		var reply bool
		err := client.Call(subscribeArg.ServiceMethod, clientArg, &reply)
		if err != nil {
			fmt.Println(fmt.Errorf("dialing: %v", err))
		}
	}
}

func (server *Server) Stop() {
	service := server.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}

type ServerService struct {
	server  *Server
	wg      *sync.WaitGroup
	started bool
}

func (service *ServerService) Register(arg *SubscribeArg, success *bool) error {
	fmt.Printf("[Server-Subscriber]arg:%#v\n",arg)
	subscribers := service.server.subscribers
	service.server.lock.Lock()
	defer service.server.lock.Unlock()
	if !service.server.HasClientSubscribed(arg) {
		rpcCallback := service.server.rpcCallback(arg)
		service.server.eventBus.Subscribe(arg.Topic, rpcCallback)
		var topicSubscribers []*SubscriberHandler
		topicSubscriber := &SubscriberHandler{*arg, sync.Mutex{}}
		if _, ok := subscribers[arg.Topic]; !ok {
			topicSubscribers = []*SubscriberHandler{topicSubscriber}
		} else {
			topicSubscribers = subscribers[arg.Topic]
			topicSubscribers = append(topicSubscribers, topicSubscriber)
		}
		subscribers[arg.Topic] = topicSubscribers
	}
	*success = true
	return nil
}

func (service *ServerService) PushEvent(handleArgs *HandleArgs, success *bool) error {
	fmt.Printf("[Server-Publisher]arg:%#v\n", handleArgs)
	service.server.EventBus().Publish(handleArgs.Topic, handleArgs.Args...)
	*success = true
	return nil
}
