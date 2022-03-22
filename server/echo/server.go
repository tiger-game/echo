package echo

import (
	"context"
	"fmt"
	"github.com/tiger-game/tiger/def"
	"github.com/tiger-game/tiger/dispatch"
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
	_ def.Handler     = (*Server)(nil)
)

type Server struct {
	base   *xserver.Server
	c      chan def.RequestWriter
	logger jlog.Logger
	smap   map[uint64]io.OuterIO
	reqCnt atomic.Uint64
	qps    atomic.Uint64
	tick   *time.Ticker
}

func (s *Server) Type() def.ServTP { return def.Server }
func (s *Server) ID() def.ServID   { return def.ServID(1) }

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
	case req := <-s.c:
		if n, ok := req.(*NotifyNewIO); ok {
			s.smap[n.ID()] = n.OuterIO
			s.logger.Info("New Session Id:", n.ID())
		} else if cl, ok := req.(*io.CloseIO); ok {
			delete(s.smap, cl.Raw().ID())
		} else {
			req.Response(ctx, req.Msg())
			s.logger.Info("Receive Info:", req.Msg().MsgID(), " Json:", req.Msg())
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

func (s *Server) Handle(ctx context.Context, req def.RequestWriter) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("server.Handle cancel deal request:%v", req)
	case s.c <- req:
	}
	return nil
}

// AsyncConnectMe run in single goroutine.
func (s *Server) AsyncConnectMe(ctx context.Context, raw net.Conn) error {
	var (
		w *io.WrapIO
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

	w = io.NewWrapIO(
		raw,
		packet.NewDefaultController(msg.NewMsgFactory()),
		dispatch.NewMultiplex(s),
		io.WrapID(serialize.Id()),
		io.Configure(conf),
	)

	// TODO: 4.notify router where I am.

	w.Go(ctx)
	// 5.add session to session manager.
	_ = s.Handle(ctx, &NotifyNewIO{OuterIO: w})
	return nil
}

func (s *Server) Stop() {
	for _, sess := range s.smap {
		sess.Close(0)
	}
}

func NewServer() *Server {
	s := &Server{
		c:    make(chan def.RequestWriter, 16),
		smap: make(map[uint64]io.OuterIO, 8),
	}
	return s
}
