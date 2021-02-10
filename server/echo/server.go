package echo

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/tiger-game/tiger/xtime"

	"github.com/tiger-game/echo/serialize"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/tiger/jlog"
	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/session/message"
	"github.com/tiger-game/tiger/xserver"
)

var _ xserver.IServer = (*Server)(nil)

type Server struct {
	base   *xserver.Server
	c      chan message.IMessage
	logger jlog.ILog
	smap   map[uint64]session.ISession
	reqCnt uint64
	qps    uint64
	tick   *time.Ticker
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
			wrap.Sender.Send(wrap.Data)
			s.logger.Info("Receive Info:", wrap.Sender.Id(), " Json:", wrap.Data)
			s.qps++
			s.reqCnt++
		} else {
			if new, ok := w.(*session.NotifyNewSession); ok {
				s.smap[new.Id()] = new.ISession
				new.ISession.Go()
				s.logger.Info("New Session Id:", new.Id())
			}
		}
	case <-s.tick.C:
		s.logger.Info("Total: ", s.reqCnt, ",QPS: ", s.qps)
		s.qps = 0
	case <-ctx.Done():
		return
	}
}

func (s *Server) handler() {

}

func (s *Server) AsyncConnectMe(raw net.Conn) error {
	var (
		sess session.ISession
		err  error
	)
	if t, ok := raw.(*net.TCPConn); ok && t == nil {
		s.logger.Infof("Raw Nil, %v", t)
		return nil
	}

	s.logger.Infof("Connect Me...%T", raw)
	// 1.authorizer verify.

	// 2.load data from db and init player's or service's data.
	conf := session.Config{
		RStreamBufferSize: 1 << 10,
	}
	conf.Init()
	if sess, err = session.NewSession(
		raw,
		serialize.Pack,
		serialize.Unpack,
		session.Id(serialize.Id()),
		session.Output(s.c),
		session.Configure(conf),
	); err != nil {
		return fmt.Errorf("xserver.Server async new session error:%v", err)
	}

	// TODO: 4.notify router where I am.

	// 5.add session to session manager.
	sess.AsyncNotify(&session.NotifyNewSession{ISession: sess})
	return nil
}

func (s *Server) Stop() {
	for _, sess := range s.smap {
		sess.Close()
	}
}

func NewServer() *Server {
	s := &Server{
		c:    make(chan message.IMessage, 16),
		smap: make(map[uint64]session.ISession),
	}
	return s
}
