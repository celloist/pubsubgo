package main

import (
	"fmt"
	"sync"
)

type Agent struct {
	mut    sync.Mutex
	subs   map[string][]chan string
	quit   chan struct{}
	closed bool
}

var msgBroker = make(chan string)

func main() { //vrty
	agent := NewAgent()

	sub := agent.Subscribe("test1")

	go agent.Publish("test1", "its a me! message!")

	fmt.Println(<-sub)
	agent.Close()
}

func (agent *Agent) Publish(topic string, msg string) {
	//make it thread safe by aquiring the lock
	agent.mut.Lock()
	defer agent.mut.Unlock()

	if agent.closed {
		return
	}

	for _, ch := range agent.subs[topic] {
		ch <- msg
	}
}

func (agent *Agent) Subscribe(topic string) <-chan string {
	agent.mut.Lock()
	defer agent.mut.Unlock()

	if agent.closed {
		return nil
	}

	ch := make(chan string)
	agent.subs[topic] = append(agent.subs[topic], ch)
	return ch

}

func (agent *Agent) Close() {
	agent.mut.Lock()
	defer agent.mut.Unlock()

	if agent.closed {
		return
	}

	agent.closed = true
	close(agent.quit)

	//close all subscriptions
	for _, ch := range agent.subs {
		for _, sub := range ch {
			close(sub)
		}
	}
}

func NewAgent() *Agent {
	return &Agent{
		subs: make(map[string][]chan string),
		quit: make(chan struct{}),
	}
}
