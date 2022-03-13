package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tiger-game/tiger/io"
	"net"
	"runtime"
	"time"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/gom"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/packet"
	"github.com/tiger-game/tiger/signal"
)

var _Mgr = NewClientMgr()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	n := flag.Int("n", 1, "client number")
	flag.Parse()
	jlog.GLogInit(jlog.LogLevel(jlog.ERROR), jlog.LogDir("./log"))
	ctx, stop := signal.Signal()
	defer stop()

	for i := 0; i < *n; i++ {
		c := &Client{}
		if err := c.Connect(ctx); err != nil {
			fmt.Println(err)
			continue
		}
		_Mgr.Add(c)
	}
	gom.Wait()
}

type Client struct {
	s    *io.WrapIO
	msgq chan packet.Msg
}

func (c *Client) ID() uint64 { return 2 }

func (c *Client) Connect(ctx context.Context) error {
	c.msgq = make(chan packet.Msg, 4)
	conn, err := net.Dial("tcp", "127.0.0.1:2233")
	if err != nil {
		return err
	}
	conf := io.Config{}
	conf.Init()
	if c.s, err = io.NewWrapIO(conn, packet.NewDefaultController(msg.NewMsgFactory()), c, io.WrapID(serialize.Id()), io.Configure(conf)); err != nil {
		return err
	}
	c.s.Go(ctx)
	gom.Go(func() {
		c.Run(ctx)
	})
	return nil
}

func (c *Client) Handle(ctx context.Context, msg packet.Msg) error {
	select {
	case <-ctx.Done():
	case c.msgq <- msg:
	}
	return nil
}

func (c *Client) Run(ctx context.Context) {
	t := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case msg := <-c.msgq:
			_ = msg
		case <-t.C:
			if err := c.s.SendMessage(&msg.Echo{Data: "echo 测试，能不能通过？答：能通过就好了"}); err != nil {
				// fmt.Println(err)
			}
		case <-ctx.Done():
			c.s.Close(0)
			return
		}
	}
}
