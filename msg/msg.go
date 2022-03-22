package msg

import (
	"fmt"
	"github.com/tiger-game/tiger/def"

	"github.com/tiger-game/echo/pb"
)

var (
	_ def.Msg = (*pb.Echo)(nil)
)

var _ def.MessageFactory = (*messageFactory)(nil)

type messageFactory struct {
	objs map[int16]func() def.Msg
}

func (m *messageFactory) NewMsgByID(Id int16) (def.Msg, error) {
	if fn, ok := m.objs[Id]; ok {
		return fn(), nil
	}
	return nil, fmt.Errorf("error")
}

func NewMsgFactory() def.MessageFactory {
	m := &messageFactory{objs: map[int16]func() def.Msg{
		1: func() def.Msg { return &pb.Echo{} },
	}}
	return m
}
