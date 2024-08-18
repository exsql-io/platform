buf-build-substrait-proto:
	buf generate

go-run:
	go run -ldflags=-checklinkname=0 main.go