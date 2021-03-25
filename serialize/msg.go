package serialize

import (
	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/tiger/channel"
	"github.com/tiger-game/tiger/channel/message"
	"github.com/tiger-game/tiger/channel/packet"
)

func Pack(s channel.DispatchSession, byteStreamAnchor *packet.ByteStream, msg message.Msg) error {
	//1.Pack MsgId.
	byteStreamAnchor.WriteInt16(msg.MsgId())
	//2.Pack data struct.
	return msg.Marshal(byteStreamAnchor)
}

func Unpack(s channel.DispatchSession, byteStream *packet.ByteStream) (message.Msg, error) {
	//1.Unpack MsgId
	msgId, err := byteStream.ReadInt16()
	if err != nil {
		return nil, err
	}
	//2.get data struct base msgId
	_ = msgId
	d := &msg.Echo{}
	if err = d.Unmarshal(byteStream); err != nil {
		return nil, err
	}
	return &msg.WrapMessage{Data: d, Sender: s}, nil
}
