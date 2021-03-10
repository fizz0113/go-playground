
PROTOC=protoc

all: clean proto

install-protoc:
ifeq ($(shell which protoc),)
	@echo " * install protoc"
	@curl -sOL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-$(shell uname -s)-$(shell uname -p).zip
	@unzip protoc-3.7.1-$(shell uname -s)-$(shell uname -p).zip -d /usr/local 'bin/protoc'
	@unzip protoc-3.7.1-$(shell uname -s)-$(shell uname -p).zip -d /usr/local 'include/*'
else ifeq ($(shell protoc --version | grep 3.7.1),)
	@echo " * remove protoc $(shell protoc --version)"
	@rm -f $(shell which protoc)
	@make --no-print-directory install-protoc
else
	@echo " * skip install protoc"
endif

install-go-tools:
	@echo " * install protoc-gen-xxx"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0

proto: install-protoc install-go-tools
	$(PROTOC) \
		--proto_path=. \
		--proto_path=/usr/local/include \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(shell find ./* -iname *.proto)

clean:
	@rm -f ./*.go
