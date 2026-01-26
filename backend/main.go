// main.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"synta-compiler/lexer"
	"synta-compiler/parser"
	"synta-compiler/token"
	"time"
)

type analyzeReq struct {
	Code     string `json:"code"`
	Filename string `json:"filename,omitempty"`
}

type tokenDTO struct {
	Lexeme        string      `json:"lexeme"`
	Type          string      `json:"type"`
	SemanticGroup string      `json:"semanticGroup"`
	Line          int         `json:"line"`
	Column        int         `json:"column"`
	Extra         interface{} `json:"extra,omitempty"`
}

type analyzeResp struct {
	Tokens    []tokenDTO      `json:"tokens,omitempty"`
	ParseTree string          `json:"parseTree,omitempty"`
	Errors    []parseErrorDTO `json:"errors,omitempty"`
	Warnings  []string        `json:"warnings,omitempty"`
	Success   bool            `json:"success"`
	Error     string          `json:"error,omitempty"`
}

type parseErrorDTO struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(analyzeResp{Success: false, Error: "method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(analyzeResp{Success: false, Error: "unable to read body"})
		return
	}

	var req analyzeReq
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(analyzeResp{Success: false, Error: "invalid json"})
		return
	}

	// Lexical analysis
	l := lexer.New(req.Code)
	toks := l.Tokenize()

	// Convert tokens to DTOs
	tokenDTOs := make([]tokenDTO, 0, len(toks))
	for _, t := range toks {
		typeName := token.TokenNames[t.Type]
		if typeName == "" {
			typeName = "UNKNOWN"
		}
		tokenDTOs = append(tokenDTOs, tokenDTO{
			Lexeme:        t.Lexeme,
			Type:          typeName,
			SemanticGroup: t.Type.GetSemanticGroup(),
			Line:          t.Line,
			Column:        t.Column,
		})
	}

	// Syntactic analysis
	p := parser.NewParser(toks)
	parseResult := p.Parse()

	// Convert parse errors to DTOs
	errorDTOs := make([]parseErrorDTO, 0, len(parseResult.Errors))
	for _, err := range parseResult.Errors {
		errorDTOs = append(errorDTOs, parseErrorDTO{
			Line:    err.Line,
			Column:  err.Column,
			Message: err.Message,
			Type:    err.Type,
		})
	}

	// Format AST
	parseTree := ""
	if parseResult.Program != nil {
		parseTree = parser.FormatAST(parseResult.Program, 0)
	}

	success := len(parseResult.Errors) == 0

	resp := analyzeResp{
		Tokens:    tokenDTOs,
		ParseTree: parseTree,
		Errors:    errorDTOs,
		Warnings:  parseResult.Warnings,
		Success:   success,
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/analyze", analyzeHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Analyzer HTTP server listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
