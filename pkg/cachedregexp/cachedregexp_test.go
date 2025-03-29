package cachedregexp

import (
	"testing"
)

func Benchmark_CachedRegexp_Compile_disable(b *testing.B) {
	re := New(false)

	patterns := providePatterns()

	for b.Loop() {
		for _, pattern := range patterns {
			_, _ = re.Compile(pattern)
		}
	}
}

func Benchmark_CachedRegexp_Compile_enable(b *testing.B) {
	re := New(true)

	patterns := providePatterns()

	for b.Loop() {
		for _, pattern := range patterns {
			_, _ = re.Compile(pattern)
		}
	}
}

func providePatterns() []string {
	return []string{
		// basic patterns
		"hello",  // simple pattern
		"^foo$",  // exact match: "foo" only
		"bar.*",  // prefix match: "bar" starts with
		".*baz$", // suffix match: "baz" ends with

		// intermediate patterns
		"^[A-Z][a-z]+$",    // word starting with uppercase
		"\\d{3}-\\d{4}",    // postal code: 123-4567
		"[a-z]+(_[a-z]+)*", // snake case variable name

		// advanced patterns
		"(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})", // named capture group: year-month-day
		"^(?:https?://)?[\\w.-]+\\.[a-z]{2,}$",               // simple URL validation
		"^(?=.*[A-Z])(?=.*[a-z])(?=.*\\d).{8,}$",             // password requirements: uppercase, lowercase, digit, 8 characters or more

		// special cases
		"\\\\.",       // escaped dot
		"[^\\s]+",     // non-space characters
		"(?i)pattern", // case-insensitive
	}
}
