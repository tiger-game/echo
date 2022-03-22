module github.com/tiger-game/echo

go 1.18

replace github.com/tiger-game/tiger => ../tiger

require (
	github.com/golang/protobuf v1.4.2
	github.com/tiger-game/tiger v0.0.0-20210202032130-591d2a863270
	go.uber.org/atomic v1.7.0
	google.golang.org/protobuf v1.23.0
)

require (
	github.com/google/gops v0.3.14 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	go.uber.org/automaxprocs v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20201207223542-d4d67f95c62d // indirect
)
