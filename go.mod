module github.com/tiger-game/echo

go 1.15

replace github.com/tiger-game/tiger => ../tiger

require (
	github.com/golang/protobuf v1.4.2
	github.com/tiger-game/tiger v0.0.0-20210202032130-591d2a863270
	go.uber.org/atomic v1.7.0
	google.golang.org/protobuf v1.23.0
)
