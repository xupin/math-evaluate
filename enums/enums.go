package enums

const (
	ILLEGAL = iota
	EOF

	// symbols
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	PERIOD // .

	// literals
	NUMBER
	IDENT

	// operators
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	XOR // ^
)
