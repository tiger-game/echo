package echo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/echo/serialize"
	"github.com/tiger-game/tiger/io"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/packet"
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
	c      chan packet.Msg
	logger jlog.Logger
	smap   map[uint64]*io.WrapIO
	reqCnt atomic.Uint64
	qps    atomic.Uint64
	tick   *time.Ticker
}

func (s *Server) ID() uint64 { return 1 }

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

func (s *Server) Loop(ctx context.Context, delta xtime.DeltaTimeMsec) {
	select {
	case msg := <-s.c:
		if n, ok := msg.(*io.NotifyNewStream); ok {
			s.smap[n.ID()] = n.WrapIO
			n.Go(ctx)
			s.logger.Info("New Session Id:", n.ID())
		} else {
			// TODO(mawei):w.SendMessage(w)
			s.logger.Info("Receive Info:", msg.MsgID(), " Json:", msg)
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
	case <-ctx.Done():
	case s.c <- msg:
	}
	return nil
}

// AsyncConnectMe run in single goroutine.
func (s *Server) AsyncConnectMe(ctx context.Context, raw net.Conn) error {
	var (
		w   *io.WrapIO
		err error
	)
	if t, ok := raw.(*net.TCPConn); ok && t == nil {
		s.logger.Infof("Raw Nil, %v", t)
		return nil
	}

	// 1.authorizer verify.

	// 2.load data from db and init player's or service's data.
	conf := io.Config{
		RStreamBufferSize: 1 << 10,
	}
	conf.Init()

	if w, err = io.NewWrapIO(
		raw,
		packet.NewDefaultController(msg.NewMsgFactory()),
		s,
		io.WrapID(serialize.Id()),
		io.Configure(conf),
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
	_ = s.Handle(ctx, &io.NotifyNewStream{WrapIO: w})
	return nil
}

func (s *Server) Stop() {
	for _, sess := range s.smap {
		sess.Close(0)
	}
}

func NewServer() *Server {
	s := &Server{
		c:    make(chan packet.Msg, 16),
		smap: make(map[uint64]*io.WrapIO),
	}
	return s
}
