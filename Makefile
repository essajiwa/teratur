export NOW=$(shell date --rfc-3339=ns)
export PKGS=$(shell go list ./... | grep -v vendor/ | grep -v cmd/)
export TEST_OPTS=-cover -race

.PHONY: dev
dev:
	@echo "${NOW} == RUN DEVELOPMENT ENVIRONMENT (AUTO BUILD)"
	@./.dev/dev.sh

run:
	@echo "${NOW} == RUNNING..."
	@go run cmd/http/main.go

test:
	@echo "${NOW} == TESTING..."
	@go test ${TEST_OPTS} ${PKGS} | tee test.out
	@.dev/test.sh test.out 10

gen:
	@echo "GENERATING MOCK FILES"
	@go generate ./...
	@echo "DONE!"

check:
	@echo "${NOW} == STATICCHECK..."
	@staticcheck ./...
	@echo "STATICCHECK DONE!"

test-all: test check
