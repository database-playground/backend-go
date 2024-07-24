package dbrunner

import "github.com/DataDog/go-sqllexer"

// FormatSQL formats the raw SQL string to the normalized form.
func FormatSQL(raw string) (string, error) {
	normalizer := sqllexer.NewNormalizer(
		sqllexer.WithCollectComments(false),
		sqllexer.WithCollectCommands(true),
		sqllexer.WithCollectTables(true),
		sqllexer.WithKeepSQLAlias(false),
		sqllexer.WithRemoveSpaceBetweenParentheses(true),
		sqllexer.WithUppercaseKeywords(true),
		sqllexer.WithKeepTrailingSemicolon(false),
	)
	normalized, _, err := normalizer.Normalize(raw)
	if err != nil {
		return "", err
	}

	return normalized, nil
}
