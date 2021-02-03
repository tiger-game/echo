package main

import (
	"context"
	"flag"
	"net"
	"runtime"
	"time"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/gom"
	"github.com/tiger-game/tiger/session"
)

var _Mgr = NewClientMgr()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	n := flag.Int("n", 1, "client number")
	flag.Parse()
	for i := 0; i < *n; i++ {
		c := &Client{}
		c.Connect()
		_Mgr.Add(c)
	}
	gom.Wait()
}

type Client struct {
	cancel context.CancelFunc
	s      session.IRSession
}

func (c *Client) Connect() error {
	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	conn, err := net.Dial("tcp", "127.0.0.1:2233")
	if err != nil {
		return nil
	}
	conf := session.Config{}
	conf.Init()
	if c.s, err = session.NewRSession(conn, serialize.Pack, serialize.Unpack, session.Id(serialize.Id()), session.Configure(conf)); err != nil {
		return err
	}
	c.s.Go()
	gom.Go(func() {
		c.Run(ctx)
	})
	return nil
}

func (c *Client) Run(ctx context.Context) {
	t := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case msg := <-c.s.Receive():
			_ = msg
		case <-t.C:
			if err := c.s.Send(&msg.Echo{Data: "echo 测试，能不能通过？答：能通过就好了"}); err != nil {
				// fmt.Println(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) Close() {
	c.s.Close()
	c.cancel()
}
