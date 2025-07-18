// analyzer/memory_optimizer.go
package analyzer

import (
	"sync"
)

// TokenPool reutiliza objetos Token para reducir allocaciones
type TokenPool struct {
	pool sync.Pool
}

// NewTokenPool crea un nuevo pool de tokens
func NewTokenPool() *TokenPool {
	return &TokenPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Token{}
			},
		},
	}
}

// Get obtiene un token del pool
func (tp *TokenPool) Get() *Token {
	return tp.pool.Get().(*Token)
}

// Put devuelve un token al pool
func (tp *TokenPool) Put(token *Token) {
	// Limpiar el token antes de devolverlo al pool
	token.Type = ""
	token.Value = ""
	token.Line = 0
	token.Col = 0
	tp.pool.Put(token)
}

// OptimizedLexer lexer optimizado con pools de memoria
type OptimizedLexer struct {
	tokenPool   *TokenPool
	stringLib   *StringLibrary
	buffer      []rune
	bufferSize  int
}

// NewOptimizedLexer crea un lexer optimizado
func NewOptimizedLexer() *OptimizedLexer {
	return &OptimizedLexer{
		tokenPool:  NewTokenPool(),
		stringLib:  NewStringLibrary(),
		bufferSize: 8192, // Buffer inicial de 8KB
	}
}

// LexOptimized análisis léxico optimizado
func (ol *OptimizedLexer) LexOptimized(input string) []Token {
	// Preallocar slice con capacidad estimada
	estimatedTokens := len(input) / 4 // Estimación heurística
	tokens := make([]Token, 0, estimatedTokens)
	
	// Reutilizar buffer si es posible
	if len(input) > ol.bufferSize {
		ol.buffer = make([]rune, len(input))
		ol.bufferSize = len(input)
	} else if ol.buffer == nil {
		ol.buffer = make([]rune, ol.bufferSize)
	}
	
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

		if isWhitespace(c) {
			col++
			i++
			continue
		}

		// Optimización: usar funciones inline para casos comunes
		if isLetter(c) || c == '_' {
			start := i
			startCol := col
			for i < len(runes) && (isAlphaNumeric(runes[i]) || runes[i] == '_') {
				i++
				col++
			}
			
			word := ol.stringLib.InternString(string(runes[start:i]))
			tokenType := "identifier"
			if keywords[word] {
				tokenType = "keyword"
			}
			
			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: word,
				Line:  line,
				Col:   startCol,
			})
			continue
		}

		if isDigit(c) {
			start := i
			startCol := col
			hasDecimal := false
			
			for i < len(runes) && (isDigit(runes[i]) || (runes[i] == '.' && !hasDecimal)) {
				if runes[i] == '.' {
					hasDecimal = true
				}
				i++
				col++
			}
			
			tokenType := "number"
			if hasDecimal {
				tokenType = "float"
			}
			
			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: ol.stringLib.InternString(string(runes[start:i])),
				Line:  line,
				Col:   startCol,
			})
			continue
		}

		// Procesar operadores de manera optimizada
		if token := ol.processOperator(runes, &i, &col, line); token != nil {
			tokens = append(tokens, *token)
			continue
		}

		// Procesar caracteres especiales
		if token := ol.processSpecialChar(runes, &i, &col, &line); token != nil {
			tokens = append(tokens, *token)
			continue
		}

		// Caracter desconocido
		tokens = append(tokens, Token{
			Type:  "unknown",
			Value: string(c),
			Line:  line,
			Col:   col,
		})
		i++
		col++
	}

	return tokens
}

