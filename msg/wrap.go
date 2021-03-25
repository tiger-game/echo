package msg

import (
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/channel/message"
	"github.com/tiger-game/tiger/channel/packet"
)

var _ message.Msg = (*WrapMessage)(nil)

type WrapMessage struct {
	Data message.Msg

	// other information
	Sender channel.DispatchSession
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
