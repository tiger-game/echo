package msg

import (
	"fmt"
	"github.com/tiger-game/tiger/jbuff"

	"github.com/tiger-game/echo/pb"
	"github.com/tiger-game/tiger/packet"
	"google.golang.org/protobuf/proto"
)

var (
	_ packet.Msg = (*Echo)(nil)
)

type Echo pb.Echo

func (e *Echo) MsgID() int16 { return 1 }

func (e *Echo) Marshal(s *jbuff.StreamLoc) error {
	p := (*pb.Echo)(e)
	size := proto.Size(p)
	b := s.Alloc(size)
	if _, err := (proto.MarshalOptions{}).MarshalAppend(b[:0], p); err != nil {
		return fmt.Errorf("message.Echo Marshal error: %w", err)
	}
	return nil
}

func (e *Echo) Unmarshal(s *jbuff.Stream) error {
	if err := proto.Unmarshal(s.Bytes(), (*pb.Echo)(e)); err != nil {
		return fmt.Errorf("message.Echo Unmarshal error: %w", err)
	}
	return nil
}

var _ packet.MessageFactory = (*messageFactory)(nil)

type messageFactory struct {
	objs map[int16]func() packet.Msg
}

func (m *messageFactory) NewMsgByID(Id int16) (packet.Msg, error) {
	if fn, ok := m.objs[Id]; ok {
		return fn(), nil
	}
	return nil, fmt.Errorf("error")
}

func NewMsgFactory() packet.MessageFactory {
	m := &messageFactory{objs: map[int16]func() packet.Msg{
		1: func() packet.Msg { return &Echo{} },
		// -1: func() packet.Msg { return &aysnc.NotifyNewSession{} },
		// -2: func() packet.Msg { return &async.NotifyCloseSession{} },
	}}
	return m
}
