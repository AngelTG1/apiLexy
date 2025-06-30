// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"java-analyzer/analyzer"
)

type Request struct {
	Code string `json:"code"`
}

type Response struct {
	Tokens         []analyzer.Token `json:"tokens"`
	SyntaxOK       bool             `json:"syntax_ok"`
	SemanticOK     bool             `json:"semantic_ok"`
	SyntaxErrors   []string         `json:"syntax_errors"`
	SemanticErrors []string         `json:"semantic_errors"`
	AnalysisTime   string           `json:"analysis_time"`
	Summary        AnalysisSummary  `json:"summary"`
}

type AnalysisSummary struct {
	TotalTokens       int `json:"total_tokens"`
	LinesAnalyzed     int `json:"lines_analyzed"`
	VariablesFound    int `json:"variables_found"`
	MethodsFound      int `json:"methods_found"`
	ErrorCount        int `json:"error_count"`
	WarningCount      int `json:"warning_count"`
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	
	// Habilitar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error leyendo request body: %v", err)
		http.Error(w, "Error leyendo el cuerpo de la petición", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Printf("Error parseando JSON: %v", err)
		http.Error(w, "Error parseando JSON", http.StatusBadRequest)
		return
	}

	// Validar que el código no esté vacío
	if len(req.Code) == 0 {
		http.Error(w, "El código no puede estar vacío", http.StatusBadRequest)
		return
	}

	log.Printf("Analizando código de %d caracteres", len(req.Code))

	// Análisis léxico
	tokens := analyzer.Lex(req.Code)
	log.Printf("Lexer generó %d tokens", len(tokens))

	// Análisis sintáctico
	syntaxOK, syntaxErrors := analyzer.Parse(tokens)
	log.Printf("Parser encontró %d errores sintácticos", len(syntaxErrors))

	// Análisis semántico
	semanticOK, semanticErrors := analyzer.AnalyzeSemantics(tokens)
	log.Printf("Analizador semántico encontró %d errores", len(semanticErrors))

	// Generar resumen
	summary := generateSummary(tokens, syntaxErrors, semanticErrors, req.Code)
	
	analysisTime := time.Since(startTime)

	res := Response{
		Tokens:         tokens,
		SyntaxOK:       syntaxOK,
		SemanticOK:     semanticOK,
		SyntaxErrors:   syntaxErrors,
		SemanticErrors: semanticErrors,
		AnalysisTime:   analysisTime.String(),
		Summary:        summary,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Error codificando respuesta: %v", err)
		http.Error(w, "Error codificando respuesta", http.StatusInternalServerError)
		return
	}

	log.Printf("Análisis completado en %v", analysisTime)
}

func generateSummary(tokens []analyzer.Token, syntaxErrors, semanticErrors []string, code string) AnalysisSummary {
	// Contar líneas
	lines := 1
	for _, char := range code {
		if char == '\n' {
			lines++
		}
	}

	// Contar variables y métodos
	variables := 0
	methods := 0
	
	for i, token := range tokens {
		// Contar declaraciones de variables
		if token.Type == "keyword" && isVariableType(token.Value) {
			if i+1 < len(tokens) && tokens[i+1].Type == "identifier" {
				variables++
			}
		}
		
		// Contar métodos (simplificado)
		if token.Value == "main" || (token.Type == "identifier" && i+1 < len(tokens) && tokens[i+1].Type == "lparen") {
			// Verificar que no sea una llamada sino una declaración
			if i > 0 && (tokens[i-1].Value == "void" || isVariableType(tokens[i-1].Value)) {
				methods++
			}
		}
	}

	// Contar errores y warnings
	errorCount := 0
	warningCount := 0
	
	for _, err := range syntaxErrors {
		if contains(err, "Error") {
			errorCount++
		} else if contains(err, "Advertencia") {
			warningCount++
		}
	}
	
	for _, err := range semanticErrors {
		if contains(err, "Error") {
			errorCount++
		} else if contains(err, "Advertencia") {
			warningCount++
		}
	}

	return AnalysisSummary{
		TotalTokens:    len(tokens),
		LinesAnalyzed:  lines,
		VariablesFound: variables,
		MethodsFound:   methods,
		ErrorCount:     errorCount,
		WarningCount:   warningCount,
	}
}

func isVariableType(value string) bool {
	types := []string{"int", "char", "String", "float", "double", "boolean", "byte", "short", "long"}
	for _, t := range types {
		if value == t {
			return true
		}
	}
	return false
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}

// Endpoint para obtener información del analizador
func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	info := map[string]interface{}{
		"name":        "Analizador de Java",
		"version":     "2.0.0",
		"description": "Analizador léxico, sintáctico y semántico para Java",
		"features": []string{
			"Análisis léxico completo",
			"Detección de palabras clave pegadas",
			"Validación de estructura Java",
			"Análisis semántico avanzado",
			"Detección de variables no utilizadas",
			"Validación de tipos",
			"Sugerencias de mejora",
		},
		"supported_constructs": []string{
			"Clases públicas",
			"Método main",
			"Variables (int, String, char, float, double, boolean)",
			"Estructuras de control (if, for, while)",
			"System.out.println",
			"Métodos de String (equals, length, etc.)",
		},
	}

	json.NewEncoder(w).Encode(info)
}

// Endpoint para validar solo sintaxis
func syntaxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error parseando JSON", http.StatusBadRequest)
		return
	}

	tokens := analyzer.Lex(req.Code)
	syntaxOK, syntaxErrors := analyzer.Parse(tokens)

	response := map[string]interface{}{
		"syntax_ok":     syntaxOK,
		"syntax_errors": syntaxErrors,
		"token_count":   len(tokens),
	}

	json.NewEncoder(w).Encode(response)
}

// Middleware para logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Iniciando %s %s", r.Method, r.URL.Path)
		
		next.ServeHTTP(w, r)
		
		log.Printf("Completado %s %s en %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	// Configurar logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Rutas
	http.Handle("/analyze", loggingMiddleware(http.HandlerFunc(analyzeHandler)))
	http.Handle("/syntax", loggingMiddleware(http.HandlerFunc(syntaxHandler)))
	http.Handle("/info", loggingMiddleware(http.HandlerFunc(infoHandler)))
	
	// Ruta de prueba
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	fmt.Println("🚀 Analizador de Java v2.0")
	fmt.Println("==============================")
	fmt.Println("📡 Servidor iniciado en http://localhost:8080")
	fmt.Println("📋 Endpoints disponibles:")
	fmt.Println("   • POST /analyze   - Análisis completo (léxico + sintáctico + semántico)")
	fmt.Println("   • POST /syntax    - Solo análisis sintáctico")
	fmt.Println("   • GET  /info      - Información del analizador")
	fmt.Println("   • GET  /health    - Estado del servidor")
	fmt.Println("==============================")
	fmt.Println("✨ Características:")
	fmt.Println("   • Detección de palabras clave pegadas (voidmain → void main)")
	fmt.Println("   • Validación completa de estructura Java")
	fmt.Println("   • Análisis semántico avanzado")
	fmt.Println("   • Detección de variables no utilizadas")
	fmt.Println("   • Sugerencias de mejora de código")
	fmt.Println("==============================")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}