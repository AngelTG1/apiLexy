// analyzer/parser.go
package analyzer

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
	errors []string
}

func Parse(tokens []Token) (bool, []string) {
	parser := &Parser{tokens: tokens, pos: 0, errors: []string{}}
	parser.validateBrackets()
	parser.validateStatements()
	parser.validateStrings()
	parser.validateSemicolons()
	return len(parser.errors) == 0, parser.errors
}

func (p *Parser) validateBrackets() {
	parenStack := []Token{}
	braceStack := []Token{}

	for _, token := range p.tokens {
		switch token.Type {
		case "lparen":
			parenStack = append(parenStack, token)
		case "rparen":
			if len(parenStack) == 0 {
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: ')' sin '(' correspondiente", token.Line))
			} else {
				parenStack = parenStack[:len(parenStack)-1]
			}
		case "lbrace":
			braceStack = append(braceStack, token)
		case "rbrace":
			if len(braceStack) == 0 {
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: '}' sin '{' correspondiente", token.Line))
			} else {
				braceStack = braceStack[:len(braceStack)-1]
			}
		}
	}

	for _, token := range parenStack {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: '(' sin cerrar", token.Line))
	}

	for _, token := range braceStack {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: '{' sin cerrar", token.Line))
	}
}

// Nueva función para validar strings mal formados
func (p *Parser) validateStrings() {
	for _, token := range p.tokens {
		if token.Type == "error" {
			if strings.Contains(token.Value, "String sin cerrar") {
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: String sin cerrar - falta comilla de cierre", token.Line))
			}
		}
	}
}

// Nueva función para validar punto y coma faltante
func (p *Parser) validateSemicolons() {
	for i := 0; i < len(p.tokens); i++ {
		// Detectar System.out.println o System.out.print
		if p.tokens[i].Value == "System" && i+4 < len(p.tokens) &&
		   p.tokens[i+1].Value == "." && p.tokens[i+2].Value == "out" &&
		   p.tokens[i+3].Value == "." && 
		   (p.tokens[i+4].Value == "println" || p.tokens[i+4].Value == "print") {
			
			// Buscar el final de la declaración println/print
			endPos := p.findPrintStatementEnd(i)
			if endPos != -1 {
				// Verificar si hay punto y coma después
				if endPos+1 >= len(p.tokens) || p.tokens[endPos+1].Type != "semicolon" {
					p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta ';' después de la declaración System.out.%s", p.tokens[i+4].Line, p.tokens[i+4].Value))
				}
			}
		}

		// Detectar println o print simple
		if (p.tokens[i].Value == "println" || p.tokens[i].Value == "print") {
			// Verificar que no sea parte de System.out.println
			isSystemOut := false
			if i >= 4 && p.tokens[i-4].Value == "System" && p.tokens[i-3].Value == "." &&
			   p.tokens[i-2].Value == "out" && p.tokens[i-1].Value == "." {
				isSystemOut = true
			}

			if !isSystemOut {
				endPos := p.findPrintStatementEnd(i)
				if endPos != -1 {
					if endPos+1 >= len(p.tokens) || p.tokens[endPos+1].Type != "semicolon" {
						p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta ';' después de la declaración %s", p.tokens[i].Line, p.tokens[i].Value))
					}
				}
			}
		}

		// Validar declaraciones de variables (int, char)
		if p.tokens[i].Type == "keyword" && (p.tokens[i].Value == "int" || p.tokens[i].Value == "char") {
			endPos := p.findVariableDeclarationEnd(i)
			if endPos != -1 {
				if endPos+1 >= len(p.tokens) || p.tokens[endPos+1].Type != "semicolon" {
					p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta ';' después de la declaración de variable", p.tokens[i].Line))
				}
			}
		}

		// Validar asignaciones simples (variable = valor)
		if p.tokens[i].Type == "identifier" && i+1 < len(p.tokens) && p.tokens[i+1].Value == "=" {
			// Verificar que no es parte de una declaración de variable
			isDeclaration := false
			if i > 0 && p.tokens[i-1].Type == "keyword" && 
			   (p.tokens[i-1].Value == "int" || p.tokens[i-1].Value == "char") {
				isDeclaration = true
			}

			if !isDeclaration {
				endPos := p.findAssignmentEnd(i)
				if endPos != -1 {
					if endPos+1 >= len(p.tokens) || p.tokens[endPos+1].Type != "semicolon" {
						p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta ';' después de la asignación", p.tokens[i].Line))
					}
				}
			}
		}
	}
}

