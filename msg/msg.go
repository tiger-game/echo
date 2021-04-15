package msg

import (
	"fmt"

	"github.com/tiger-game/tiger/codec/message"
	"github.com/tiger-game/tiger/codec/packet"

	"github.com/tiger-game/echo/pb"
	"google.golang.org/protobuf/proto"
)

var (
	_ message.Msg = (*Echo)(nil)
)

type Echo pb.Echo

func (e *Echo) MsgId() int16 { return 1 }

func (e *Echo) Marshal(pack *packet.WPacket) error {
	p := (*pb.Echo)(e)
	size := proto.Size(p)
	b := pack.Alloc(size)
	if _, err := (proto.MarshalOptions{}).MarshalAppend(b[:0], p); err != nil {
		return fmt.Errorf("message.Echo Marshal error: %w", err)
	}
	return nil
}

func (e *Echo) Unmarshal(pack *packet.RPacket) error {
	if err := proto.Unmarshal(pack.Bytes(), (*pb.Echo)(e)); err != nil {
		return fmt.Errorf("message.Echo Unmarshal error: %w", err)
	}
	return nil
}
