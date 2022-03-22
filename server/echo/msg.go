package echo

import (
	"context"
	"github.com/tiger-game/tiger/def"
	"github.com/tiger-game/tiger/io"
)

var _ def.RequestWriter = (*NotifyNewIO)(nil)

type NotifyNewIO struct {
	io.OuterIO
}

func (n *NotifyNewIO) PeerServType() def.ServTP                        { return 0 }
func (n *NotifyNewIO) PeerServID() def.ServID                          { return 0 }
func (n *NotifyNewIO) ReqID() int32                                    { return 0 }
func (n *NotifyNewIO) Msg() def.Msg                                    { return nil }
func (n *NotifyNewIO) Response(ctx context.Context, msg def.Msg) error { return nil }
