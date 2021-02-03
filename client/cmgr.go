package main

import (
	"sync"

	"github.com/tiger-game/tiger/signal"
)

type ClientMgr struct {
	mgr  map[uint64]*Client
	lock sync.Mutex
	sig  *signal.SigM
}

func (c *ClientMgr) Add(client *Client) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.mgr[client.s.Id()] = client
}

func (c *ClientMgr) Del(id uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.mgr, id)
}

func (c *ClientMgr) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, client := range c.mgr {
		client.Close()
	}
}

func NewClientMgr() *ClientMgr {
	c := &ClientMgr{
		mgr: make(map[uint64]*Client),
		sig: signal.NewSigM(),
	}
	c.sig.RegisterSignalAction(signal.SIGINT, c.Close)
	c.sig.Listen()
	return c
}
