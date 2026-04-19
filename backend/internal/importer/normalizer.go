package importer

import "strings"

func IsCommentColumn(col string) bool {
	c := strings.ToLower(col)
	return strings.Contains(c, "comment") || strings.Contains(c, "recommendation")
}

func IsScoreableValue(v string) bool {
	n := strings.ToLower(strings.TrimSpace(v))
	return n == "yes" || n == "no" || n == "1" || n == "0" || n == "true" || n == "false"
}
