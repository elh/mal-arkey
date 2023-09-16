.PHONY: help
help:
	@echo "You must specify a target: repl, bin"

.PHONY: repl
repl:
	@go run cmd/malarkey.go

.PHONY: bin
bin:
	@mkdir -p bin
	@go build -o bin/malarkey cmd/malarkey.go
	@cp bin/malarkey bin/step0_repl
	@cp bin/malarkey bin/step1_read_print
	@cp bin/malarkey bin/step2_eval
	@cp bin/malarkey bin/step3_env
	@cp bin/malarkey bin/step4_if_fn_do
	@cp bin/malarkey bin/step5_tco
	@cp bin/malarkey bin/step6_file
	@cp bin/malarkey bin/step7_quote
	@cp bin/malarkey bin/step8_macros
	@cp bin/malarkey bin/step9_try
	@cp bin/malarkey bin/stepA_mal

.PHONY: lint
lint:
	@if golint ./... 2>&1 | grep '^'; then exit 1; fi; # Requires comments for exported functions
	@golangci-lint run
