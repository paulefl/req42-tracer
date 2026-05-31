package testresult

import "testing"

// [test-spec,id=TS-TR-025,req="REQ-TESTING-001",aspice="SWE.5.BP3"]
// TestShortPackage verifies extraction of last segment from slash- and dot-separated package paths.
func TestShortPackage(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"github.com/paulefl/req42-tracer/src/internal/parser", "parser"},
		{"github.com/paulefl/req42-tracer/src/internal/graph", "graph"},
		{"com.example.MyTests", "MyTests"},
		{"pkg", "pkg"},
		{"", ""},
	}
	for _, c := range cases {
		got := shortPackage(c.input)
		if got != c.want {
			t.Errorf("shortPackage(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
