package checker

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"github.com/cloverrose/pkgdep/pkg/cachedregexp"
	"github.com/cloverrose/pkgdep/pkg/orderedmap"
)

type Checker struct {
	dependencies orderedmap.OrderedMap
	recorder     recorder
	regexpCache  *cachedregexp.CachedRegexp
}

func New(dependencies orderedmap.OrderedMap, recorder recorder) *Checker {
	return &Checker{
		dependencies: dependencies,
		recorder:     recorder,
		regexpCache:  cachedregexp.New(true),
	}
}

func (c *Checker) IsAllowedDependency(from, to string) bool {
	for fromPattern, toTemplateStrings := range c.dependencies.Iter() {
		data, err := c.matchAndExtract(fromPattern, from)
		if err != nil {
			continue
		}
		for _, toTemplateString := range toTemplateStrings {
			toPattern, err := buildPattern(toTemplateString, data)
			if err != nil {
				continue
			}
			re, err := c.regexpCache.Compile(toPattern)
			if err != nil {
				continue
			}
			if re.MatchString(to) {
				c.recorder.RecordUsage(fromPattern, toTemplateString)
				return true
			}
		}
	}
	return false
}

func (c *Checker) matchAndExtract(pattern, text string) (map[string]string, error) {
	re, err := c.regexpCache.Compile(pattern)
	if err != nil {
		return nil, err
	}
	match := re.FindStringSubmatch(text)
	if match == nil {
		return nil, errors.New("no match")
	}
	names := re.SubexpNames()
	data := make(map[string]string)
	for i, name := range names {
		if name != "" {
			data[name] = match[i]
		}
	}
	return data, nil
}

func buildPattern(templateString string, data map[string]string) (string, error) {
	// By using missingkey=error, if templateString refers a key that is not present in the data.
	// https://pkg.go.dev/text/template#Template.Option
	tmpl, err := template.New("example").Option("missingkey=error").Parse(templateString)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	ret := result.String()
	if strings.Contains(ret, "<no value>") {
		// text/template missingkey=error option does not affect "index" https://github.com/golang/go/issues/60008
		// ret is used to match with go package and go package does not contain `<` and `>`.
		// We can ignore the possibility that user defines template with `<no value>`.
		// Returning error manually to align non index case.
		return "", errors.New("missing key")
	}

	return ret, nil
}
