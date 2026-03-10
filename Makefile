# Makefile

.PHONY: test test-auth test-user sync tidy

sync:
	go work sync

tidy:
	go mod tidy -C apps/gateway
	go mod tidy -C apps/auth
	go mod tidy -C apps/user
	go mod tidy -C apps/order
	go mod tidy -C apps/product
	#---
	go mod tidy -C modules/auth
	go mod tidy -C modules/user
	go mod tidy -C modules/order
	go mod tidy -C modules/product
	
test:
	go test ./...

test-auth:
	cd modules/auth && go test ./...

test-user:
	cd modules/user && go test ./...

test-order:
	cd modules/order && go test ./...

# Thêm rule khác: build, lint, docker, etc.
