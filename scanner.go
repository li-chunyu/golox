package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Scanner struct {
	src    string
	toks   []*Token
	reader *strings.Reader

	start   int
	current int // 当前未被消耗的
	line    int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:     src,
		toks:    nil,
		reader:  strings.NewReader(src),
		start:   0,
		current: 0,
		line:    0,
	}
}

func (s *Scanner) scanTokens() []*Token {
	for !s.isAtEnd() {
		s.start = s.current
		// fmt.Println(string(s.peek()), s.start, s.current)
		s.scanToken()
	}
	s.toks = append(s.toks, NewToken(EOF, nil, "", s.line))
	return s.toks
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)

	case ')':
		s.addToken(RIGHT_PAREN)

	case '{':
		s.addToken(LEFT_BRACE)

	case '}':
		s.addToken(RIGHT_BRACE)

	case ',':
		s.addToken(COMMA)

	case '.':
		s.addToken(DOT)

	case '-':
		s.addToken(MINUS)

	case '+':
		s.addToken(PLUS)

	case ';':
		s.addToken(SEMICOLON)

	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// 换行符不会被消耗，在下一次循环中，进入 case '\n'，更新 s.line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line += 1
	case '"':
		s.string()
	default:
		// parse number
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			// 假设所有以字母和下划线开头的 lexeme 是一个 identifier
			s.identifier()
		} else {
			perror(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) addToken(t TokenType) {
	s.addTokenLiteral(t, nil)
}

func (s *Scanner) addTokenLiteral(t TokenType, literal interface{}) {
	lex := s.src[s.start:s.current]
	s.toks = append(s.toks, NewToken(t, literal, lex, s.line))
}

// advance lookahead and consume.
func (s *Scanner) advance() rune {
	c, sz, err := s.reader.ReadRune()
	if err != nil {
		perror(s.line, fmt.Sprintf("%v", err))
		return 0
	}
	s.current += sz
	return c
}

// match is like a conditional advance(). We only consume the current character if it’s what we’re looking for.
func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	c, sz, _ := s.reader.ReadRune()
	if c != expected {
		if err := s.reader.UnreadRune(); err != nil {
			panic(fmt.Sprintf("Error match UnreadRune, %v", err))
		} else {
			return false
		}
	}
	s.current += sz
	return true
}

// peek is sort of like advance(), but doesn’t consume the character.
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	c, _, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("Error peak. %v", err))
	}
	err = s.reader.UnreadRune()
	if err != nil {
		panic(fmt.Sprintf("Error peak. %v", err))
	}
	return c
}

func (s *Scanner) peekNext() rune {
	var (
		err error
		c   rune
		i   int
	)

	for ; i < 2 && s.reader.Len() != 0; i++ {
		c, _, err = s.reader.ReadRune()
	}

	if err != nil {
		panic(errorMsg(s.line, "peekNext", err.Error()))
	}

	if _, err := s.reader.Seek(int64(s.current), io.SeekStart); err != nil {
		panic(errorMsg(s.line, "peekNext", err.Error()))
	}

	if i < 2 {
		return 0
	}

	return c
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		panic(errorMsg(s.line, "func string", "unterminated string."))
	}
	// now current is index of ".
	s.advance()
	lit := s.src[s.start+1 : s.current-1]
	s.addTokenLiteral(STRING, lit)
}

func (s *Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) && !s.isAtEnd() {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) && !s.isAtEnd() {
			s.advance()
		}
	}

	a := s.src[s.start:s.current]
	num, err := strconv.ParseFloat(a, 64)
	if err != nil {
		panic(errorMsg(s.line, "parse number", err.Error()))
	}
	s.addTokenLiteral(NUMBER, num)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func (s *Scanner) isAlphaDigit(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) identifier() {
	for !s.isAtEnd() && s.isAlphaDigit(s.peek()) {
		s.advance()
	}

	text := s.src[s.start:s.current]
	if t, ok := keywords[text]; ok {
		s.addToken(t)
	} else {
		s.addToken(IDENTIFIER)
	}
}