// Función auxiliar para encontrar el final de una declaración print
func (p *Parser) findPrintStatementEnd(start int) int {
	// Buscar el paréntesis de apertura
	parenStart := -1
	for i := start; i < len(p.tokens) && i < start+10; i++ {
		if p.tokens[i].Type == "lparen" {
			parenStart = i
			break
		}
	}

	if parenStart == -1 {
		return -1
	}

	// Buscar el paréntesis de cierre correspondiente
	parenCount := 0
	for i := parenStart; i < len(p.tokens); i++ {
		if p.tokens[i].Type == "lparen" {
			parenCount++
		} else if p.tokens[i].Type == "rparen" {
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}

	return -1
}

// Función auxiliar para encontrar el final de una declaración de variable
func (p *Parser) findVariableDeclarationEnd(start int) int {
	// Buscar hasta encontrar el final de la declaración
	for i := start + 1; i < len(p.tokens); i++ {
		// Si encontramos un semicolon, el final es el token anterior
		if p.tokens[i].Type == "semicolon" {
			return i - 1
		}
		// Si encontramos una nueva línea o un token que indica nueva declaración
		if p.tokens[i].Type == "keyword" || p.tokens[i].Value == "{" || p.tokens[i].Value == "}" {
			return i - 1
		}
	}
	// Si llegamos al final sin encontrar semicolon
	if len(p.tokens) > start+1 {
		return len(p.tokens) - 1
	}
	return -1
}

// Función auxiliar para encontrar el final de una asignación
func (p *Parser) findAssignmentEnd(start int) int {
	// Buscar el = y luego el valor
	equalPos := -1
	for i := start; i < len(p.tokens) && i < start+5; i++ {
		if p.tokens[i].Value == "=" {
			equalPos = i
			break
		}
	}

	if equalPos == -1 {
		return -1
	}

	// Buscar el final de la asignación
	for i := equalPos + 1; i < len(p.tokens); i++ {
		// Si encontramos un semicolon, el final es el token anterior
		if p.tokens[i].Type == "semicolon" {
			return i - 1
		}
		// Si encontramos tokens que indican nueva declaración
		if p.tokens[i].Type == "keyword" || p.tokens[i].Value == "{" || p.tokens[i].Value == "}" {
			return i - 1
		}
	}
	
	// Si llegamos al final sin encontrar semicolon
	if len(p.tokens) > equalPos+1 {
		return len(p.tokens) - 1
	}
	return -1
}

func (p *Parser) validateStatements() {
	for i := 0; i < len(p.tokens); i++ {
		if p.tokens[i].Value == "for" {
			p.validateForLoop(i)
		}
		if p.tokens[i].Value == "println" || p.tokens[i].Value == "print" {
			p.validatePrintStatement(i)
		}
	}
}

func (p *Parser) validateForLoop(start int) {
	if start+1 >= len(p.tokens) || p.tokens[start+1].Type != "lparen" {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta '(' después de 'for'", p.tokens[start].Line))
		return
	}

	// Buscar el paréntesis de cierre
	parenCount := 0
	endParen := -1
	for i := start + 1; i < len(p.tokens); i++ {
		if p.tokens[i].Type == "lparen" {
			parenCount++
		} else if p.tokens[i].Type == "rparen" {
			parenCount--
			if parenCount == 0 {
				endParen = i
				break
			}
		}
	}

	if endParen == -1 {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta ')' en el for", p.tokens[start].Line))
		return
	}

	// Validar contenido del for
	p.validateForContent(start+2, endParen)

	// Validar que hay llave de apertura después del for
	if endParen+1 >= len(p.tokens) || p.tokens[endParen+1].Type != "lbrace" {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta '{' después del for", p.tokens[endParen].Line))
	}
}

func (p *Parser) validateForContent(start, end int) {
	semicolonCount := 0
	semicolonPos := []int{}

	// Contar punto y coma
	for i := start; i < end; i++ {
		if p.tokens[i].Type == "semicolon" {
			semicolonCount++
			semicolonPos = append(semicolonPos, i)
		}
	}

	if semicolonCount != 2 {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: For debe tener exactamente 2 ';'", p.tokens[start].Line))
		return
	}

	// Validar inicialización
	p.validateForInit(start, semicolonPos[0])

	// Validar condición
	p.validateForCondition(semicolonPos[0]+1, semicolonPos[1])

	// Validar incremento
	p.validateForIncrement(semicolonPos[1]+1, end)
}

func (p *Parser) validateForInit(start, end int) {
	if start >= end {
		p.errors = append(p.errors, "Error: Inicialización del for vacía")
		return
	}

	// Caso 1: Declaración completa (int i = valor)
	if end-start >= 4 && p.tokens[start].Type == "keyword" && 
	   (p.tokens[start].Value == "int" || p.tokens[start].Value == "char") {
		
		if p.tokens[start+1].Type != "identifier" {
			p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba identificador después del tipo", p.tokens[start+1].Line))
			return
		}

		if p.tokens[start+2].Value != "=" {
			p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba '=' en la asignación", p.tokens[start+2].Line))
			return
		}

		// Validar valor inicial según el tipo
		if p.tokens[start].Value == "int" {
			if p.tokens[start+3].Type != "number" {
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba número para variable int", p.tokens[start+3].Line))
			}
		} else if p.tokens[start].Value == "char" {
			if p.tokens[start+3].Type != "char" {
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba char literal para variable char", p.tokens[start+3].Line))
			}
		}
		return
	}

	// Caso 2: Solo identificador (variable ya declarada)
	if end-start == 1 && p.tokens[start].Type == "identifier" {
		// Es válido, solo debe ser un identificador
		return
	}

	// Caso 3: Asignación a variable existente (i = valor)
	if end-start >= 3 && p.tokens[start].Type == "identifier" && p.tokens[start+1].Value == "=" {
		// Validar que hay un valor después del =
		if p.tokens[start+2].Type != "number" && p.tokens[start+2].Type != "char" && p.tokens[start+2].Type != "identifier" {
			p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Valor inválido en asignación", p.tokens[start+2].Line))
		}
		return
	}

	// Si no coincide con ningún patrón válido
	p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Inicialización del for inválida", p.tokens[start].Line))
}

func (p *Parser) validateForCondition(start, end int) {
	if start >= end {
		p.errors = append(p.errors, "Error: Condición del for vacía")
		return
	}

	if end-start < 3 {
		p.errors = append(p.errors, "Error: Condición del for incompleta")
		return
	}

	if p.tokens[start].Type != "identifier" {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba variable en condición", p.tokens[start].Line))
	}

	validOperators := []string{"<", "<=", ">", ">=", "==", "!="}
	isValidOp := false
	for _, op := range validOperators {
		if p.tokens[start+1].Value == op {
			isValidOp = true
			break
		}
	}

	if !isValidOp {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Operador inválido en condición", p.tokens[start+1].Line))
	}

	if p.tokens[start+2].Type != "number" && p.tokens[start+2].Type != "char" && p.tokens[start+2].Type != "identifier" {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Valor inválido en condición", p.tokens[start+2].Line))
	}
}

func (p *Parser) validateForIncrement(start, end int) {
	if start >= end {
		p.errors = append(p.errors, "Error: Incremento del for vacío")
		return
	}

	if end-start < 2 {
		p.errors = append(p.errors, "Error: Incremento del for incompleto")
		return
	}

	if p.tokens[start].Type != "identifier" {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Se esperaba variable en incremento", p.tokens[start].Line))
		return
	}

	// Verificar que no haya operadores múltiples como c+++
	if end-start > 2 {
		// Verificar si hay múltiples operadores seguidos
		for i := start + 1; i < end-1; i++ {
			if p.tokens[i].Type == "operator" && p.tokens[i+1].Type == "operator" {
				// Casos válidos: ++ o --
				if (p.tokens[i].Value == "+" && p.tokens[i+1].Value == "+") ||
				   (p.tokens[i].Value == "-" && p.tokens[i+1].Value == "-") {
					continue
				}
				// Casos inválidos como c+++
				p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Operadores múltiples inválidos en incremento", p.tokens[i].Line))
				return
			}
		}
	}

	validIncrements := []string{"++", "--", "+=", "-="}
	isValidIncrement := false
	for _, inc := range validIncrements {
		if p.tokens[start+1].Value == inc {
			isValidIncrement = true
			break
		}
	}

	if !isValidIncrement {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Operador de incremento inválido", p.tokens[start+1].Line))
		return
	}

	// Verificar que no haya tokens adicionales después del incremento válido
	if (p.tokens[start+1].Value == "++" || p.tokens[start+1].Value == "--") && end-start > 2 {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Operadores adicionales después de %s", p.tokens[start+1].Line, p.tokens[start+1].Value))
		return
	}

	// Si es += o -=, debe haber un valor después
	if (p.tokens[start+1].Value == "+=" || p.tokens[start+1].Value == "-=") && end-start < 3 {
		p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Falta valor después de %s", p.tokens[start+1].Line, p.tokens[start+1].Value))
	}
}

func (p *Parser) validatePrintStatement(printPos int) {
	// Buscar el paréntesis de apertura
	parenStart := -1
	for i := printPos + 1; i < len(p.tokens) && i < printPos + 3; i++ {
		if p.tokens[i].Type == "lparen" {
			parenStart = i
			break
		}
	}

	if parenStart == -1 {
		return // No es un error crítico, puede ser válido
	}

	// Buscar el paréntesis de cierre
	parenCount := 0
	parenEnd := -1
	for i := parenStart; i < len(p.tokens); i++ {
		if p.tokens[i].Type == "lparen" {
			parenCount++
		} else if p.tokens[i].Type == "rparen" {
			parenCount--
			if parenCount == 0 {
				parenEnd = i
				break
			}
		}
	}

	if parenEnd == -1 {
		return // Error ya manejado por validateBrackets
	}

	// Validar el contenido dentro del println
	p.validatePrintContent(parenStart + 1, parenEnd)
}

func (p *Parser) validatePrintContent(start, end int) {
	for i := start; i < end-1; i++ {
		// Buscar patrones problemáticos en concatenación
		if p.tokens[i].Type == "string" && i+1 < end {
			// Verificar concatenación después de string
			if p.tokens[i+1].Value == "+" && i+2 < end {
				// Caso válido: "texto" + variable
				if p.tokens[i+2].Type == "identifier" {
					continue
				}
				// Caso inválido: "texto" ++ variable
				if p.tokens[i+1].Value == "+" && i+2 < end && p.tokens[i+2].Value == "+" {
					p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Uso incorrecto de '++' en concatenación, use solo '+'", p.tokens[i+1].Line))
				}
			}
		}

		// Verificar el patrón específico: " ++ variable"
		if p.tokens[i].Value == "+" && i+1 < end && p.tokens[i+1].Value == "+" && i+2 < end && p.tokens[i+2].Type == "identifier" {
			p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Uso incorrecto de '++' en concatenación, use solo '+'", p.tokens[i].Line))
		}

		// También verificar si hay ++ usado como concatenación en cualquier contexto
		if p.tokens[i].Value == "++" && 
		   ((i > start && (p.tokens[i-1].Type == "string" || p.tokens[i-1].Type == "identifier")) ||
		    (i+1 < end && (p.tokens[i+1].Type == "string" || p.tokens[i+1].Type == "identifier"))) {
			p.errors = append(p.errors, fmt.Sprintf("Error línea %d: Uso incorrecto de '++' para concatenación, use '+' para concatenar", p.tokens[i].Line))
		}
	}
}