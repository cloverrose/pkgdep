package cachedregexp

import (
	"regexp"
	"sync"
)

// CachedRegexp is a wrapper around regexp.Compile that caches the compiled regexps.
type CachedRegexp struct {
	enabledCache bool
	data         map[string]*regexp.Regexp
	mu           sync.RWMutex
}

// New creates a new CachedRegexp.
func New(enabledCache bool) *CachedRegexp {
	return &CachedRegexp{
		enabledCache: enabledCache,
		data:         make(map[string]*regexp.Regexp, 64),
	}
}

// Compile compiles a regular expression and returns, as a compiled regular expression, the expression and an error, if any.
func (c *CachedRegexp) Compile(pattern string) (*regexp.Regexp, error) {
	if !c.enabledCache {
		return regexp.Compile(pattern)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if re, ok := c.data[pattern]; ok {
		return re, nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	c.data[pattern] = re
	return re, nil
}
