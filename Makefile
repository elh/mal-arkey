.PHONY: mal
mal:
	@mkdir -p bin
	@go build -o bin/mal cmd/mal.go
	@cp bin/mal bin/step0_repl
	@cp bin/mal bin/step1_read_print
	@cp bin/mal bin/step2_eval
	@cp bin/mal bin/step3_env
	@cp bin/mal bin/step4_if_fn_do
	@cp bin/mal bin/step5_tco
	@cp bin/mal bin/step6_file
	@cp bin/mal bin/step7_quote
	@cp bin/mal bin/step8_macros
	@cp bin/mal bin/step9_try
	@cp bin/mal bin/stepA_mal