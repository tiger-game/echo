package serialize

import (
	"github.com/tiger-game/echo/msg"
	"github.com/tiger-game/tiger/session"
	"github.com/tiger-game/tiger/session/bytep"
	"github.com/tiger-game/tiger/session/message"
)

func Pack(s session.ISession, byteStreamAnchor *bytep.ByteStream, msg message.IMessage) error {
	//1.Pack MsgId.
	byteStreamAnchor.WriteInt16(msg.MsgId())
	//2.Pack data struct.
	return msg.Marshal(byteStreamAnchor)
}

func Unpack(s session.ISession, byteStream *bytep.ByteStream) (message.IMessage, error) {
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
