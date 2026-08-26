package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFetchUsage(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Errorf("Authorization = %q", got)
		}
		if got := r.Header.Get("ChatGPT-Account-Id"); got != "account" {
			t.Errorf("account header = %q", got)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"plan_type":"plus"}`)),
		}, nil
	})}
	var auth authFile
	auth.Tokens.AccessToken = "secret"
	auth.Tokens.AccountID = "account"
	b, err := fetchUsage(context.Background(), client, "https://example.test/usage", auth)
	if err != nil || !strings.Contains(string(b), "plus") {
		t.Fatalf("body=%s err=%v", b, err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestPrintUsage(t *testing.T) {
	reset := time.Date(2026, 8, 21, 13, 0, 0, 0, time.Local)
	u := usage{PlanType: "plus", RateLimit: rateLimit{PrimaryWindow: &window{
		UsedPercent: 35, ResetAt: reset.Unix(), LimitWindowSeconds: 5 * 60 * 60,
	}}}
	var out bytes.Buffer
	printUsage(&out, u, reset.Add(-90*time.Minute))
	got := out.String()
	for _, want := range []string{"Codex 套餐: plus", "剩余 65.0%", "1小时30分钟后重置"} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestUsageAcceptsStringCreditBalance(t *testing.T) {
	var u usage
	if err := json.Unmarshal([]byte(`{"credits":{"balance":"0","has_credits":false}}`), &u); err != nil {
		t.Fatal(err)
	}
	if u.Credits.Balance == nil || float64(*u.Credits.Balance) != 0 {
		t.Fatalf("balance = %v, want 0", u.Credits.Balance)
	}
}

func TestUsageAcceptsNumericCreditBalance(t *testing.T) {
	var u usage
	if err := json.Unmarshal([]byte(`{"credits":{"balance":12.5}}`), &u); err != nil {
		t.Fatal(err)
	}
	if u.Credits.Balance == nil || float64(*u.Credits.Balance) != 12.5 {
		t.Fatalf("balance = %v, want 12.5", u.Credits.Balance)
	}
}

func TestCompactDuration(t *testing.T) {
	if got := compactDuration(50 * time.Hour); got != "2天2小时" {
		t.Fatalf("got %q", got)
	}
}

func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("# comment\nPLAIN=value\nexport QUOTED=\"hello world\"\nCOMMENTED=x # note\nEMPTY=\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"PLAIN", "QUOTED", "COMMENTED", "EMPTY"} {
		t.Setenv(key, "before")
	}
	if err := loadDotEnv(path); err != nil {
		t.Fatal(err)
	}
	wants := map[string]string{"PLAIN": "value", "QUOTED": "hello world", "COMMENTED": "x", "EMPTY": ""}
	for key, want := range wants {
		if got := os.Getenv(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}

func TestParseEnvLineRejectsInvalidKey(t *testing.T) {
	if _, _, _, err := parseEnvLine("BAD-KEY=value"); err == nil {
		t.Fatal("expected invalid key error")
	}
}
