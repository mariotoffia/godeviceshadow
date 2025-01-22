package reutils

import (
	"regexp"
	"sync"
)

var Shared = &RegexCache{}

// RegexCache is a cache for compiled regex patterns.
type RegexCache struct {
	cache sync.Map
}

func (rc *RegexCache) GetOrCompile(pattern string) (*regexp.Regexp, error) {
	if v, ok := rc.cache.Load(pattern); ok {
		return v.(*regexp.Regexp), nil
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	rc.cache.Store(pattern, re)

	return re, nil
}

func (rc *RegexCache) MustCompile(pattern string) *regexp.Regexp {
	re, err := rc.GetOrCompile(pattern)
	if err != nil {
		panic(err)
	}

	return re
}

func (rc *RegexCache) Get(pattern string) (*regexp.Regexp, bool) {
	if v, ok := rc.cache.Load(pattern); ok {
		return v.(*regexp.Regexp), true
	}

	return nil, false
}

func (rc *RegexCache) Delete(pattern string) {
	rc.cache.Delete(pattern)
}
