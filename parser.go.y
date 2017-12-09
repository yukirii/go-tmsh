%{
package tmsh

import (
	"fmt"
)

type Expression interface{}

type Token struct {
	token   int
	literal string
}

type Ltm struct {
	restype string
	fqdn    string
	object  Object
}

type Object struct {
	members []Member
}

type Member struct {
	key     string
	value   string
	members []Member
}

%}

%union{
	token Token
	expr  Expression
}

%type<expr> ltm_obj
%type<expr> object
%type<expr> members
%type<expr> value

%token<token> ILLEGAL
%token<token> EOF
%token<token> WS
%token<token> NEWLINE
%token<token> L_BRACE
%token<token> R_BRACE
%token<token> IDENT
%token<token> LTM

%%

ltm_obj
	: LTM IDENT IDENT object
	{
		$$ = $1
	}

object
	: L_BRACE R_BRACE
	{
		$$ = $1
	}
	| L_BRACE NEWLINE R_BRACE
	{
		$$ = $1
	}
	| L_BRACE NEWLINE members R_BRACE
	{
		$$ = $1
	}

members
  : IDENT value NEWLINE
	{
		$$ = $1
	}
	| IDENT object NEWLINE
	{
		$$ = $1
	}
	| members IDENT object NEWLINE
	{
		$$ = $1
	}
	| members IDENT value NEWLINE
	{
		$$ = $1
	}

value
	: IDENT
	{
		$$ = $1
	}
	| value IDENT
	{
		$$ = $1
	}

%%

type Lexer struct {
	s      *Scanner
	result Expression
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

	fmt.Println(tok, lit)

	if tok == EOF {
		return 0
	}

	return tok
}

func (l *Lexer) Error(e string) {
	panic(e)
}

func Decode(data string, node *Node) error {
	l := Lexer{s: NewScanner(data)}
	if yyParse(&l) != 0 {
		panic("Parse error")
	}
	return nil
}
