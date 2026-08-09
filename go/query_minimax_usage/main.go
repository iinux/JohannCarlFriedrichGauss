// 查询 MiniMax Token Plan 订阅剩余额度
//
// 用法:
//
//	export MINIMAX_API_KEY=sk-cp-xxx
//	go run .                 # 自动检测 key 类型并尝试两个平台, 显示进度条
//	go run . -raw            # 只打印原始 JSON
//	go run . -k sk-cp-xxx    # 通过命令行传 key
//	go run . -platform cn    # 强制用中国站 (默认)
//	go run . -platform intl  # 强制用国际站
//	go run . -platform auto  # 先中国站, 再国际站
//	go run . -no-color       # 关闭颜色
//
// Token Plan 使用 Subscription Key（管理位置：控制台 Billing > Token Plan）,
// 而非 Pay-as-you-go API Key. 两平台 key 不互通.
//
// 端点:
//
//	GET {base}/v1/token_plan/remains
//	Authorization: Bearer <key>
//	Content-Type: application/json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultIntlEndpoint = "https://api.minimax.io/v1/token_plan/remains"
	defaultCNEndpoint   = "https://api.minimaxi.com/v1/token_plan/remains"

	barWidth = 30 // 进度条字符宽度
)

// 顶层响应
type response struct {
	Code        any            `json:"code"`
	Msg         string         `json:"msg"`
	Data        map[string]any `json:"data"`
	BaseResp    *baseResp      `json:"base_resp"`
	ModelRemains []modelRemain `json:"model_remains"`
}

type baseResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type modelRemain struct {
	ModelName                       string `json:"model_name"`
	CurrentIntervalRemainingPercent int    `json:"current_interval_remaining_percent"`
	CurrentIntervalStatus           int    `json:"current_interval_status"`
	CurrentIntervalTotalCount       int64  `json:"current_interval_total_count"`
	CurrentIntervalUsageCount       int64  `json:"current_interval_usage_count"`
	CurrentWeeklyRemainingPercent   int    `json:"current_weekly_remaining_percent"`
	CurrentWeeklyStatus             int    `json:"current_weekly_status"`
	CurrentWeeklyTotalCount         int64  `json:"current_weekly_total_count"`
	CurrentWeeklyUsageCount         int64  `json:"current_weekly_usage_count"`
	EndTime                         int64  `json:"end_time"`
	StartTime                       int64  `json:"start_time"`
	RemainsTime                     int64  `json:"remains_time"`
	WeeklyBoostPermille             int    `json:"weekly_boost_permille"`
	WeeklyEndTime                   int64  `json:"weekly_end_time"`
	WeeklyRemainsTime               int64  `json:"weekly_remains_time"`
	WeeklyStartTime                 int64  `json:"weekly_start_time"`
}

type platform struct {
	name     string
	endpoint string
}

var noColor bool

