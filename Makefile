# ShadowNet Makefile

.PHONY: build-client run-node clean

build-client:
	cd client && go build -o shadownet-client main.go

run-node:
	sudo bash nodes/run_node.sh

clean:
	rm -f client/shadownet-client
	rm -rf ~/.shadownet
