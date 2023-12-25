package HengBus_test

import (
	HengBus "EventBus"
	"fmt"
	"testing"
)

func calculator(a, b int) {
	fmt.Printf("%d+%d=%d\n", a, b, a+b)
}

func multiply(a, b int) {
	fmt.Printf("%d*%d=%d\n", a, b, a*b)
}

func minus(a, b int) {
	fmt.Printf("%d-%d=%d\n", a, b, a-b)
}

func TestSubscriber1(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Server not found -", r)
		}
	}()
	subscriber := HengBus.NewSubscriber(":2023", "/_subscriber_", ":2024", "/_server_", HengBus.New())
	defer subscriber.Stop()
	fmt.Println(subscriber)
	fmt.Println(subscriber.Start())
	subscriber.Subscribe("main:calculator", calculator)
	subscriber.Subscribe("main:minus", minus)
	for true {
		;
	}
}

func TestSubscriber2(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Server not found -", r)
		}
	}()
	subscriber := HengBus.NewSubscriber(":2026", "/_subscriber_", ":2024", "/_server_", HengBus.New())
	defer subscriber.Stop()
	fmt.Println(subscriber)
	fmt.Println(subscriber.Start())
	subscriber.Subscribe("main:multiply", multiply)
	subscriber.Subscribe("main:minus", minus)
	for true {
		;
	}

}
