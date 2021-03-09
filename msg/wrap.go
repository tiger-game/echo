package msg

import (
	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/session/message"
	"github.com/tiger-game/tiger/session/packet"
)

var _ message.Messager = (*WrapMessage)(nil)

type WrapMessage struct {
	Data message.Messager

	// other information
	Sender session.Sessioner
}

func (w *WrapMessage) MsgId() int16 {
	panic("implement me")
}

func (w *WrapMessage) Marshal(buffer *packet.ByteStream) error {
	panic("implement me")
}

func (w *WrapMessage) Unmarshal(buffer *packet.ByteStream) error {
	panic("implement me")
}
