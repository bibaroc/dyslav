GO_ENV ?= CGO_ENABLED=0 

THIS_FILE := $(lastword $(MAKEFILE_LIST))

.PHONY: bin

tools:
	@go get \
		google.golang.org/grpc@v1.32.0 \
		github.com/golang/protobuf/protoc-gen-go@v1.4.3 \
		github.com/bufbuild/buf/cmd/buf@v0.28.0 \
		github.com/google/wire/cmd/wire@v0.4.0
	@go get -d github.com/envoyproxy/protoc-gen-validate@v0.4.1
	@cd $(shell go env GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		&& make build

regen:
	@find . -type f -name '*.pb.*.go' -o -name '*.pb.go' -delete
	@for PROTO in $(shell find . -type f -name '*.proto' | grep -v proto/google/api) ; do \
		echo "Compiling" $${PROTO} ; \
		protoc \
			-I . \
			-I $(shell go env GOPATH)/src \
			-I $(shell go env GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
			--go_out=plugins=grpc,paths=source_relative:./ \
			--validate_out=lang=go,paths=source_relative:. \
			$${PROTO}; \
	done;
	@buf check lint || :;
	@wire gen ./...

bin/%:
	@echo "building cmd/$(shell basename $@)" \
		&& $(GO_ENV) go build \
			-trimpath \
			-gcflags='-e -l' \
			-ldflags='-w -s -extldflags "-static" -X main.version=${VERSION} -X main.commit=${COMMIT}' \
			-o $@ \
			./cmd/$(shell basename $@)
bin:
	@for CMD in $(shell find ./cmd/* -maxdepth 0 -type d -name '*' ) ; do \
		$(MAKE) -B -f $(THIS_FILE) bin/$$(basename $${CMD}); \
	done;
