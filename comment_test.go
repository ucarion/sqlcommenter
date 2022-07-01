package sqlcommenter

import "testing"

func TestEncode(t *testing.T) {
	testCases := []struct {
		s, k, v string
	}{
		{"", "", "''"},
		{"a", "a", "'a'"},
		{"/", "%2F", "'%2F'"},
		{"'", "%27", "'%27'"},
		{"\\", "%5C", "'%5C'"},
		{"a/'/\\/a", "a%2F%27%2F%5C%2Fa", "'a%2F%27%2F%5C%2Fa'"},

		// sql injection risks
		{"/* */", "%2F%2A%20%2A%2F", "'%2F%2A%20%2A%2F'"},
	}

	for _, tt := range testCases {
		t.Run(tt.s, func(t *testing.T) {
			if encodeKey(tt.s) != tt.k {
				t.Fatalf("bad encode key, want: %q, got: %q", tt.k, encodeKey(tt.s))
			}

			if encodeValue(tt.s) != tt.v {
				t.Fatalf("bad encode val, want: %q, got: %q", tt.v, encodeValue(tt.s))
			}
		})
	}
}

func TestAppendComment(t *testing.T) {
	testCases := []struct {
		want  string
		query string
		attrs attrs
	}{
		{
			want:  "select 1 /*a='b',c='d'*/",
			query: "select 1",
			attrs: attrs{{"a", "'b'"}, {"c", "'d'"}},
		},
		{
			want:  "select 1",
			query: "select 1",
			attrs: attrs{},
		},
		{
			want:  "select 1 -- a",
			query: "select 1 -- a",
			attrs: attrs{{"a", "b"}},
		},
		{
			want:  "select /* a */ 1",
			query: "select /* a */ 1",
			attrs: attrs{{"a", "b"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.want, func(t *testing.T) {
			if appendComment(tt.attrs, tt.query) != tt.want {
				t.Fatalf("bad appendComment, want: %q, got: %q", tt.want, appendComment(tt.attrs, tt.query))
			}
		})
	}
}
