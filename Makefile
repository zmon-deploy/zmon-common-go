

.PHONY: test
test:
	@echo "\n\033[1;33m+ $@\033[0m"
	go test -short ./...

