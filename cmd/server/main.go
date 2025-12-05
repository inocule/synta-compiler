package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"synta-lexical/lexer"
	"synta-lexical/token"
)

type analyzeReq struct {
	Code     string `json:"code"`
	Filename string `json:"filename,omitempty"`
}

type tokenDTO struct {
	Lexeme string      `json:"lexeme"`
	Type   string      `json:"type"`
	Line   int         `json:"line"`
	Column int         `json:"column"`
	Extra  interface{} `json:"extra,omitempty"`
}

type analyzeResp struct {
	Tokens []tokenDTO `json:"tokens,omitempty"`
	Error  string     `json:"error,omitempty"`
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
		json.NewEncoder(w).Encode(analyzeResp{Error: "method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(analyzeResp{Error: "unable to read body"})
		return
	}

	var req analyzeReq
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(analyzeResp{Error: "invalid json"})
		return
	}

	l := lexer.New(req.Code)
	toks := l.Tokenize()

	out := make([]tokenDTO, 0, len(toks))
	for _, t := range toks {
		// Skip NEWLINE tokens if desired - keep everything for now
		typeName := token.TokenNames[t.Type]
		if typeName == "" {
			typeName = "UNKNOWN"
		}
		out = append(out, tokenDTO{
			Lexeme: t.Lexeme,
			Type:   typeName,
			Line:   t.Line,
			Column: t.Column,
		})
	}

	resp := analyzeResp{Tokens: out}
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
