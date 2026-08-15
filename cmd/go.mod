module github.com/ranjith035/morego/cmd

go 1.25.0

replace github.com/ranjith035/morego/proto => ../proto

replace github.com/ranjith035/morego/drivers => ../drivers

replace github.com/ranjith035/morego/ai => ../ai

replace github.com/ranjith035/morego/core => ../core

require google.golang.org/grpc v1.83.0

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
