package echo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/codec/message"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/xserver"
	"github.com/tiger-game/tiger/xtime"
	"go.uber.org/atomic"
)

var _ xserver.IServer = (*Server)(nil)

type Server struct {
	base   *xserver.Server
	c      chan message.Msg
	logger jlog.Logger
	smap   map[uint64]*channel.ConnChannel
	reqCnt atomic.Uint64
	qps    atomic.Uint64
	tick   *time.Ticker
}

func (s *Server) AfterInit() error {
	return nil
}

func (s *Server) BeforeStop() {

}

func (s *Server) Init(srv *xserver.Server) error {
	s.base = srv
	s.logger = jlog.NewLogByPrefix("Server")
	s.tick = time.NewTicker(time.Second)
	return nil
}

func (s *Server) Run(ctx context.Context, delta xtime.DeltaTimeMsec) {
	select {
	case w := <-s.c:
		if wrap, ok := w.(*msg.WrapMessage); ok {
			wrap.Sender.SendMessage(wrap.Data)
			s.logger.Info("Receive Info:", wrap.Sender.Id(), " Json:", wrap.Data)
			s.qps.Inc()
			s.reqCnt.Inc()
		} else {
			if n, ok := w.(*channel.NotifyNewSession); ok {
				s.smap[n.Id()] = n.ConnChannel
				n.Go()
				s.logger.Info("New Session Id:", n.Id())
			}
		}
	case <-s.tick.C:
		s.logger.Info("Total: ", s.reqCnt.Load(), ",QPS: ", s.qps.Load())
		s.qps.Store(0)
	case <-ctx.Done():
		return
	}
}

func (s *Server) handler() {

}

// AsyncConnectMe run in single goroutine.
func (s *Server) AsyncConnectMe(ctx context.Context, raw net.Conn) error {
	var (
		ch  *channel.ConnChannel
		err error
	)
	if t, ok := raw.(*net.TCPConn); ok && t == nil {
		s.logger.Infof("Raw Nil, %v", t)
		return nil
	}

	// 1.authorizer verify.

	// 2.load data from db and init player's or service's data.
	conf := channel.Config{
		RStreamBufferSize: 1 << 10,
	}
	conf.Init()

	if ch, err = channel.NewChannelWithChan(
		raw,
		s.c,
		serialize.Pack,
		serialize.Unpack,
		channel.Id(serialize.Id()),
		channel.Configure(conf),
	); err != nil {
		return fmt.Errorf("xserver.Server async new session error:%v", err)
	}

	/*
		if sess, err = session.NewRSession(raw,
			serialize.Pack,
			serialize.Unpack,
			session.Id(serialize.Id()),
			session.Configure(conf),
		); err != nil {
			return fmt.Errorf("xserver.Server async new session error:%v", err)
		}

		gom.Go(func() {
			s.gogo(ctx, sess.(session.RSessioner))
		})
	*/

	// TODO: 4.notify router where I am.

	// 5.add session to session manager.
	_ = ch.Signal(&channel.NotifyNewSession{ConnChannel: ch})
	return nil
}

func (s *Server) gogo(ctx context.Context, r *channel.ConnChannel) {
	for {
		select {
		case w := <-r.Receive():
			if wrap, ok := w.(*msg.WrapMessage); ok {
				wrap.Sender.SendMessage(wrap.Data)
				s.logger.Info("Receive Info:", wrap.Sender.Id(), " Json:", wrap.Data)
				s.qps.Inc()
				s.reqCnt.Inc()
			}
		case <-ctx.Done():
			s.logger.Errorf("session(%v) gogo quit", r.Id())
			return
		}
	}
}

func (s *Server) Stop() {
	for _, sess := range s.smap {
		sess.Close()
	}
}

func NewServer() *Server {
	s := &Server{
		c:    make(chan message.Msg, 16),
		smap: make(map[uint64]*channel.ConnChannel),
	}
	return s
}
