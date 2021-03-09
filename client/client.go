package main

import (
	"context"
	"flag"
	"net"
	"runtime"
	"time"

	"github.com/tiger-game/tiger/jlog"

	"github.com/tiger-game/tiger/signal"

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
	jlog.GLogInit(jlog.LogLevel(jlog.ERROR), jlog.LogDir("./log"))
	ctx, stop := signal.Monitor()
	defer stop()

	for i := 0; i < *n; i++ {
		c := &Client{}
		c.Connect(ctx)
		_Mgr.Add(c)
	}
	gom.Wait()
}

type Client struct {
	s session.RSessioner
}

func (c *Client) Connect(ctx context.Context) error {
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
			c.s.Close()
			return
		}
	}
}
