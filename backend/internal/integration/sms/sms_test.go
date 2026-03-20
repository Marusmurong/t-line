package sms

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSendCode_NoFullCodeInSource verifies that the SendCode function
// never logs the raw verification code. We parse the actual source code
// to ensure no fmt.Printf / log line contains the raw `code` variable
// in a way that would leak it.
func TestSendCode_NoFullCodeInSource(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "sms.go", nil, parser.AllErrors)
	require.NoError(t, err, "should be able to parse sms.go")

	// Walk the AST for the SendCode function body
	var sendCodeFunc *ast.FuncDecl
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "SendCode" {
			sendCodeFunc = fn
			break
		}
	}
	require.NotNil(t, sendCodeFunc, "SendCode function must exist in sms.go")

	// Collect all string literals in SendCode to check for code leakage patterns
	var stringLiterals []string
	ast.Inspect(sendCodeFunc.Body, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			stringLiterals = append(stringLiterals, lit.Value)
		}
		return true
	})

	// Verify no format string contains the full code (e.g., "code=%s" or printing code directly)
	for _, s := range stringLiterals {
		lower := strings.ToLower(s)
		assert.NotContains(t, lower, "code=",
			"SendCode should not contain format strings that log the raw code")
		assert.NotContains(t, lower, "code: %s",
			"SendCode should not log the raw verification code")
		assert.NotContains(t, lower, "code: %d",
			"SendCode should not log the raw verification code")
	}
}

// TestSendCode_PhoneMaskedInDebugLog verifies that the debug log format string
// masks the phone number (shows only first 3 and last 4 characters).
func TestSendCode_PhoneMaskedInDebugLog(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "sms.go", nil, parser.AllErrors)
	require.NoError(t, err)

	var sendCodeFunc *ast.FuncDecl
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "SendCode" {
			sendCodeFunc = fn
			break
		}
	}
	require.NotNil(t, sendCodeFunc)

	// Check that any Debugf call uses masked phone format (phone[:3] ... phone[len-4:])
	// We verify the format string contains "***" which indicates masking
	found := false
	ast.Inspect(sendCodeFunc.Body, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING && strings.Contains(lit.Value, "***") {
			found = true
		}
		return true
	})

	assert.True(t, found, "SendCode debug log must mask the phone number with '***'")
}

// TestSendCode_DebugModeGuarded verifies that the log statement is guarded
// by APP_MODE == "debug" check.
func TestSendCode_DebugModeGuarded(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "sms.go", nil, parser.AllErrors)
	require.NoError(t, err)

	var sendCodeFunc *ast.FuncDecl
	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if ok && fn.Name.Name == "SendCode" {
			sendCodeFunc = fn
			break
		}
	}
	require.NotNil(t, sendCodeFunc)

	// Check for an if-statement containing "debug" guard
	hasDebugGuard := false
	ast.Inspect(sendCodeFunc.Body, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}
		// Check if the condition contains a comparison with "debug"
		ast.Inspect(ifStmt.Cond, func(inner ast.Node) bool {
			lit, ok := inner.(*ast.BasicLit)
			if ok && lit.Kind == token.STRING && strings.Contains(lit.Value, "debug") {
				hasDebugGuard = true
			}
			return true
		})
		return true
	})

	assert.True(t, hasDebugGuard, "SendCode log must be guarded by APP_MODE == 'debug' check")
}

// TestGenerateCode_SixDigits verifies the code is always 6 digits.
func TestGenerateCode_SixDigits(t *testing.T) {
	for i := 0; i < 100; i++ {
		code := generateCode()
		assert.Len(t, code, 6, "verification code must be exactly 6 digits")
	}
}
