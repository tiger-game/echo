package msg

import (
	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/session/bytep"
	"github.com/tiger-game/tiger/session/message"
)

var _ message.IMessage = (*WrapMessage)(nil)

type WrapMessage struct {
	Data message.IMessage

	// other information
	Sender session.ISession
}

func (w *WrapMessage) MsgId() int16 {
	panic("implement me")
}

func (w *WrapMessage) Marshal(buffer *bytep.ByteStream) error {
	panic("implement me")
}

func (w *WrapMessage) Unmarshal(buffer *bytep.ByteStream) error {
	panic("implement me")
}
