package HengBus_test

import (
	HengBus "EventBus"
	"fmt"
	"testing"
)

func MyPrint(a int) {
	fmt.Printf("MyPrint:%d\n", a)
}

func TestBaseBus(t *testing.T) {
	bus := HengBus.New()
	bus.Subscribe("main:Print", MyPrint)
	bus.Publish("main:Print", 20)
	bus.Unsubscribe("main:Print", MyPrint)
}
