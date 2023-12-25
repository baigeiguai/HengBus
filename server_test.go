package HengBus_test

import (
	HengBus "EventBus"
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	server := HengBus.NewServer(":2024", "/_server_")
	defer server.Stop()
	fmt.Println(server.Start(),'1')
}
