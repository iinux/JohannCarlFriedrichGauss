all:
	go build -ldflags="-s -w" client.go
	go build -ldflags="-s -w" server.go
	go build -ldflags="-s -w" -o libp2p-server ./cmd/libp2p-server
	go build -ldflags="-s -w" -o libp2p-client ./cmd/libp2p-client
	go build -ldflags="-s -w" -o natcheck ./cmd/natcheck
