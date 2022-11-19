build::
	GOOS=linux GOARCH=amd64 go build -o ./build/handler ./src/handler.go
	zip -j ./build/handler.zip ./build/handler

macos::
	go build -o ./buld/macHandler ./src/handler.go
