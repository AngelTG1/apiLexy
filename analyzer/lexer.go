package analyzer

import (
	"unicode"
)

func Lex(input string) []Token {
	var tokens []Token
	runes := []rune(input)
	i := 0
	line := 1
	col := 1

	for i < len(runes) {
		c := runes[i]

		if c == '\n' {
			line++
			col = 1
			i++
			continue
		}

		if unicode.IsSpace(c) {
			col++
			i++
			continue
		}

		if unicode.IsLetter(c) || c == '_' {
			start := i
			startCol := col
			for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
				col++
			}
			word := string(runes[start:i])
			tokenType := "identifier"
			if keywords[word] {
				tokenType = "keyword"
			}
			tokens = append(tokens, Token{Type: tokenType, Value: word, Line: line, Col: startCol})
			continue
		}

		if unicode.IsDigit(c) {
			start := i
			startCol := col
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
				col++
			}
			tokens = append(tokens, Token{Type: "number", Value: string(runes[start:i]), Line: line, Col: startCol})
			continue
		}

		switch c {
		case '+':
			if i+1 < len(runes) && runes[i+1] == '+' {
				// Verificar si hay un tercer + (c+++)
				if i+2 < len(runes) && runes[i+2] == '+' {
					// Tokenizar como ++ y luego +
					tokens = append(tokens, Token{Type: "operator", Value: "++", Line: line, Col: col})
					tokens = append(tokens, Token{Type: "operator", Value: "+", Line: line, Col: col + 2})
					i += 3
					col += 3
				} else {
					tokens = append(tokens, Token{Type: "operator", Value: "++", Line: line, Col: col})
					i += 2
					col += 2
				}
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "+=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "+", Line: line, Col: col})
				i++
				col++
			}
		case '-':
			if i+1 < len(runes) && runes[i+1] == '-' {
				// Verificar si hay un tercer - (c---)
				if i+2 < len(runes) && runes[i+2] == '-' {
					// Tokenizar como -- y luego -
					tokens = append(tokens, Token{Type: "operator", Value: "--", Line: line, Col: col})
					tokens = append(tokens, Token{Type: "operator", Value: "-", Line: line, Col: col + 2})
					i += 3
					col += 3
				} else {
					tokens = append(tokens, Token{Type: "operator", Value: "--", Line: line, Col: col})
					i += 2
					col += 2
				}
			} else if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "-=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "-", Line: line, Col: col})
				i++
				col++
			}
		case '*':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "*=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "*", Line: line, Col: col})
				i++
				col++
			}
		case '/':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "/=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "/", Line: line, Col: col})
				i++
				col++
			}
		case '<':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "<=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "<", Line: line, Col: col})
				i++
				col++
			}
		case '>':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: ">=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: ">", Line: line, Col: col})
				i++
				col++
			}
		case '=':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "==", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "=", Line: line, Col: col})
				i++
				col++
			}
		case '!':
			if i+1 < len(runes) && runes[i+1] == '=' {
				tokens = append(tokens, Token{Type: "operator", Value: "!=", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "!", Line: line, Col: col})
				i++
				col++
			}
		case '&':
			if i+1 < len(runes) && runes[i+1] == '&' {
				tokens = append(tokens, Token{Type: "operator", Value: "&&", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "&", Line: line, Col: col})
				i++
				col++
			}
		case '|':
			if i+1 < len(runes) && runes[i+1] == '|' {
				tokens = append(tokens, Token{Type: "operator", Value: "||", Line: line, Col: col})
				i += 2
				col += 2
			} else {
				tokens = append(tokens, Token{Type: "operator", Value: "|", Line: line, Col: col})
				i++
				col++
			}
		case ';':
			tokens = append(tokens, Token{Type: "semicolon", Value: ";", Line: line, Col: col})
			i++
			col++
		case '(':
			tokens = append(tokens, Token{Type: "lparen", Value: "(", Line: line, Col: col})
			i++
			col++
		case ')':
			tokens = append(tokens, Token{Type: "rparen", Value: ")", Line: line, Col: col})
			i++
			col++
		case '{':
			tokens = append(tokens, Token{Type: "lbrace", Value: "{", Line: line, Col: col})
			i++
			col++
		case '}':
			tokens = append(tokens, Token{Type: "rbrace", Value: "}", Line: line, Col: col})
			i++
			col++
		case ',':
			tokens = append(tokens, Token{Type: "comma", Value: ",", Line: line, Col: col})
			i++
			col++
		case '.':
			tokens = append(tokens, Token{Type: "dot", Value: ".", Line: line, Col: col})
			i++
			col++
		case '"':
			start := i + 1
			startCol := col
			i++
			col++
			for i < len(runes) && runes[i] != '"' {
				if runes[i] == '\n' {
					line++
					col = 1
				} else {
					col++
				}
				i++
			}
			if i >= len(runes) {
				tokens = append(tokens, Token{Type: "error", Value: "String sin cerrar", Line: line, Col: startCol})
			} else {
				value := string(runes[start:i])
				tokens = append(tokens, Token{Type: "string", Value: value, Line: line, Col: startCol})
				i++ // cerrar comillas
				col++
			}
		case '\'':
			start := i + 1
			startCol := col
			i++
			col++
			if i < len(runes) && runes[i] != '\'' {
				if i+1 < len(runes) && runes[i+1] == '\'' {
					// Char literal válido
					value := string(runes[start : i+1])
					tokens = append(tokens, Token{Type: "char", Value: value, Line: line, Col: startCol})
					i += 2 // saltar char y comilla de cierre
					col += 2
				} else {
					tokens = append(tokens, Token{Type: "error", Value: "Char literal inválido", Line: line, Col: startCol})
					i++
					col++
				}
			} else {
				tokens = append(tokens, Token{Type: "error", Value: "Char literal vacío", Line: line, Col: startCol})
				i++
				col++
			}
		default:
			tokens = append(tokens, Token{Type: "unknown", Value: string(c), Line: line, Col: col})
			i++
			col++
		}
	}

	return tokens
}
