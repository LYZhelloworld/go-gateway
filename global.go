package gateway

import (
	"regexp"
	"strings"
)

const (
	baseServiceHandler = "*"
)

var pathRegexp = regexp.MustCompile("^(?:/|(?:/(?:[A-Za-z0-9-._~]|%[0-9A-Fa-f]{2})+)+)$")
var serviceRegexp = regexp.MustCompile("^[a-zA-Z0-9_-]+(?:\\.[a-zA-Z0-9_-]+)*$")

// isValidPath checks if the path is a valid path, or a valid prefix which has a path and a suffix "/*".
func isValidPath(path string) bool {
	return pathRegexp.MatchString(path)
}

// trimPrefix trims the "/*" at the end of the prefix and returns the path.
// If the path does not have "/*", it remains unchanged.
func trimPrefix(path string) string {
	if path == "/*" {
		return "/"
	}
	return strings.TrimSuffix(path, "/*")
}

// removeLastDir removes last directory of a path. For example: removeLastDir("/foo/bar/baz") gives "/foo/bar".
// Removing "/" gives empty string, but removing "/foo" will produce "/".
func removeLastDir(path string) string {
	if path == "/" {
		return ""
	}
	r := path[:strings.LastIndex(path, "/")]
	if r == "" {
		r = "/"
	}
	return r
}

// isValidService checks if the Service name is valid.
func isValidService(serviceName string) bool {
	return serviceName == baseServiceHandler || serviceRegexp.MatchString(serviceName)
}

// removeLastSubService removes last sub-Service of a Service name.
// For example: removeLastSubService("foo.bar.baz") gives "foo.bar".
// removeLastSubService("foo") gives empty string.
func removeLastSubService(serviceName string) string {
	idx := strings.LastIndex(serviceName, ".")
	if idx == -1 {
		return ""
	}
	return serviceName[:idx]
}
