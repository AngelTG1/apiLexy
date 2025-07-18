// analyzer/string_utils.go
package analyzer


// StringLibrary proporciona funcionalidades optimizadas para el análisis de strings
type StringLibrary struct {
	// Pool de strings reutilizables para reducir allocaciones
	stringPool map[string]string
	// Cache de resultados de validación
	validationCache map[string]bool
}

// NewStringLibrary crea una nueva instancia de la librería de strings
func NewStringLibrary() *StringLibrary {
	return &StringLibrary{
		stringPool:      make(map[string]string, 1000),
		validationCache: make(map[string]bool, 500),
	}
}

// InternString reutiliza strings para reducir memoria
func (sl *StringLibrary) InternString(s string) string {
	if interned, exists := sl.stringPool[s]; exists {
		return interned
	}
	sl.stringPool[s] = s
	return s
}

// ValidateJavaString valida si un string cumple con las reglas de Java
func (sl *StringLibrary) ValidateJavaString(s string) bool {
	if cached, exists := sl.validationCache[s]; exists {
		return cached
	}
	
	result := sl.validateStringContent(s)
	sl.validationCache[s] = result
	return result
}

func (sl *StringLibrary) validateStringContent(s string) bool {
	// Validar caracteres escapados válidos
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			next := s[i+1]
			switch next {
			case 'n', 't', 'r', '\\', '"', '\'':
				i++ // saltar el carácter escapado
			default:
				return false
			}
		}
	}
	return true
}

// GetStringMethods retorna métodos disponibles para String en Java
func (sl *StringLibrary) GetStringMethods() []string {
	return []string{
		"length", "charAt", "substring", "indexOf", "lastIndexOf",
		"equals", "equalsIgnoreCase", "compareTo", "compareToIgnoreCase",
		"startsWith", "endsWith", "contains", "toUpperCase", "toLowerCase",
		"trim", "replace", "replaceAll", "split", "valueOf", "toString",
	}
}

// ValidateStringMethod valida si un método de String es válido
func (sl *StringLibrary) ValidateStringMethod(method string) bool {
	methods := sl.GetStringMethods()
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

