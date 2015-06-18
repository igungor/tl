/*
Package	tl implements Telegram's TL(Type Language) lexer and parser. It takes a
[]byte as source which can be tokenized through the scanner and then they fed
to the parser. The output is an abstract syntax tree (AST) representing the TL
program.
*/
package tl
