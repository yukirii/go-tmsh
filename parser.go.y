%{
package tmsh

import (
	"fmt"
)

type Token struct {
	token   int
	literal string
}

type LTMObject struct {
	resType string
	fqdn    string
	object  Object
}

type Object struct {
	members []Member
}

type Member struct {
	key   string
	value interface{}
}

%}

%union{
	token   Token
	ltm     LTMObject
	object  Object
	members []Member
	value   interface{}
}

%type<ltm>     ltm
%type<object>  object
%type<members> members
%type<value>   value

%token<token> ILLEGAL
%token<token> EOF
%token<token> WS
%token<token> NEWLINE
%token<token> L_BRACE
%token<token> R_BRACE
%token<token> IDENT
%token<token> LTM

%%

ltm
	: LTM IDENT IDENT object
	{
		yylex.(*Lexer).result = LTMObject{
			resType: $2.literal,
			fqdn: $3.literal,
			object: $4,
		}
	}

object
	: L_BRACE R_BRACE
	{
		$$ = Object{}
	}
	| L_BRACE NEWLINE R_BRACE
	{
		$$ = Object{}
	}
	| L_BRACE NEWLINE members R_BRACE
	{
		$$ = Object{members: $3}
	}

members
  : IDENT value NEWLINE
	{
		$$ = []Member{Member{key: $1.literal, value: $2}}
	}
	| IDENT object NEWLINE
	{
		$$ = []Member{Member{key: $1.literal, value: $2}}
	}
	| members IDENT object NEWLINE
	{
		m := Member{key: $2.literal, value: $3}
		$$ = append($1, m)
	}
	| members IDENT value NEWLINE
	{
		m := Member{key: $2.literal, value: $3}
		$$ = append($1, m)
	}

value
	: IDENT value
	{
		$$ = fmt.Sprintf("%s %s", $1.literal, $2)
	}
	| IDENT
	{
		$$ = $1.literal
	}

%%

type Lexer struct {
	s      *Scanner
	result LTMObject
}

func (l *Lexer) Lex(lval *yySymType) int {
	var tok int
	var lit string

	for {
		tok, lit = l.s.Scan()
		if tok != WS {
			break
		}
	}

	lval.token = Token{token: tok, literal: lit}

	if tok == EOF {
		return 0
	}

	return tok
}

func (l *Lexer) Error(e string) {
	panic(e)
}
