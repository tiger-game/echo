package msg

import (
	"fmt"

	"github.com/tiger-game/tiger/channel"

	"github.com/tiger-game/echo/pb"
	"github.com/tiger-game/tiger/packet"
	"google.golang.org/protobuf/proto"
)

var (
	_ packet.Msg = (*Echo)(nil)
)

type Echo pb.Echo

func (e *Echo) MsgId() int16 { return 1 }

func (e *Echo) Marshal(pack packet.Packet) error {
	p := (*pb.Echo)(e)
	size := proto.Size(p)
	b := pack.Alloc(size)
	if _, err := (proto.MarshalOptions{}).MarshalAppend(b[:0], p); err != nil {
		return fmt.Errorf("message.Echo Marshal error: %w", err)
	}
	return nil
}

func (e *Echo) Unmarshal(pack packet.Packet) error {
	if err := proto.Unmarshal(pack.Bytes(), (*pb.Echo)(e)); err != nil {
		return fmt.Errorf("message.Echo Unmarshal error: %w", err)
	}
	return nil
}

var _ packet.MessageFactory = (*messageFactory)(nil)

type messageFactory struct {
	objs map[int16]func() packet.Msg
}

func (m *messageFactory) GetMsgById(Id int16) (packet.Msg, error) {
	if fn, ok := m.objs[Id]; ok {
		return fn(), nil
	}
	return nil, fmt.Errorf("error")
}

func NewMsgFactory() packet.MessageFactory {
	m := &messageFactory{objs: map[int16]func() packet.Msg{
		1:  func() packet.Msg { return &Echo{} },
		-1: func() packet.Msg { return &channel.NotifyNewSession{} },
		-2: func() packet.Msg { return &channel.NotifyCloseSession{} },
	}}
	return m
}
