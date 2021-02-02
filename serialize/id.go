package serialize

import "go.uber.org/atomic"

var id atomic.Uint64

func Id() uint64 { return id.Inc() }