func main() {
	apiKey := flag.String("k", os.Getenv("MINIMAX_API_KEY"), "API/Subscription key（也可通过 MINIMAX_API_KEY 设置）")
	platformFlag := flag.String("platform", "cn", "使用哪个平台: cn (默认) | intl | auto (cn→intl)")
	raw := flag.Bool("raw", false, "只打印原始 JSON")
	timeout := flag.Duration("timeout", 30*time.Second, "HTTP 超时时间")
	flag.BoolVar(&noColor, "no-color", false, "关闭颜色输出")
	flag.Parse()

	key := strings.TrimSpace(*apiKey)
	if key == "" {
		fmt.Fprintln(os.Stderr, "错误: 需要 key，通过 -k 或 MINIMAX_API_KEY 环境变量设置")
		flag.Usage()
		os.Exit(2)
	}

	fmt.Fprintf(os.Stderr, "key 类型: %s\n", detectKeyType(key))

	platforms := pickPlatforms(*platformFlag)
	fmt.Fprintf(os.Stderr, "尝试: %s\n\n", joinNames(platforms))

	var lastErr error
	var lastBody []byte
	for _, p := range platforms {
		fmt.Fprintf(os.Stderr, "→ %s  %s\n", p.name, p.endpoint)
		body, err := fetch(p.endpoint, key, *timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ %v\n\n", err)
			lastErr = err
			lastBody = nil
			continue
		}

		var resp response
		_ = json.Unmarshal(body, &resp)
		if resp.BaseResp != nil && resp.BaseResp.StatusCode != 0 {
			fmt.Fprintf(os.Stderr, "  ✗ %d %s\n\n", resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
			lastErr = fmt.Errorf("%s: %d %s", p.name, resp.BaseResp.StatusCode, resp.BaseResp.StatusMsg)
			lastBody = body
			continue
		}
		if v, ok := resp.Code.(float64); ok && v != 0 {
			fmt.Fprintf(os.Stderr, "  ✗ code=%v %s\n\n", v, resp.Msg)
			lastErr = fmt.Errorf("%s: code=%v %s", p.name, v, resp.Msg)
			lastBody = body
			continue
		}

		fmt.Fprintf(os.Stderr, "  ✓ 成功\n\n")
		output(body, *raw, resp)
		return
	}

	// 全部失败
	fmt.Fprintln(os.Stderr, "所有平台均失败。可能原因:")
	fmt.Fprintln(os.Stderr, "  1. Token Plan 需要 Subscription Key (sk-cp- 前缀), 不是 Pay-as-you-go API Key")
	fmt.Fprintln(os.Stderr, "     获取位置: 控制台 → Billing → Token Plan")
	fmt.Fprintln(os.Stderr, "  2. key 已过期、被撤销, 或订阅已停止")
	fmt.Fprintln(os.Stderr, "  3. 当前是默认中国站; 国际站 key 用 -platform intl")
	if lastBody != nil {
		fmt.Fprintf(os.Stderr, "\n最后一次原始响应:\n%s\n", string(lastBody))
	}
	if lastErr != nil {
		fmt.Fprintf(os.Stderr, "\n最后一次错误: %v\n", lastErr)
	}
	os.Exit(1)
}

func detectKeyType(key string) string {
	switch {
	case strings.HasPrefix(key, "sk-cp-"):
		return "Subscription Key (Token Plan, sk-cp- 前缀) ✓"
	case strings.HasPrefix(key, "sk-"):
		return "Pay-as-you-go API Key (sk- 前缀) — Token Plan 通常需要 Subscription Key"
	case strings.HasPrefix(key, "eyJ"):
		return "JWT token"
	default:
		return "未知格式"
	}
}

func pickPlatforms(p string) []platform {
	switch p {
	case "intl":
		return []platform{{"intl", defaultIntlEndpoint}}
	case "cn":
		return []platform{{"cn", defaultCNEndpoint}}
	case "auto":
		return []platform{
			{"cn", defaultCNEndpoint},
			{"intl", defaultIntlEndpoint},
		}
	default:
		return []platform{{"cn", defaultCNEndpoint}}
	}
}

func joinNames(ps []platform) string {
	names := make([]string, len(ps))
	for i, p := range ps {
		names[i] = p.name
	}
	return strings.Join(names, " → ")
}

func fetch(endpoint, apiKey string, timeout time.Duration) ([]byte, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

// output 主输出函数
func output(body []byte, raw bool, resp response) {
	pretty, err := prettyJSON(body)
	if err != nil {
		fmt.Println(string(body))
		return
	}

	// raw 模式: 只打印原始 JSON
	if raw {
		fmt.Println(string(body))
		return
	}

	fmt.Println(pretty)

	// 优先用 model_remains 数组渲染进度条
	if len(resp.ModelRemains) > 0 {
		printModelRemains(resp.ModelRemains)
		return
	}

	// 兜底: 尝试从 data 解析
	var data map[string]any
	if err := json.Unmarshal(body, &data); err == nil {
		printSummary(data)
	}
}

func prettyJSON(b []byte) (string, error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return "", err
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// statusName 把 status 字段翻译成可读文本
func statusName(s int) string {
	switch s {
	case 0:
		return "未启用"
	case 1:
		return "正常"
	case 2:
		return "受限"
	case 3:
		return "未订阅"
	default:
		return fmt.Sprintf("status=%d", s)
	}
}

// printModelRemains 按模型分组打印 5h + 周窗口进度条
func printModelRemains(items []modelRemain) {
	fmt.Println("\n═══ 订阅额度 ═══")
	for _, m := range items {
		label := modelDisplayName(m.ModelName)
		fmt.Printf("\n▶ %s  [%s]\n", label, statusName(m.CurrentIntervalStatus))

		// 5 小时窗口
		if m.CurrentIntervalRemainingPercent > 0 || m.CurrentIntervalStatus != 0 {
			fmt.Println("  5 小时窗口:")
			printBar(m.CurrentIntervalRemainingPercent,
				fmt.Sprintf("%d%% 剩余", m.CurrentIntervalRemainingPercent))
			printResetLine("重置倒计时", m.RemainsTime, m.EndTime)
		}

		// 周窗口
		if m.CurrentWeeklyRemainingPercent > 0 || m.CurrentWeeklyStatus != 0 {
			fmt.Println("  周窗口:")
			printBar(m.CurrentWeeklyRemainingPercent,
				fmt.Sprintf("%d%% 剩余", m.CurrentWeeklyRemainingPercent))
			if m.WeeklyBoostPermille > 0 {
				fmt.Printf("    周加成: +%g%%\n", float64(m.WeeklyBoostPermille)/10.0)
			}
			printResetLine("周重置倒计时", m.WeeklyRemainsTime, m.WeeklyEndTime)
		}
	}
	fmt.Println()
}

// modelDisplayName 把 internal model_name 翻译为友好名
func modelDisplayName(name string) string {
	switch strings.ToLower(name) {
	case "general":
		return "通用文本 (general)"
	case "video":
		return "视频生成 (video)"
	case "image":
		return "图像生成 (image)"
	case "speech", "audio", "voice":
		return "语音 (speech)"
	case "music":
		return "音乐 (music)"
	default:
		if name == "" {
			return "(未命名模型)"
		}
		return name
	}
}

// printBar 渲染一个进度条: [████████░░░░░░░░░░░░] 75%
func printBar(percent int, label string) {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * barWidth / 100
	empty := barWidth - filled

	color := pickColor(percent)
	reset := ""
	if !noColor {
		reset = "\033[0m"
		c := "\033[32m" // green
		switch color {
		case "yellow":
			c = "\033[33m"
		case "red":
			c = "\033[31m"
		}
		color = c
	} else {
		color = ""
	}

	bar := color + strings.Repeat("█", filled) + reset + strings.Repeat("░", empty)
	fmt.Printf("    [%s] %s\n", bar, label)
}

// pickColor 根据剩余百分比返回颜色档位
// 本地测试示例响应 (用 -sample 触发)
func renderSample() {
	noColor = true
	data, _ := os.ReadFile("/tmp/sample.json")
	var resp response
	_ = json.Unmarshal(data, &resp)
	printModelRemains(resp.ModelRemains)
}

func pickColor(percent int) string {
	switch {
	case percent >= 50:
		return "green"
	case percent >= 20:
		return "yellow"
	default:
		return "red"
	}
}

// printResetLine 打印一行重置倒计时
// 自动探测 remainsSec 单位: 若 > 1e8 则视为毫秒, 否则视为秒.
// 优先使用 endMs 计算实际倒计时.
func printResetLine(label string, remainsSec, endMs int64) {
	now := time.Now().UnixMilli()
	parts := []string{}
	if endMs > 0 {
		diffMs := endMs - now
		if diffMs > 0 {
			parts = append(parts, "剩余 "+formatDuration(time.Duration(diffMs)*time.Millisecond))
		} else {
			parts = append(parts, "已到重置时间")
		}
		parts = append(parts, "时间 "+formatTime(endMs))
	} else if remainsSec > 0 {
		// 探测单位: 14_169_864 显然是毫秒(≈4小时), 不是秒(164天)
		var d time.Duration
		if remainsSec > 100_000_000 {
			d = time.Duration(remainsSec) * time.Millisecond
		} else {
			d = time.Duration(remainsSec) * time.Second
		}
		parts = append(parts, "剩余 "+formatDuration(d))
	}
	if len(parts) == 0 {
		return
	}
	fmt.Printf("    %s: %s\n", label, strings.Join(parts, "  "))
}

// formatDuration 把秒数格式化为 d/h/m/s
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "已过期"
	}
	days := int(d / (24 * time.Hour))
	h := int((d % (24 * time.Hour)) / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	s := int((d % time.Minute) / time.Second)
	if days > 0 {
		return fmt.Sprintf("%d天%d时%d分", days, h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%d时%d分%d秒", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%d分%d秒", m, s)
	}
	return fmt.Sprintf("%d秒", s)
}

// formatTime 把毫秒时间戳格式化为本地时间字符串
func formatTime(ms int64) string {
	if ms <= 0 {
		return "-"
	}
	t := time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond))
	return t.Local().Format("2006-01-02 15:04:05")
}

