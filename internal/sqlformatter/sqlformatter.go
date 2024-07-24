// sqlformatter is based on the work https://github.com/cockroachdb/sqlfmt/blob/v0.4.0/main.go.
package sqlformatter

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	_ "github.com/cockroachdb/cockroach/pkg/sql/sem/builtins"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
)

var ignoreComments = regexp.MustCompile(`^--.*\s*`)

func fmtsql(cfg tree.PrettyCfg, stmts []string) (string, error) {
	var prettied strings.Builder
	for _, stmt := range stmts {
		for len(stmt) > 0 {
			stmt = strings.TrimSpace(stmt)
			hasContent := false
			// Trim comments, preserving whitespace after them.
			for {
				found := ignoreComments.FindString(stmt)
				if found == "" {
					break
				}
				// Remove trailing whitespace but keep up to 2 newlines.
				prettied.WriteString(strings.TrimRightFunc(found, unicode.IsSpace))
				newlines := strings.Count(found, "\n")
				if newlines > 2 {
					newlines = 2
				}
				prettied.WriteString(strings.Repeat("\n", newlines))
				stmt = stmt[len(found):]
				hasContent = true
			}
			// Split by semicolons
			next := stmt
			if pos, _ := parser.SplitFirstStatement(stmt); pos > 0 {
				next = stmt[:pos]
				stmt = stmt[pos:]
			} else {
				stmt = ""
			}
			// This should only return 0 or 1 responses.
			allParsed, err := parser.Parse(next)
			if err != nil {
				return "", err
			}
			for _, parsed := range allParsed {
				prettied.WriteString(cfg.Pretty(parsed.AST))
				prettied.WriteString(";\n")
				hasContent = true
			}
			if hasContent {
				prettied.WriteString("\n")
			}
		}
	}

	return strings.TrimRightFunc(prettied.String(), unicode.IsSpace), nil
}
