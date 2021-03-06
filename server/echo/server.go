package echo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/tiger/packet"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/xserver"
	"github.com/tiger-game/tiger/xtime"
	"go.uber.org/atomic"
)

var _ xserver.IServer = (*Server)(nil)

type Server struct {
	base   *xserver.Server
	c      chan channel.Message
	logger jlog.Logger
	smap   map[uint64]*channel.NetChan
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
		if n, ok := w.Msg.(*channel.NotifyNewSession); ok {
			s.smap[n.Id()] = n.NetChan
			n.Go()
			s.logger.Info("New Session Id:", n.Id())
		} else {
			w.SendMessage(w)
			s.logger.Info("Receive Info:", w.MsgId(), " Json:", w.Msg)
			s.qps.Inc()
			s.reqCnt.Inc()
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
		ch  *channel.NetChan
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
		packet.NewDefaultController(msg.NewMsgFactory(), 0),
		s.c,
		channel.Id(serialize.Id()),
		channel.Configure(conf),
	); err != nil {
		return fmt.Errorf("xserver.Server async new session error:%v", err)
	}

	/*
		if sess, err = session.NewRSession(raw,
			packet.NewDefaultController(msg.NewMsgFactory(), 0),
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
	_ = ch.NotifyApp(&channel.NotifyNewSession{NetChan: ch})
	return nil
}

func (s *Server) gogo(ctx context.Context, r *channel.NetChan) {
	for {
		select {
		case w := <-r.ReceiveMessage():
			w.SendMessage(w.Msg)
			s.logger.Info("Receive Info:", w.MsgId(), " Json:", w.Msg)
			s.qps.Inc()
			s.reqCnt.Inc()
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
		c:    make(chan channel.Message, 16),
		smap: make(map[uint64]*channel.NetChan),
	}
	return s
}
