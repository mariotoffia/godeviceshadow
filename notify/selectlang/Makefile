.PHONY: generate
generate:
	@echo "Generating language files..."
	@antlr -Dlanguage=Go -no-visitor -package selectlang -o . selectlang.g4
	@rm selectlang.interp selectlang.tokens selectlangLexer.interp selectlangLexer.tokens

