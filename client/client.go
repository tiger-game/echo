package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/echo/msg"

	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/gom"
	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/signal"
)

func main() {
	c := &Client{}
	c.Connect()
	gom.Wait()
}

type Client struct {
	cancel context.CancelFunc
	s      session.IRSession
	sig    *signal.SigM
}

func (c *Client) Connect() error {
	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	c.sig = signal.NewSigM()
	c.sig.RegisterSignalAction(signal.SIGINT, c.Close)
	c.sig.Listen()
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
	t := time.NewTicker(time.Second)
	for {
		select {
		case msg := <-c.s.Receive():
			_ = msg
		case <-t.C:
			if err := c.s.Send(&msg.Echo{Data: []byte(`asdasdlasjdlasjdlkasjdkljsdlajsdlkjasdljalsdjsl`)}); err != nil {
				fmt.Println(err)
			}
			fmt.Println("=======")
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) Close() {
	c.s.Close()
	c.cancel()
}
