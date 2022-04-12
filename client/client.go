package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tiger-game/echo/pb"
	"github.com/tiger-game/tiger/def"
	"github.com/tiger-game/tiger/dispatch"
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
	msgq chan def.RequestWriter
}

func (c *Client) Type() def.ServTP { return 0 }
func (c *Client) ID() def.ServID   { return 0 }

func (c *Client) Connect(ctx context.Context) error {
	c.msgq = make(chan def.RequestWriter, 4)
	conn, err := net.Dial("tcp", "127.0.0.1:2233")
	if err != nil {
		return err
	}
	conf := io.Config{}
	conf.Init()
	c.s = io.NewWrapIO(
		conn,
		packet.NewDefaultController(msg.NewMsgFactory()),
		dispatch.NewMultiplex(c),
		io.WrapID(serialize.Id()),
		io.Configure(conf))
	c.s.Go(ctx)
	gom.Go(func() {
		c.Run(ctx)
	})
	return nil
}

func (c *Client) Handle(ctx context.Context, req def.RequestWriter) error {
	select {
	case <-ctx.Done():
	case c.msgq <- req:
	}
	return nil
}

func (c *Client) Run(ctx context.Context) {
	t := time.NewTicker(100 * time.Millisecond)
	data := &pb.Echo{
		Data: "echo 测试，能不能通过？答：能通过就好了",
	}

	for {
		select {
		case req := <-c.msgq:
			if _, ok := req.(*io.CloseIO); ok {
				c.s.Close(0)
				return
			} else {
				jlog.Infof("Receive Msg Id:%d, info:%v", req.Msg().MsgID(), req.Msg())
			}
		case <-t.C:
			if err := c.s.WritePacket(packet.NewPacket(packet.Request, data)); err != nil {
				// fmt.Println(err)
			}
		case <-ctx.Done():
			c.s.Close(0)
			return
		}
	}
}
