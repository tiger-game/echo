package echo

import (
	"context"
	"fmt"
	"net"

	"github.com/tiger-game/echo/serialize"

	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/tiger/jlog"

	"github.com/tiger-game/tiger/session/message"

	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/xserver"
)

var _ xserver.IServer = (*Server)(nil)

type Server struct {
	base   *xserver.Server
	c      chan message.IMessage
	logger jlog.ILog
	smap   map[uint64]session.ISession
}

func (s *Server) Init(srv *xserver.Server) error {
	s.base = srv
	s.logger = jlog.NewLogByPrefix("Server")
	return nil
}

func (s *Server) Run(ctx context.Context, delta int64) {
	select {
	case w := <-s.c:
		s.logger.Infof("Receive Msg: %T\n", w)
		if wrap, ok := w.(*msg.WrapMessage); ok {
			wrap.Sender.Send(wrap.Data)
			s.logger.Info("Receive Info:", wrap.Sender.Id(), " Json:", wrap.Data)
		} else {
			if new, ok := w.(*session.NotifyNewSession); ok {
				s.smap[new.Id()] = new.ISession
				new.ISession.Go()
			}
			// internal message
		}
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

	s.logger.Infof("Connect Me...")
	// 1.authorizer verify.

	// 2.load data from db and init player's or service's data.
	conf := session.Config{}
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
