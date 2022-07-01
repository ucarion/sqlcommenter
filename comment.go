package sqlcommenter

import (
	"net/url"
	"strings"
)

// The sqlcommenter spec specifies that keys and values should be URL-encoded,
// and then the meta-characters "\" and "'" are to be backslash-escaped. This is
// so that it's unambiguous how to parse values, which are surrounded in "'".
//
// A wrinkle is that net/url will always percent-encode "\" and "'", so those
// characters will never appear in url.PathEscape's return value. The decoding
// specification requires backslash-unescaping followed by URL decoding, so it's
// fine if URL encoding is "overzealous".

func encodeKey(s string) string {
	return url.PathEscape(s)
}

func encodeValue(s string) string {
	return "'" + url.PathEscape(s) + "'"
}

func appendComment(attrs attrs, query string) string {
	// do not append empty comments
	if len(attrs) == 0 {
		return query
	}

	// do not modify queries containing comments
	if strings.Contains(query, "--") || strings.Contains(query, "/*") {
		return query
	}

	var b strings.Builder
	b.WriteString(query)
	b.WriteString(" /*")
	for i, kv := range attrs {
		if i != 0 {
			b.WriteString(",")
		}

		b.WriteString(kv[0])
		b.WriteString("=")
		b.WriteString(kv[1])
	}

	b.WriteString("*/")
	return b.String()
}
