package main

import (
	"sync"
)

type ClientMgr struct {
	mgr  map[uint64]*Client
	lock sync.Mutex
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

func NewClientMgr() *ClientMgr {
	c := &ClientMgr{
		mgr: make(map[uint64]*Client),
	}
	return c
}
