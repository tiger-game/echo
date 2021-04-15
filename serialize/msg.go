package serialize

import (
	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/codec/message"
	"github.com/tiger-game/tiger/codec/packet"
)

func Pack(s *channel.ConnChannel, pack *packet.WPacket, msg message.Msg) error {
	//1.Pack MsgId.
	pack.WriteInt16(msg.MsgId())
	//2.Pack data struct.
	return msg.Marshal(pack)
}

func Unpack(s *channel.ConnChannel, pack *packet.RPacket) (message.Msg, error) {
	//1.Unpack MsgId
	msgId, err := pack.ReadInt16()
	if err != nil {
		return nil, err
	}
	//2.get data struct base msgId
	_ = msgId
	d := &msg.Echo{}
	if err = d.Unmarshal(pack); err != nil {
		return nil, err
	}
	return &msg.WrapMessage{Data: d, Sender: s}, nil
}
