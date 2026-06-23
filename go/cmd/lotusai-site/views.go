package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

var botMarkers = []string{
	"bot", "crawl", "spider", "slurp", "bing", "google", "yandex",
	"facebookexternal", "semrush", "ahrefs", "petal", "gpt", "claude",
	"python-requests", "curl", "go-http", "libredtail", "wget", "scan",
}

func isBot(ua string) bool {
	l := strings.ToLower(ua)
	if l == "" {
		return true
	}
	for _, m := range botMarkers {
		if strings.Contains(l, m) {
			return true
		}
	}
	return false
}

type viewStore struct {
	path string
	mu   sync.Mutex
	data map[string]int
}

func newViewStore(path string) *viewStore {
	vs := &viewStore{path: path, data: map[string]int{}}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &vs.data)
	}
	return vs
}

func (vs *viewStore) hit(slug, ua string) {
	if isBot(ua) {
		return
	}
	vs.mu.Lock()
	vs.data[slug]++
	vs.mu.Unlock()
}

func (vs *viewStore) count(slug string) int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return vs.data[slug]
}

func (vs *viewStore) flush() error {
	vs.mu.Lock()
	b, err := json.MarshalIndent(vs.data, "", "  ")
	vs.mu.Unlock()
	if err != nil {
		return err
	}
	return os.WriteFile(vs.path, b, 0o644)
}

func (vs *viewStore) snapshot() map[string]int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	out := make(map[string]int, len(vs.data))
	for k, v := range vs.data {
		out[k] = v
	}
	return out
}
