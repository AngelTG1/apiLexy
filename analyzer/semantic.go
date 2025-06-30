package analyzer

import (
	"fmt"
	"strconv"
)

type Variable struct {
	Name  string
	Type  string
	Value interface{}
	Line  int
}

func AnalyzeSemantics(tokens []Token) (bool, []string) {
	errors := []string{}
	declaredVars := make(map[string]Variable)

	// Primera pasada: declaraciones de variables (incluyendo las del for)
	for i := 0; i < len(tokens); i++ {
		// AGREGAR String A LOS TIPOS RECONOCIDOS
		if tokens[i].Type == "keyword" && (tokens[i].Value == "int" || tokens[i].Value == "char" || tokens[i].Value == "float" || tokens[i].Value == "String") {
			if i+1 < len(tokens) && tokens[i+1].Type == "identifier" {
				varName := tokens[i+1].Value
				varType := tokens[i].Value

				// Verificar si ya está declarada
				if _, exists := declaredVars[varName]; exists {
					errors = append(errors, fmt.Sprintf("Error línea %d: Variable '%s' ya está declarada", tokens[i+1].Line, varName))
				} else {
					// Verificar si hay inicialización
					var value interface{}
					
					if i+3 < len(tokens) && tokens[i+2].Value == "=" {
						// VALIDACIÓN DE TIPOS CORREGIDA
						valueToken := tokens[i+3]
						
						if varType == "int" {
							// Para int solo aceptamos números enteros
							if valueToken.Type == "number" {
								if val, err := strconv.Atoi(valueToken.Value); err == nil {
									value = val
								} else {
									errors = append(errors, fmt.Sprintf("Error línea %d: Valor '%s' no es un entero válido", valueToken.Line, valueToken.Value))
								}
							} else if valueToken.Type == "float" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "string" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "char" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
							} else {
								errors = append(errors, fmt.Sprintf("Error línea %d: Tipo incompatible para variable int '%s'", valueToken.Line, varName))
							}
						} else if varType == "float" {
							// Para float solo aceptamos números decimales
							if valueToken.Type == "float" {
								if val, err := strconv.ParseFloat(valueToken.Value, 64); err == nil {
									value = val
								} else {
									errors = append(errors, fmt.Sprintf("Error línea %d: Valor '%s' no es un float válido", valueToken.Line, valueToken.Value))
								}
							} else if valueToken.Type == "number" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar entero '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "string" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "char" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
							} else {
								errors = append(errors, fmt.Sprintf("Error línea %d: Tipo incompatible para variable float '%s'", valueToken.Line, varName))
							}
						} else if varType == "char" {
							// Para char solo aceptamos caracteres
							if valueToken.Type == "char" {
								if len(valueToken.Value) == 1 {
									value = rune(valueToken.Value[0])
								} else {
									errors = append(errors, fmt.Sprintf("Error línea %d: Char literal inválido '%s'", valueToken.Line, valueToken.Value))
								}
							} else if valueToken.Type == "number" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar número '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "float" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "string" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
							} else {
								errors = append(errors, fmt.Sprintf("Error línea %d: Tipo incompatible para variable char '%s'", valueToken.Line, varName))
							}
						} else if varType == "String" {
							// NUEVO: Para String solo aceptamos strings
							if valueToken.Type == "string" {
								value = valueToken.Value
							} else if valueToken.Type == "number" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar número '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "float" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
							} else if valueToken.Type == "char" {
								errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
							} else {
								errors = append(errors, fmt.Sprintf("Error línea %d: Tipo incompatible para variable String '%s'", valueToken.Line, varName))
							}
						}
					}

					// CAMBIO CRÍTICO: Declarar la variable SIEMPRE, incluso si hay errores de tipo
					// Esto evita errores posteriores de "variable no declarada"
					declaredVars[varName] = Variable{
						Name:  varName,
						Type:  varType,
						Value: value,
						Line:  tokens[i+1].Line,
					}
				}
			}
		}
	}

	// Segunda pasada: uso de variables y validación de asignaciones
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "identifier" {
			varName := tokens[i].Value

			// Verificar si es una palabra reservada
			if keywords[varName] {
				continue
			}

			// CORREGIR: Verificar si estamos en un contexto de declaración de variable
			// Incluir String en la verificación
			if i > 0 && tokens[i-1].Type == "keyword" && 
			   (tokens[i-1].Value == "int" || tokens[i-1].Value == "char" || tokens[i-1].Value == "float" || tokens[i-1].Value == "String") {
				continue // Es una declaración, no un uso
			}

			// Verificar asignaciones a variables ya declaradas
			if i+2 < len(tokens) && tokens[i+1].Value == "=" {
				if declaredVar, exists := declaredVars[varName]; exists {
					valueToken := tokens[i+2]
					
					// Validar compatibilidad de tipos en asignación
					if declaredVar.Type == "int" {
						if valueToken.Type == "string" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "char" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "float" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable int '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "number" {
							if _, err := strconv.Atoi(valueToken.Value); err != nil {
								errors = append(errors, fmt.Sprintf("Error línea %d: Valor '%s' no es un entero válido", valueToken.Line, valueToken.Value))
							}
						}
					} else if declaredVar.Type == "float" {
						if valueToken.Type == "string" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "char" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "number" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar entero '%s' a variable float '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "float" {
							if _, err := strconv.ParseFloat(valueToken.Value, 64); err != nil {
								errors = append(errors, fmt.Sprintf("Error línea %d: Valor '%s' no es un float válido", valueToken.Line, valueToken.Value))
							}
						}
					} else if declaredVar.Type == "char" {
						if valueToken.Type == "string" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar string '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "number" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar número '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "float" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable char '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "char" && len(valueToken.Value) != 1 {
							errors = append(errors, fmt.Sprintf("Error línea %d: Char literal inválido '%s'", valueToken.Line, valueToken.Value))
						}
					} else if declaredVar.Type == "String" {
						// NUEVO: Validación para String
						if valueToken.Type == "number" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar número '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "char" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar char '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
						} else if valueToken.Type == "float" {
							errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar float '%s' a variable String '%s'", valueToken.Line, valueToken.Value, varName))
						}
					}
				} else {
					errors = append(errors, fmt.Sprintf("Error línea %d: Variable '%s' usada sin declarar", tokens[i].Line, varName))
				}
			} else {
				// Verificar si la variable está declarada (uso normal)
				if _, exists := declaredVars[varName]; !exists {
					errors = append(errors, fmt.Sprintf("Error línea %d: Variable '%s' usada sin declarar", tokens[i].Line, varName))
				}
			}
		}
	}

	// Tercera pasada: validación específica de for loops
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Value == "for" {
			errors = append(errors, validateForSemantics(tokens, i, declaredVars)...)
		}
	}

	// Validación de rangos para char
	errors = append(errors, validateCharRanges(tokens, declaredVars)...)

	return len(errors) == 0, errors
}

