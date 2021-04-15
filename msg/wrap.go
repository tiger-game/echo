package msg

import (
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/codec/message"
	"github.com/tiger-game/tiger/codec/packet"
)

var _ message.Msg = (*WrapMessage)(nil)

type WrapMessage struct {
	Data message.Msg

	// other information
	Sender *channel.ConnChannel
}

func (w *WrapMessage) MsgId() int16 {
	panic("implement me")
}

func (w *WrapMessage) Marshal(pack *packet.WPacket) error {
	panic("implement me")
}

func (w *WrapMessage) Unmarshal(pack *packet.RPacket) error {
	panic("implement me")
}
