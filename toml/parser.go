package toml

type Parser struct {
	lexer    lexer
	previous token
	current  token

	root         Table
	currentTable Table
}
