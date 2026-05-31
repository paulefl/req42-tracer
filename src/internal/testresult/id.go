package testresult

import "strings"

// shortPackage extracts the last segment of a package path.
// Handles slash-separated Go paths (github.com/foo/bar → bar)
// and dot-separated Java classnames (com.example.MyTests → MyTests).
func shortPackage(pkg string) string {
	if i := strings.LastIndex(pkg, "/"); i >= 0 {
		return pkg[i+1:]
	}
	if i := strings.LastIndex(pkg, "."); i >= 0 {
		return pkg[i+1:]
	}
	return pkg
}
