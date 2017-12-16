package tmsh

import (
	"bytes"
	"strings"
)

type Scanner struct {
	r    *strings.Reader
	line int
}

func NewScanner(data string) *Scanner {
	return &Scanner{r: strings.NewReader(data)}
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
		ch == '.' || ch == '_' || ch == '-' || ch == ':' ||
		ch == '/' || ch == '(' || ch == ')'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return rune(0)
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan() (tok int, lit string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.unread()
		return s.scanIdent()
	}

	switch ch {
	case rune(0):
		return EOF, ""
	case '\n':
		s.line++
		return NEWLINE, string(ch)
	case '{':
		return L_BRACE, string(ch)
	case '}':
		return R_BRACE, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok int, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == rune(0) {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok int, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == rune(0) {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch buf.String() {
	case "ltm":
		return LTM, buf.String()
	}

	return IDENT, buf.String()
}
