package HengBus_test

import (
	HengBus "EventBus"
	"testing"
)

func TestPublisher(t *testing.T) {
	publish := HengBus.NewPublisher(":2024","/_server_")
	publish.Publish("main:calculator",1,2)
	publish.Publish("main:minus",1,2)
	publish.Publish("main:multiply",1,2)
}
