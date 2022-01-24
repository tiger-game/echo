package echo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/tiger/async"
	"github.com/tiger-game/tiger/packet"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/xserver"
	"github.com/tiger-game/tiger/xtime"
	"go.uber.org/atomic"
)

var (
	_ xserver.IServer = (*Server)(nil)
	_ packet.Handler  = (*Server)(nil)
)

type Server struct {
	base   *xserver.Server
	c      chan async.Message
	logger jlog.Logger
	smap   map[uint64]*async.Stream
	reqCnt atomic.Uint64
	qps    atomic.Uint64
	tick   *time.Ticker
}

func (s *Server) ID() uint64 { return 0 }

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
		if n, ok := w.Msg.(*async.NotifyNewSession); ok {
			s.smap[n.Id()] = n.Stream
			n.Go(ctx)
			s.logger.Info("New Session Id:", n.Id())
		} else {
			w.SendMessage(w)
			//s.logger.Info("Receive Info:", w.MsgId(), " Json:", w.Msg)
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

func (s *Server) Handle(ctx context.Context, msg packet.Msg) error {
	select {
	case s.c = <-msg:
	case <-ctx.Done():
	}
	return nil
}

// AsyncConnectMe run in single goroutine.
func (s *Server) AsyncConnectMe(ctx context.Context, raw net.Conn) error {
	if t, ok := raw.(*net.TCPConn); ok && t == nil {
		s.logger.Infof("Raw Nil, %v", t)
		return nil
	}

	// 1.authorizer verify.

	// 2.load data from db and init player's or service's data.
	conf := async.Config{
		RStreamBufferSize: 1 << 10,
	}
	conf.Init()

	stream, err := async.NewStream(
		raw,
		packet.NewDefaultController(msg.NewMsgFactory(), 0),
		s,
		async.Id(serialize.Id()),
		async.Configure(conf),
	)
	if err != nil {
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
	_ = s.Handle(ctx, &async.NotifyNewSession{Stream: stream})
	return nil
}

func (s *Server) gogo(ctx context.Context, r *async.Stream) {
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
		c:    make(chan async.Message, 16),
		smap: make(map[uint64]*async.Stream),
	}
	return s
}
