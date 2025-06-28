package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/samnart1/golang/008parser/internal/parser"
	"github.com/samnart1/golang/008parser/pkg/types"
)

// HandleParse processes JSON parse requests
func (s *Server) HandleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req types.ParseRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = s.config.DefaultFormat
	}
	if req.Indent == 0 {
		req.Indent = s.config.DefaultIndent
	}

	p := parser.New(s.config)
	result, err := p.ParseString(req.Data, "request")
	if err != nil {
		s.writeErrorResponse(w, fmt.Sprintf("JSON parse error: %v", err), http.StatusBadRequest)
		return
	}

	formatter := parser.NewFormatter(s.config)
	options := types.OutputOptions{
		Format: req.Format,
		Indent: req.Indent,
		Colors: false,
	}

	formatted, err := formatter.Format(result.Data, options)
	if err != nil {
		s.writeErrorResponse(w, fmt.Sprintf("Format error: %v", err), http.StatusInternalServerError)
		return
	}

	response := types.ParseResponse{
		Success:   true,
		Data:      result.Data,
		Formatted: formatted,
		Metadata: &types.Metadata{
			Type:        result.Type,
			Size:        result.Size,
			KeyCount:    result.KeyCount,
			ArrayLength: result.ArrayLength,
		},
	}

	s.writeJSONResponse(w, response, http.StatusOK)
}

// HandleValidate processes JSON validation requests
func (s *Server) HandleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req types.ParseRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}

	strict := r.URL.Query().Get("strict") == "true"

	validator := parser.NewValidator(s.config)
	result := validator.ValidateString(req.Data, strict)

	response := types.ValidationResponse{
		Success:  true,
		Valid:    result.IsValid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}

	statusCode := http.StatusOK
	if !result.IsValid {
		statusCode = http.StatusBadRequest
	}

	s.writeJSONResponse(w, response, statusCode)
}

// HandleFormat processes JSON formatting requests
func (s *Server) HandleFormat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var req types.ParseRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeErrorResponse(w, "Invalid request JSON", http.StatusBadRequest)
		return
	}

	if req.Format == "" {
		req.Format = "pretty"
	}
	if req.Indent == 0 {
		req.Indent = s.config.DefaultIndent
	}

	p := parser.New(s.config)
	result, err := p.ParseString(req.Data, "request")
	if err != nil {
		s.writeErrorResponse(w, fmt.Sprintf("JSON parse error: %v", err), http.StatusBadRequest)
		return
	}

	formatter := parser.NewFormatter(s.config)
	options := types.OutputOptions{
		Format:   req.Format,
		Indent:   req.Indent,
		Colors:   false,
		SortKeys: r.URL.Query().Get("sort") == "true",
		UseTabs:  r.URL.Query().Get("tabs") == "true",
	}

	if maxDepthStr := r.URL.Query().Get("max_depth"); maxDepthStr != "" {
		if maxDepth, err := strconv.Atoi(maxDepthStr); err == nil {
			options.MaxDepth = maxDepth
		}
	}

	formatted, err := formatter.Format(result.Data, options)
	if err != nil {
		s.writeErrorResponse(w, fmt.Sprintf("Format error: %v", err), http.StatusInternalServerError)
		return
	}

	response := types.ParseResponse{
		Success:   true,
		Formatted: formatted,
	}

	s.writeJSONResponse(w, response, http.StatusOK)
}