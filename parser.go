//line parser.go.y:2
package tmsh

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"fmt"
)

const (
	ltmNodeNode = iota
	ltmPoolNode
	ltmVirtualNode
	keyNode
	structNode
	scalarNode
)

type nodeType int

type node struct {
	kind     nodeType
	value    string
	children []*node
}

type Token struct {
	token   int
	literal string
}

//line parser.go.y:32
type yySymType struct {
	yys     int
	token   Token
	ltm     *node
	object  *node
	pair    *node
	members []*node
	value   *node
}

const ILLEGAL = 57346
const EOF = 57347
const WS = 57348
const NEWLINE = 57349
const L_BRACE = 57350
const R_BRACE = 57351
const IDENT = 57352
const LTM = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ILLEGAL",
	"EOF",
	"WS",
	"NEWLINE",
	"L_BRACE",
	"R_BRACE",
	"IDENT",
	"LTM",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line parser.go.y:126

type Lexer struct {
	s      *Scanner
	result *node
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
	panic(fmt.Sprintf("line %d : %s", l.s.line, e))
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 24

var yyAct = [...]int{

	16, 5, 15, 6, 2, 18, 13, 12, 9, 12,
	11, 18, 4, 8, 17, 7, 6, 3, 20, 21,
	19, 14, 10, 1,
}
var yyPact = [...]int{

	-7, -1000, 7, 2, 8, -1000, 6, -1000, -1, -1000,
	-3, -1000, -5, -1000, -1000, -1000, 13, 11, 1, -1000,
	-1000, -1000,
}
var yyPgo = [...]int{

	0, 23, 1, 10, 22, 0,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 4, 4, 3, 3, 3,
	5, 5,
}
var yyR2 = [...]int{

	0, 4, 2, 3, 4, 1, 2, 2, 3, 3,
	2, 1,
}
var yyChk = [...]int{

	-1000, -1, 11, 10, 10, -2, 8, 9, 7, 9,
	-4, -3, 10, 9, -3, 7, -5, -2, 10, 7,
	7, -5,
}
var yyDef = [...]int{

	0, -2, 0, 0, 0, 1, 0, 2, 0, 3,
	0, 5, 0, 4, 6, 7, 0, 0, 11, 8,
	9, 10,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:60
		{
			var kind nodeType
			switch yyDollar[2].token.literal {
			case "node":
				kind = ltmNodeNode
			case "pool":
				kind = ltmPoolNode
			case "virtual":
				kind = ltmVirtualNode
			}
			yylex.(*Lexer).result = &node{
				kind:     kind,
				value:    yyDollar[3].token.literal,
				children: []*node{yyDollar[4].object},
			}
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:79
		{
			yyVAL.object = &node{kind: structNode, value: "", children: []*node{}}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:83
		{
			yyVAL.object = &node{kind: structNode, value: "", children: []*node{}}
		}
	case 4:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line parser.go.y:87
		{
			yyVAL.object = &node{kind: structNode, value: "", children: yyDollar[3].members}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:93
		{
			yyVAL.members = []*node{yyDollar[1].pair}
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:97
		{
			yyVAL.members = append(yyDollar[1].members, yyDollar[2].pair)
		}
	case 7:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:103
		{
			yyVAL.pair = &node{kind: keyNode, value: yyDollar[1].token.literal, children: []*node{}}
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:107
		{
			yyVAL.pair = &node{kind: keyNode, value: yyDollar[1].token.literal, children: []*node{yyDollar[2].value}}
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line parser.go.y:111
		{
			yyVAL.pair = &node{kind: keyNode, value: yyDollar[1].token.literal, children: []*node{yyDollar[2].object}}
		}
	case 10:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line parser.go.y:117
		{
			s := fmt.Sprintf("%s %s", yyDollar[1].token.literal, yyDollar[2].value.value)
			yyVAL.value = &node{kind: scalarNode, value: s, children: []*node{}}
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line parser.go.y:122
		{
			yyVAL.value = &node{kind: scalarNode, value: yyDollar[1].token.literal, children: []*node{}}
		}
	}
	goto yystack /* stack new state and value */
}
