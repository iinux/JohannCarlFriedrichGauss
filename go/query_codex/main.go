package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultEndpoint = "https://chatgpt.com/backend-api/wham/usage"

type authFile struct {
	Tokens struct {
		AccessToken string `json:"access_token"`
		AccountID   string `json:"account_id"`
	} `json:"tokens"`
}

type window struct {
	UsedPercent        float64 `json:"used_percent"`
	ResetAt            int64   `json:"reset_at"`
	LimitWindowSeconds int64   `json:"limit_window_seconds"`
}

type rateLimit struct {
	PrimaryWindow   *window `json:"primary_window"`
	SecondaryWindow *window `json:"secondary_window"`
}

type credits struct {
	HasCredits bool            `json:"has_credits"`
	Unlimited  bool            `json:"unlimited"`
	Balance    *flexibleNumber `json:"balance"`
}

type flexibleNumber float64

func (n *flexibleNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var value float64
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		text, err := strconv.Unquote(string(data))
		if err != nil {
			return err
		}
		value, err = strconv.ParseFloat(text, 64)
		if err != nil {
			return fmt.Errorf("无效的额度余额 %q: %w", text, err)
		}
	} else {
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
	}
	*n = flexibleNumber(value)
	return nil
}

type usage struct {
	PlanType  string    `json:"plan_type"`
	RateLimit rateLimit `json:"rate_limit"`
	Credits   credits   `json:"credits"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, http.DefaultClient))
}

func run(args []string, stdout, stderr io.Writer, client *http.Client) int {
	home, _ := os.UserHomeDir()
	if err := loadDotEnv(filepath.Join(home, ".codex", ".env")); err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(stderr, "读取 ~/.codex/.env 失败: %v\n", err)
		return 1
	}
	fs := flag.NewFlagSet("query-codex", flag.ContinueOnError)
	fs.SetOutput(stderr)
	authPath := fs.String("auth-file", filepath.Join(home, ".codex", "auth.json"), "Codex auth.json 路径")
	endpoint := fs.String("endpoint", envOr("CODEX_USAGE_ENDPOINT", defaultEndpoint), "额度接口地址")
	jsonOutput := fs.Bool("json", false, "原样输出 JSON")
	timeout := fs.Duration("timeout", 8*time.Second, "请求超时")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	auth, err := loadAuth(*authPath)
	if err != nil {
		fmt.Fprintf(stderr, "读取 Codex 登录信息失败: %v\n请先运行 codex login，或用 -auth-file 指定凭据文件。\n", err)
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	body, err := fetchUsage(ctx, client, *endpoint, auth)
	if err != nil {
		fmt.Fprintf(stderr, "查询失败: %v\n", err)
		return 1
	}
	if *jsonOutput {
		var pretty any
		if err := json.Unmarshal(body, &pretty); err != nil {
			fmt.Fprintf(stderr, "接口返回的不是合法 JSON: %v\n", err)
			return 1
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(pretty); err != nil {
			fmt.Fprintf(stderr, "输出失败: %v\n", err)
			return 1
		}
		return 0
	}

	var u usage
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Fprintf(stderr, "解析额度失败: %v\n可加 -json 查看原始响应。\n", err)
		return 1
	}
	printUsage(stdout, u, time.Now())
	return 0
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for lineNo := 1; scanner.Scan(); lineNo++ {
		key, value, ok, err := parseEnvLine(scanner.Text())
		if err != nil {
			return fmt.Errorf("%s:%d: %w", path, lineNo, err)
		}
		if ok {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("%s:%d: %w", path, lineNo, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 %s: %w", path, err)
	}
	return nil
}

func parseEnvLine(line string) (string, string, bool, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false, nil
	}
	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
	}
	key, value, found := strings.Cut(line, "=")
	key = strings.TrimSpace(key)
	if !found || !validEnvKey(key) {
		return "", "", false, errors.New("无效的环境变量定义")
	}
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '\'' && value[len(value)-1] == '\'') || (value[0] == '"' && value[len(value)-1] == '"')) {
		value = value[1 : len(value)-1]
	} else if i := strings.Index(value, " #"); i >= 0 {
		value = strings.TrimSpace(value[:i])
	}
	return key, value, true, nil
}

func validEnvKey(key string) bool {
	if key == "" || !(key[0] == '_' || key[0] >= 'A' && key[0] <= 'Z' || key[0] >= 'a' && key[0] <= 'z') {
		return false
	}
	for i := 1; i < len(key); i++ {
		c := key[i]
		if !(c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func loadAuth(path string) (authFile, error) {
	var auth authFile
	b, err := os.ReadFile(path)
	if err != nil {
		return auth, err
	}
	if err := json.Unmarshal(b, &auth); err != nil {
		return auth, fmt.Errorf("解析 %s: %w", path, err)
	}
	if strings.TrimSpace(auth.Tokens.AccessToken) == "" {
		return auth, errors.New("auth.json 中缺少 tokens.access_token")
	}
	return auth, nil
}

func fetchUsage(ctx context.Context, client *http.Client, endpoint string, auth authFile) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+auth.Tokens.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "query-codex/1.0")
	if auth.Tokens.AccountID != "" {
		req.Header.Set("ChatGPT-Account-Id", auth.Tokens.AccountID)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(string(body))
		if len(msg) > 300 {
			msg = msg[:300] + "..."
		}
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, msg)
	}
	return body, nil
}

func printUsage(w io.Writer, u usage, now time.Time) {
	if u.PlanType != "" {
		fmt.Fprintf(w, "Codex 套餐: %s\n", u.PlanType)
	}
	printWindow(w, "短周期", u.RateLimit.PrimaryWindow, now)
	printWindow(w, "长周期", u.RateLimit.SecondaryWindow, now)
	if u.Credits.Unlimited {
		fmt.Fprintln(w, "额外点数: 不限")
	} else if u.Credits.Balance != nil {
		fmt.Fprintf(w, "额外点数: %.2f\n", float64(*u.Credits.Balance))
	}
}

func printWindow(w io.Writer, name string, v *window, now time.Time) {
	if v == nil {
		return
	}
	remaining := 100 - v.UsedPercent
	if remaining < 0 {
		remaining = 0
	}
	label := name
	if v.LimitWindowSeconds > 0 {
		label = fmt.Sprintf("%s(%s)", name, compactDuration(time.Duration(v.LimitWindowSeconds)*time.Second))
	}
	fmt.Fprintf(w, "%s: 剩余 %.1f%%", label, remaining)
	if v.ResetAt > 0 {
		reset := time.Unix(v.ResetAt, 0)
		until := time.Until(reset)
		if !now.IsZero() {
			until = reset.Sub(now)
		}
		if until < 0 {
			until = 0
		}
		fmt.Fprintf(w, "，%s后重置 (%s)", compactDuration(until), reset.Local().Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintln(w)
}

func compactDuration(d time.Duration) string {
	if d < time.Minute {
		return "<1分钟"
	}
	d = d.Round(time.Minute)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	}
	return strings.Join(parts, "")
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
