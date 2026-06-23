package main

import (
	"path/filepath"
	"testing"
)

func TestIsBot(t *testing.T) {
	bots := []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1)",
		"python-requests/2.31",
		"curl/8.1",
		"Go-http-client/1.1",
	}
	for _, ua := range bots {
		if !isBot(ua) {
			t.Errorf("isBot(%q) = false, want true", ua)
		}
	}
	human := "Mozilla/5.0 (Linux; Android 14; Pixel 8) Chrome/135 Mobile Safari/537.36"
	if isBot(human) {
		t.Errorf("isBot(human) = true, want false")
	}
}

func TestViewStoreIncrementAndPersist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "views.json")
	vs := newViewStore(path)
	vs.hit("2026-06-23-event-mover-hpg", "Chrome/135 Mobile")
	vs.hit("2026-06-23-event-mover-hpg", "Chrome/135 Mobile")
	vs.hit("2026-06-23-event-mover-hpg", "Googlebot/2.1") // bot — ignored
	if err := vs.flush(); err != nil {
		t.Fatal(err)
	}
	reloaded := newViewStore(path)
	if got := reloaded.count("2026-06-23-event-mover-hpg"); got != 2 {
		t.Fatalf("count = %d, want 2 (bot excluded, persisted)", got)
	}
}
