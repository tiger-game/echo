package msg

import (
	"fmt"

	"github.com/tiger-game/echo/pb"

	"github.com/tiger-game/tiger/session/message"
	"github.com/tiger-game/tiger/session/packet"
	"google.golang.org/protobuf/proto"
)

var (
	_ message.Messager = (*Echo)(nil)
)

type Echo pb.Echo

func (e *Echo) MsgId() int16 { return 1 }

func (e *Echo) Marshal(buffer *packet.ByteStream) error {
	p := (*pb.Echo)(e)
	size := proto.Size(p)
	b := buffer.Alloc(size)
	if _, err := (proto.MarshalOptions{}).MarshalAppend(b[:0], p); err != nil {
		return fmt.Errorf("message.Echo Marshal error: %w", err)
	}
	return nil
}

func (e *Echo) Unmarshal(buffer *packet.ByteStream) error {
	if err := proto.Unmarshal(buffer.Bytes(), (*pb.Echo)(e)); err != nil {
		return fmt.Errorf("message.Echo Unmarshal error: %w", err)
	}
	return nil
}
