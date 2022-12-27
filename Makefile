GOPATHS=`go env GOPATH`
init:
	@echo "init mqant tools"

bin:
	#rm protoc-gen-hip
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -o "protoc-gen-hip" main.go

example:
	 @protoc  --plugin=./protoc-gen-hip --go_out=./ --hip_out=./ examples/greeter.proto

