// analyzer/enhanced_semantic.go
package analyzer

import (
	"fmt"
	"strconv"
)

// EnhancedSemanticAnalyzer analizador semántico mejorado con optimizaciones
type EnhancedSemanticAnalyzer struct {
	stringLib     *StringLibrary
	variableCache map[string]Variable
	errorBuffer   []string
}

// NewEnhancedSemanticAnalyzer crea un nuevo analizador semántico optimizado
func NewEnhancedSemanticAnalyzer() *EnhancedSemanticAnalyzer {
	return &EnhancedSemanticAnalyzer{
		stringLib:     NewStringLibrary(),
		variableCache: make(map[string]Variable, 100),
		errorBuffer:   make([]string, 0, 50),
	}
}

// AnalyzeOptimized análisis semántico optimizado
func (esa *EnhancedSemanticAnalyzer) AnalyzeOptimized(tokens []Token) (bool, []string) {
	// Limpiar cache y buffer para nuevo análisis
	for k := range esa.variableCache {
		delete(esa.variableCache, k)
	}
	esa.errorBuffer = esa.errorBuffer[:0]
	
	// Análisis en múltiples pasadas optimizadas
	esa.analyzeDeclarations(tokens)
	esa.analyzeUsage(tokens)
	esa.analyzeStringMethods(tokens)
	esa.analyzeTypeCompatibility(tokens)
	
	return len(esa.errorBuffer) == 0, esa.errorBuffer
}

// analyzeDeclarations analiza declaraciones de variables
func (esa *EnhancedSemanticAnalyzer) analyzeDeclarations(tokens []Token) {
	supportedTypes := map[string]bool{
		"int": true, "char": true, "float": true, "String": true,
		"double": true, "boolean": true, "byte": true, "short": true, "long": true,
	}
	
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "keyword" && supportedTypes[tokens[i].Value] {
			if i+1 < len(tokens) && tokens[i+1].Type == "identifier" {
				varName := esa.stringLib.InternString(tokens[i+1].Value)
				varType := tokens[i].Value
				
				if _, exists := esa.variableCache[varName]; exists {
					esa.addError(fmt.Sprintf("Error línea %d: Variable '%s' ya está declarada", tokens[i+1].Line, varName))
					continue
				}
				
				var value interface{}
				if i+3 < len(tokens) && tokens[i+2].Value == "=" {
					value = esa.validateInitialization(tokens[i+3], varType)
				}
				
				esa.variableCache[varName] = Variable{
					Name:  varName,
					Type:  varType,
					Value: value,
					Line:  tokens[i+1].Line,
				}
			}
		}
	}
}

// analyzeStringMethods analiza métodos de String
func (esa *EnhancedSemanticAnalyzer) analyzeStringMethods(tokens []Token) {
	for i := 0; i < len(tokens)-2; i++ {
		if tokens[i].Type == "identifier" && tokens[i+1].Value == "." && tokens[i+2].Type == "identifier" {
			varName := tokens[i].Value
			methodName := tokens[i+2].Value
			
			if variable, exists := esa.variableCache[varName]; exists && variable.Type == "String" {
				if !esa.stringLib.ValidateStringMethod(methodName) {
					esa.addError(fmt.Sprintf("Error línea %d: Método '%s' no válido para String", tokens[i+2].Line, methodName))
				}
			}
		}
	}
}

// validateInitialization valida la inicialización de variables
func (esa *EnhancedSemanticAnalyzer) validateInitialization(valueToken Token, varType string) interface{} {
	switch varType {
	case "int":
		if valueToken.Type == "number" {
			if val, err := strconv.Atoi(valueToken.Value); err == nil {
				return val
			}
		}
		esa.addError(fmt.Sprintf("Error línea %d: Valor inválido para tipo int", valueToken.Line))
	case "float", "double":
		if valueToken.Type == "float" || valueToken.Type == "number" {
			if val, err := strconv.ParseFloat(valueToken.Value, 64); err == nil {
				return val
			}
		}
		esa.addError(fmt.Sprintf("Error línea %d: Valor inválido para tipo %s", valueToken.Line, varType))
	case "String":
		if valueToken.Type == "string" {
			return valueToken.Value
		}
		esa.addError(fmt.Sprintf("Error línea %d: Se esperaba string para tipo String", valueToken.Line))
	case "char":
		if valueToken.Type == "char" && len(valueToken.Value) == 1 {
			return rune(valueToken.Value[0])
		}
		esa.addError(fmt.Sprintf("Error línea %d: Valor inválido para tipo char", valueToken.Line))
	case "boolean":
		if valueToken.Value == "true" || valueToken.Value == "false" {
			return valueToken.Value == "true"
		}
		esa.addError(fmt.Sprintf("Error línea %d: Se esperaba true o false para tipo boolean", valueToken.Line))
	}
	return nil
}

// analyzeUsage analiza el uso de variables
func (esa *EnhancedSemanticAnalyzer) analyzeUsage(tokens []Token) {
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == "identifier" {
			varName := tokens[i].Value
			
			// Verificar si es palabra reservada
			if keywords[varName] {
				continue
			}
			
			// Verificar si es declaración
			if i > 0 && tokens[i-1].Type == "keyword" {
				continue
			}
			
			// Verificar si la variable está declarada
			if _, exists := esa.variableCache[varName]; !exists {
				esa.addError(fmt.Sprintf("Error línea %d: Variable '%s' usada sin declarar", tokens[i].Line, varName))
			}
		}
	}
}

// analyzeTypeCompatibility analiza compatibilidad de tipos
func (esa *EnhancedSemanticAnalyzer) analyzeTypeCompatibility(tokens []Token) {
	for i := 0; i < len(tokens)-2; i++ {
		if tokens[i].Type == "identifier" && tokens[i+1].Value == "=" {
			varName := tokens[i].Value
			valueToken := tokens[i+2]
			
			if variable, exists := esa.variableCache[varName]; exists {
				esa.validateAssignment(variable, valueToken)
			}
		}
	}
}

// validateAssignment valida asignaciones
func (esa *EnhancedSemanticAnalyzer) validateAssignment(variable Variable, valueToken Token) {
	switch variable.Type {
	case "int":
		if valueToken.Type != "number" {
			esa.addError(fmt.Sprintf("Error línea %d: No se puede asignar %s a variable int '%s'", valueToken.Line, valueToken.Type, variable.Name))
		}
	case "String":
		if valueToken.Type != "string" {
			esa.addError(fmt.Sprintf("Error línea %d: No se puede asignar %s a variable String '%s'", valueToken.Line, valueToken.Type, variable.Name))
		}
	case "char":
		if valueToken.Type != "char" {
			esa.addError(fmt.Sprintf("Error línea %d: No se puede asignar %s a variable char '%s'", valueToken.Line, valueToken.Type, variable.Name))
		}
	case "float", "double":
		if valueToken.Type != "float" && valueToken.Type != "number" {
			esa.addError(fmt.Sprintf("Error línea %d: No se puede asignar %s a variable %s '%s'", valueToken.Line, valueToken.Type, variable.Type, variable.Name))
		}
	}
}

// addError agrega un error al buffer
func (esa *EnhancedSemanticAnalyzer) addError(error string) {
	esa.errorBuffer = append(esa.errorBuffer, error)
}