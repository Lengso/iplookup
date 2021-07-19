GOCMD			:=$(shell which go)
GOBUILD			:=$(GOCMD) build
FLAG			:="-w -s "

BINARY_DIR=bin
BINARY_NAME:=ipgo

# linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -trimpath -o $(BINARY_DIR)/$(BINARY_NAME)-linux

#mac
build-darwin:
	CGO_ENABLED=0 GOOS=darwin $(GOBUILD) -ldflags $(FLAG) -trimpath -o $(BINARY_DIR)/$(BINARY_NAME)-darwin

# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags $(FLAG) -trimpath -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe
# 全平台
build-all:
	make build-linux
	make build-darwin
	make build-win
	upx $(BINARY_DIR)/$(BINARY_NAME)-*