// formatNum 把 JSON 数字格式化为不强制带 .0 的字符串
func formatNum(v any) string {
	switch x := v.(type) {
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// printSummary 兜底打印: 从任意 map 抽取已知字段
func printSummary(data map[string]any) {
	if len(data) == 0 {
		return
	}
	fmt.Println("\n── 摘要 ──")

	if v, ok := data["current_package_name"]; ok {
		fmt.Printf("套餐:      %v\n", v)
	}
	if v, ok := data["plan_name"]; ok {
		fmt.Printf("套餐:      %v\n", v)
	}

	fields := []struct {
		key string
		tag string
	}{
		{"remains", "剩余"},
		{"remaining", "剩余"},
		{"quota_remains", "剩余"},
		{"limit", "上限"},
		{"quota_limit", "上限"},
		{"used", "已用"},
		{"quota_used", "已用"},
	}
	for _, f := range fields {
		if v, ok := data[f.key]; ok {
			fmt.Printf("%s:  %v\n", f.tag, formatNum(v))
		}
	}

	if v, ok := data["five_hour"]; ok {
		fmt.Println("\n[5 小时窗口]")
		printQuotaBlock(v)
	}
	if v, ok := data["weekly"]; ok {
		fmt.Println("\n[周窗口]")
		printQuotaBlock(v)
	}
}

func printQuotaBlock(v any) {
	m, ok := v.(map[string]any)
	if !ok {
		fmt.Printf("  %v\n", v)
		return
	}
	if val, ok := m["limit"]; ok {
		fmt.Printf("  上限: %v\n", formatNum(val))
	}
	if val, ok := m["used"]; ok {
		fmt.Printf("  已用: %v\n", formatNum(val))
	}
	if val, ok := m["remaining"]; ok {
		fmt.Printf("  剩余: %v\n", formatNum(val))
	}
	if val, ok := m["reset_seconds"]; ok {
		if secs, ok := val.(float64); ok {
			fmt.Printf("  距离重置: %s\n", formatDuration(time.Duration(secs)*time.Second))
		} else {
			fmt.Printf("  距离重置: %v\n", val)
		}
	}
}
