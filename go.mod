module github.com/tiger-game/echo

go 1.15

replace github.com/tiger-game/tiger => ../tiger

require (
	github.com/golang/protobuf v1.4.2
	github.com/ip2location/ip2location-go/v9 v9.0.0 // indirect
	github.com/tiger-game/tiger v0.0.0-20210812102611-76a4663fa837
	go.uber.org/atomic v1.7.0
	google.golang.org/protobuf v1.23.0
)
