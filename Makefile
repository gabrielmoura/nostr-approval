ALL:
	@echo "Usage: make [linux-pc|linux-raspberry|windows]"
	@echo "  linux-pc:       Build for Linux PC"
	@echo "  linux-raspberry:Build for Raspberry Pi"
	@echo "  windows:        Build for Windows"

linux-pc:
	@echo "Building for Linux PC"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -trimpath -ldflags="-w -s" -o nostr-approval main.go

linux-raspberry:
	@echo "Building for Raspberry Pi"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -trimpath -ldflags="-w -s" -o nostr-approval main.go

windows:
	@echo "Building for Windows"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -trimpath -ldflags="-w -s" -o nostr-approval.exe main.go

.PHONY: ALL linux-pc linux-raspberry windows