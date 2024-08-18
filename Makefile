buf-build-substrait-proto:
	buf generate "https://github.com/substrait-io/substrait.git#tag=v0.54.0" --path proto/substrait

go-run:
	go run -ldflags=-checklinkname=0 main.go