func validateForSemantics(tokens []Token, forPos int, declaredVars map[string]Variable) []string {
	errors := []string{}

	// Buscar los componentes del for
	parenStart := -1
	parenEnd := -1
	parenCount := 0

	for i := forPos + 1; i < len(tokens); i++ {
		if tokens[i].Type == "lparen" {
			if parenStart == -1 {
				parenStart = i
			}
			parenCount++
		} else if tokens[i].Type == "rparen" {
			parenCount--
			if parenCount == 0 {
				parenEnd = i
				break
			}
		}
	}

	if parenStart == -1 || parenEnd == -1 {
		return errors
	}

	// Buscar punto y coma
	semicolons := []int{}
	for i := parenStart + 1; i < parenEnd; i++ {
		if tokens[i].Type == "semicolon" {
			semicolons = append(semicolons, i)
		}
	}

	if len(semicolons) != 2 {
		return errors
	}

	// Obtener los rangos de cada parte del for
	initStart := parenStart + 1
	initEnd := semicolons[0]
	condStart := semicolons[0] + 1
	condEnd := semicolons[1]
	incrStart := semicolons[1] + 1
	incrEnd := parenEnd

	var forVar string
	var forVarType string

	// Analizar la inicialización del for con validación de tipos
	if initStart < initEnd {
		// Caso 1: Declaración completa (int i = valor)
		if tokens[initStart].Type == "keyword" && (tokens[initStart].Value == "int" || tokens[initStart].Value == "char" || tokens[initStart].Value == "float" || tokens[initStart].Value == "String") {
			if initStart+1 < initEnd && tokens[initStart+1].Type == "identifier" {
				forVar = tokens[initStart+1].Value
				forVarType = tokens[initStart].Value
				
				// Validar la inicialización en el for
				if initStart+3 < initEnd && tokens[initStart+2].Value == "=" {
					valueToken := tokens[initStart+3]
					
					if forVarType == "int" && valueToken.Type != "number" {
						errors = append(errors, fmt.Sprintf("Error línea %d: Variable int '%s' en for debe inicializarse con número", valueToken.Line, forVar))
					} else if forVarType == "float" && valueToken.Type != "float" {
						errors = append(errors, fmt.Sprintf("Error línea %d: Variable float '%s' en for debe inicializarse con float", valueToken.Line, forVar))
					} else if forVarType == "char" && valueToken.Type != "char" {
						errors = append(errors, fmt.Sprintf("Error línea %d: Variable char '%s' en for debe inicializarse con char", valueToken.Line, forVar))
					} else if forVarType == "String" && valueToken.Type != "string" {
						errors = append(errors, fmt.Sprintf("Error línea %d: Variable String '%s' en for debe inicializarse con string", valueToken.Line, forVar))
					}
				}
			}
		} else if tokens[initStart].Type == "identifier" {
			// Caso 2: Solo variable (i) o asignación (i = valor)
			forVar = tokens[initStart].Value
			
			// Verificar que la variable esté declarada previamente
			if declaredVar, exists := declaredVars[forVar]; exists {
				forVarType = declaredVar.Type
				
				// Si hay asignación, validar tipo
				if initStart+2 < initEnd && tokens[initStart+1].Value == "=" {
					valueToken := tokens[initStart+2]
					
					if forVarType == "int" && valueToken.Type != "number" {
						errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar '%s' a variable int '%s'", valueToken.Line, valueToken.Value, forVar))
					} else if forVarType == "float" && valueToken.Type != "float" {
						errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar '%s' a variable float '%s'", valueToken.Line, valueToken.Value, forVar))
					} else if forVarType == "char" && valueToken.Type != "char" {
						errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar '%s' a variable char '%s'", valueToken.Line, valueToken.Value, forVar))
					} else if forVarType == "String" && valueToken.Type != "string" {
						errors = append(errors, fmt.Sprintf("Error línea %d: No se puede asignar '%s' a variable String '%s'", valueToken.Line, valueToken.Value, forVar))
					}
				}
			} else {
				errors = append(errors, fmt.Sprintf("Error línea %d: Variable '%s' en for no está declarada", tokens[initStart].Line, forVar))
				return errors
			}
		}
	}

	// Validar que la variable de condición coincida
	if condStart < condEnd && tokens[condStart].Type == "identifier" {
		condVar := tokens[condStart].Value
		if forVar != "" && condVar != forVar {
			errors = append(errors, fmt.Sprintf("Error línea %d: Variable en condición '%s' no coincide con variable del for '%s'", tokens[condStart].Line, condVar, forVar))
		}
	}

	// Validar que la variable de incremento coincida
	if incrStart < incrEnd && tokens[incrStart].Type == "identifier" {
		incrVar := tokens[incrStart].Value
		if forVar != "" && incrVar != forVar {
			errors = append(errors, fmt.Sprintf("Error línea %d: Variable en incremento '%s' no coincide con variable del for '%s'", tokens[incrStart].Line, incrVar, forVar))
		}
	}

	// Validación específica de tipos en la condición
	if forVarType != "" && condStart+2 < condEnd {
		condValueToken := tokens[condStart+2]
		if forVarType == "int" && condValueToken.Type == "number" {
			// Validar que sea un número entero válido
			if _, err := strconv.Atoi(condValueToken.Value); err != nil {
				errors = append(errors, fmt.Sprintf("Error línea %d: Valor inválido para comparación con int", condValueToken.Line))
			}
		} else if forVarType == "float" && condValueToken.Type == "float" {
			// Validar que sea un float válido
			if _, err := strconv.ParseFloat(condValueToken.Value, 64); err != nil {
				errors = append(errors, fmt.Sprintf("Error línea %d: Valor inválido para comparación con float", condValueToken.Line))
			}
		} else if forVarType == "char" && condValueToken.Type == "char" {
			// Validar que sea un char válido
			if len(condValueToken.Value) != 1 {
				errors = append(errors, fmt.Sprintf("Error línea %d: Char literal inválido en condición", condValueToken.Line))
			}
		} else if forVarType == "String" && condValueToken.Type == "string" {
			// String comparisons are valid
		} else if forVarType == "int" && condValueToken.Type != "number" {
			errors = append(errors, fmt.Sprintf("Error línea %d: No se puede comparar int con %s", condValueToken.Line, condValueToken.Type))
		} else if forVarType == "float" && condValueToken.Type != "float" {
			errors = append(errors, fmt.Sprintf("Error línea %d: No se puede comparar float con %s", condValueToken.Line, condValueToken.Type))
		} else if forVarType == "char" && condValueToken.Type != "char" {
			errors = append(errors, fmt.Sprintf("Error línea %d: No se puede comparar char con %s", condValueToken.Line, condValueToken.Type))
		} else if forVarType == "String" && condValueToken.Type != "string" {
			errors = append(errors, fmt.Sprintf("Error línea %d: No se puede comparar String con %s", condValueToken.Line, condValueToken.Type))
		}
	}

	return errors
}

func validateCharRanges(tokens []Token, declaredVars map[string]Variable) []string {
	errors := []string{}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "char" {
			char := tokens[i].Value
			if len(char) == 1 {
				r := rune(char[0])
				if r < 'a' || r > 'z' {
					errors = append(errors, fmt.Sprintf("Error línea %d: Char '%s' fuera del rango permitido (a-z)", tokens[i].Line, char))
				}
			}
		}
	}

	return errors
}