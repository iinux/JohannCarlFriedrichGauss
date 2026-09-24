// go run ping_report.go [list_file] [count]
  package main

  import (
      "bufio"
      "fmt"
      "math"
      "os"
      "os/exec"
      "regexp"
      "sort"
      "strconv"
      "strings"
      "time"
  )

  type sample struct {
      target string
      ok     bool
      lat    float64
  }

  type summary struct {
      Target string
      Sent   int
      Recv   int
      Loss   float64
      Min    float64
      Avg    float64
      Max    float64
      Mdev   float64
  }

  var timeRE = regexp.MustCompile(`time[=<]\s*([\d.]+)\s*ms`)

  func pingOnce(target string, timeout int) (bool, float64) {
      cmd := exec.Command("ping", "-c", "1", "-W", strconv.Itoa(timeout), target)
      out, err := cmd.Output()
      if err != nil {
          return false, 0
      }
      m := timeRE.FindStringSubmatch(string(out))
      if len(m) < 2 {
          return false, 0
      }
      lat, err := strconv.ParseFloat(m[1], 64)
      if err != nil {
          return false, 0
      }
      return true, lat
  }

  func probe(target string, count, timeout int) []sample {
      samples := make([]sample, 0, count)
      for i := 0; i < count; i++ {
          ok, lat := pingOnce(target, timeout)
          samples = append(samples, sample{target, ok, lat})
      }
      return samples
  }

  func summarize(samples []sample) summary {
      var s summary
      s.Target = samples[0].target
      s.Sent = len(samples)
      var vals []float64
      for _, x := range samples {
          if x.ok {
              s.Recv++
              vals = append(vals, x.lat)
          }
      }
      s.Loss = float64(s.Sent-s.Recv) / float64(s.Sent) * 100
      if len(vals) > 0 {
          s.Min, s.Max = vals[0], vals[0]
          var sum float64
          for _, v := range vals {
              sum += v
              if v < s.Min {
                  s.Min = v
              }
              if v > s.Max {
                  s.Max = v
              }
          }
          s.Avg = sum / float64(len(vals))
          if len(vals) > 1 {
              var sq float64
              for _, v := range vals {
                  sq += (v - s.Avg) * (v - s.Avg)
              }
              s.Mdev = math.Sqrt(sq / float64(len(vals)-1))
          }
      }
      return s
  }

  func loadTargets(path string) ([]string, error) {
      f, err := os.Open(path)
      if err != nil {
          return nil, err
      }
      defer f.Close()
      var out []string
      sc := bufio.NewScanner(f)
      for sc.Scan() {
          t := strings.TrimSpace(sc.Text())
          if t == "" || strings.HasPrefix(t, "#") {
              continue
          }
          out = append(out, t)
      }
      return out, sc.Err()
  }

  func fmtOrDash(v float64) string {
      if v == 0 {
          return "-"
      }
      return strconv.FormatFloat(v, 'f', 2, 64)
  }

  func main() {
      listFile := "test_list"
      count := 4
      timeout := 2
      if len(os.Args) > 1 {
          listFile = os.Args[1]
      }
      if len(os.Args) > 2 {
          if n, err := strconv.Atoi(os.Args[2]); err == nil {
              count = n
          }
      }

      targets, err := loadTargets(listFile)
      if err != nil {
          fmt.Fprintln(os.Stderr, "打开文件失败:", err)
          os.Exit(1)
      }
      if len(targets) == 0 {
          fmt.Fprintln(os.Stderr, "无有效目标")
          os.Exit(0)
      }

      fmt.Printf("目标数: %d  每目标 ping次数: %d  超时: %ds\n", len(targets), count, timeout)
      fmt.Println("探测中...")

      summaries := make([]summary, 0, len(targets))
      for _, t := range targets {
          s := summarize(probe(t, count, timeout))
          summaries = append(summaries, s)
      }

      fmt.Printf("\n=== Ping 报告 === %s  每个目标 ping %d 次\n\n",
          time.Now().Format("2006-01-02 15:04:05"), count)
      fmt.Printf("%-24s %4s %4s %7s %9s %9s %9s %9s\n",
          "目标", "发送", "收到", "丢包%", "最小ms", "平均ms", "最大ms", "mdevms")
      fmt.Println(strings.Repeat("-", 80))

      byLat := append([]summary{}, summaries...)
      sort.SliceStable(byLat, func(i, j int) bool {
          if byLat[i].Loss != byLat[j].Loss {
              return byLat[i].Loss < byLat[j].Loss
          }
          return byLat[i].Avg < byLat[j].Avg
      })
      for _, s := range byLat {
          fmt.Printf("%-24s %4d %4d %6.1f%% %9s %9s %9s %9s\n",
              s.Target, s.Sent, s.Recv, s.Loss, fmtOrDash(s.Min), fmtOrDash(s.Avg),
              fmtOrDash(s.Max), fmtOrDash(s.Mdev))
      }

      fmt.Println("\n--- 延迟排行 (按平均延迟升序, 排除全丢包) ---")
      alive := append([]summary{}, summaries...)
      sort.SliceStable(alive, func(i, j int) bool { return alive[i].Avg < alive[j].Avg })
      i := 0
      for _, s := range alive {
          if s.Recv == 0 {
              continue
          }
          i++
          fmt.Printf("  %2d. %-24s 平均 %.2f ms  丢包 %.0f%%\n", i, s.Target, s.Avg, s.Loss)
      }

      fmt.Println("\n--- 丢包排行 (按丢包率降序) ---")
      byLoss := append([]summary{}, summaries...)
      sort.SliceStable(byLoss, func(i, j int) bool { return byLoss[i].Loss > byLoss[j].Loss })
      for k, s := range byLoss {
          flag := ""
          if s.Recv == 0 {
              flag = " ★全丢"
          }
          fmt.Printf("  %2d. %-24s 丢包 %.0f%% (%d/%d)%s\n",
              k+1, s.Target, s.Loss, s.Sent-s.Recv, s.Sent, flag)
      }
  }
