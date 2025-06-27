run:
	go run cmd/web/main.go

test:
	go run github.com/bitfield/gotestdox/cmd/gotestdox@latest ./...
