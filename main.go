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

func main() { //vrt
	agent := NewAgent()

	sub := agent.Subscribe("test1")

	go agent.Publish("test1", "its a me! message!")

	fmt.Println(<-sub)
	agent.Close()
}

func (agent *Agent) Publish(topic string, msg string) {
	//make it thread safe by aquiring the lock
	agent.mut.Lock()

	if agent.closed {
		agent.mut.Unlock()
		return
	}

	for _, ch := range agent.subs[topic] {
		ch <- msg
	}
	agent.mut.Unlock()
}

func (agent *Agent) Subscribe(topic string) <-chan string {
	agent.mut.Lock()

	if agent.closed {
		agent.mut.Unlock()
		return nil
	}

	ch := make(chan string)
	agent.subs[topic] = append(agent.subs[topic], ch)
	agent.mut.Unlock()
	return ch

}

func (agent *Agent) Close() {
	agent.mut.Lock()

	if agent.closed {
		defer agent.mut.Unlock()
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
	agent.mut.Unlock()
}

func NewAgent() *Agent {
	return &Agent{
		subs: make(map[string][]chan string),
		quit: make(chan struct{}),
	}
}
