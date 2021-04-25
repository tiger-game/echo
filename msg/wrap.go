package msg

import (
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/packet"
)

var _ packet.Msg = (*WrapMessage)(nil)

type WrapMessage struct {
	Data packet.Msg

	// other information
	Sender *channel.NetChan
}

func (w *WrapMessage) MsgId() int16 {
	panic("implement me")
}

func (w *WrapMessage) Marshal(pack packet.Packet) error {
	panic("implement me")
}

func (w *WrapMessage) Unmarshal(pack packet.Packet) error {
	panic("implement me")
}
