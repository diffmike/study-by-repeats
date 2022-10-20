build::
	GOOS=linux GOARCH=amd64 go build -o ./handler/handler ./handler/handler.go
	zip -j ./handler/handler.zip ./handler/handler

macos::
	go build -o ./handler/macHandler ./handler/handler.go