// processOperator procesa operadores de manera optimizada
func (ol *OptimizedLexer) processOperator(runes []rune, i *int, col *int, line int) *Token {
	c := runes[*i]
	
	switch c {
	case '+', '-', '*', '/', '<', '>', '=', '!', '&', '|':
		return ol.handleComplexOperator(runes, i, col, line, c)
	case ';':
		*i++
		*col++
		return &Token{Type: "semicolon", Value: ";", Line: line, Col: *col - 1}
	case '(':
		*i++
		*col++
		return &Token{Type: "lparen", Value: "(", Line: line, Col: *col - 1}
	case ')':
		*i++
		*col++
		return &Token{Type: "rparen", Value: ")", Line: line, Col: *col - 1}
	case '{':
		*i++
		*col++
		return &Token{Type: "lbrace", Value: "{", Line: line, Col: *col - 1}
	case '}':
		*i++
		*col++
		return &Token{Type: "rbrace", Value: "}", Line: line, Col: *col - 1}
	case ',':
		*i++
		*col++
		return &Token{Type: "comma", Value: ",", Line: line, Col: *col - 1}
	case '.':
		*i++
		*col++
		return &Token{Type: "dot", Value: ".", Line: line, Col: *col - 1}
	}
	
	return nil
}

// handleComplexOperator maneja operadores complejos
func (ol *OptimizedLexer) handleComplexOperator(runes []rune, i *int, col *int, line int, c rune) *Token {
	startCol := *col
	value := string(c)
	*i++
	*col++
	
	// Verificar operadores de dos caracteres
	if *i < len(runes) {
		next := runes[*i]
		switch c {
		case '+':
			if next == '+' || next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '-':
			if next == '-' || next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '*', '/':
			if next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '<', '>':
			if next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '=':
			if next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '!':
			if next == '=' {
				value += string(next)
				*i++
				*col++
			}
		case '&':
			if next == '&' {
				value += string(next)
				*i++
				*col++
			}
		case '|':
			if next == '|' {
				value += string(next)
				*i++
				*col++
			}
		}
	}
	
	return &Token{
		Type:  "operator",
		Value: ol.stringLib.InternString(value),
		Line:  line,
		Col:   startCol,
	}
}

// processSpecialChar procesa caracteres especiales como strings y chars
func (ol *OptimizedLexer) processSpecialChar(runes []rune, i *int, col *int, line *int) *Token {
	c := runes[*i]
	
	switch c {
	case '"':
		return ol.processString(runes, i, col, line)
	case '\'':
		return ol.processChar(runes, i, col, line)
	}
	
	return nil
}

// processString procesa strings con validación
func (ol *OptimizedLexer) processString(runes []rune, i *int, col *int, line *int) *Token {
	start := *i + 1
	startCol := *col
	*i++
	*col++
	
	for *i < len(runes) && runes[*i] != '"' {
		if runes[*i] == '\n' {
			*line++
			*col = 1
		} else {
			*col++
		}
		*i++
	}
	
	if *i >= len(runes) {
		return &Token{Type: "error", Value: "String sin cerrar", Line: *line, Col: startCol}
	}
	
	value := string(runes[start:*i])
	*i++ // cerrar comillas
	*col++
	
	// Validar el contenido del string
	if !ol.stringLib.ValidateJavaString(value) {
		return &Token{Type: "error", Value: "String con caracteres inválidos", Line: *line, Col: startCol}
	}
	
	return &Token{
		Type:  "string",
		Value: ol.stringLib.InternString(value),
		Line:  *line,
		Col:   startCol,
	}
}

// processChar procesa caracteres con validación
func (ol *OptimizedLexer) processChar(runes []rune, i *int, col *int, line *int) *Token {
	start := *i + 1
	startCol := *col
	*i++
	*col++
	
	if *i < len(runes) && runes[*i] != '\'' {
		if *i+1 < len(runes) && runes[*i+1] == '\'' {
			// Char literal válido
			value := string(runes[start : *i+1])
			*i += 2 // saltar char y comilla de cierre
			*col += 2
			return &Token{
				Type:  "char",
				Value: ol.stringLib.InternString(value),
				Line:  *line,
				Col:   startCol,
			}
		} else {
			return &Token{Type: "error", Value: "Char literal inválido", Line: *line, Col: startCol}
		}
	} else {
		return &Token{Type: "error", Value: "Char literal vacío", Line: *line, Col: startCol}
	}
}

// Funciones inline optimizadas
func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c rune) bool {
	return isLetter(c) || isDigit(c)
}

