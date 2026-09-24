#!/usr/bin/env python3
"""
读取 test_list，对每个 IP/域名执行 ping，统计延迟和丢包率，排行后输出报告。
"""

import subprocess
import re
import sys
import platform
from statistics import mean, stdev
from datetime import datetime

LIST_FILE = "test_list"
PING_COUNT = 4
TIMEOUT = 2  # 单次 ping 超时（秒）


def ping_once(target):
    """调用系统 ping 一次，返回 (success: bool, latency_ms: float|None)。"""
    system = platform.system().lower()
    if system == "windows":
        cmd = ["ping", "-n", "1", "-w", str(TIMEOUT * 1000), target]
    else:
        cmd = ["ping", "-c", "1", "-W", str(TIMEOUT), target]
    try:
        out = subprocess.run(cmd, capture_output=True, text=True, timeout=TIMEOUT + 2).stdout
    except (subprocess.TimeoutExpired, Exception):
        return False, None
    m = re.search(r"time[=<]\s*([\d.]+)\s*ms", out)
    if m:
        return True, float(m.group(1))
    return False, None


def probe(target, count):
    """对单个目标 ping count 次。返回 dict 含延迟统计与丢包率。"""
    samples = []
    success = 0
    for _ in range(count):
        ok, lat = ping_once(target)
        if ok:
            samples.append(lat)
            success += 1
    loss_pct = (count - success) / count * 100
    return {
        "target": target,
        "sent": count,
        "recv": success,
        "loss_pct": loss_pct,
        "min": min(samples) if samples else None,
        "avg": mean(samples) if samples else None,
        "max": max(samples) if samples else None,
        "mdev": stdev(samples) if len(samples) > 1 else 0.0,
 "samples": samples,
    }


def fmt(v, suf=""):
    return f"{v:.2f}{suf}" if v is not None else "—"


def load_targets(path):
    targets = []
    with open(path, "r", encoding="utf-8") as f:
        for line in f:
            t = line.strip()
            if not t or t.startswith("#"):
                continue
            targets.append(t)
    return targets


def print_report(results, count):
    print(f"\n=== Ping报告 === 时间: {datetime.now():%Y-%m-%d %H:%M:%S}  每个目标 ping {count} 次\n")
    header = f"{'目标':<24}{'发送':>4}{'收到':>4}{'丢包%':>8}{'最小ms':>9}{'平均ms':>9}{'最大ms':>9}{'mdevms':>9}"
    print(header)
    print("-" * len(header.encode("gbk", "ignore") or header.encode()))
    print("-" * 80)
    for r in sorted(results, key=lambda x: (x["loss_pct"], x["avg"] or 1e9)):
        print(
            f"{r['target']:<24}{r['sent']:>4}{r['recv']:>4}{r['loss_pct']:>8.1f}"
            f"{fmt(r['min']):>9}{fmt(r['avg']):>9}{fmt(r['max']):>9}{fmt(r['mdev']):>9}"
        )

    print("\n---延迟排行 (按平均延迟升序，排除全丢包) ---")
    alive = [r for r in results if r["recv"] > 0]
    for i, r in enumerate(sorted(alive, key=lambda x: x["avg"]), 1):
        print(f"  {i:>2}. {r['target']:<24} 平均 {r['avg']:.2f} ms  丢包 {r['loss_pct']:.0f}%")

    print("\n--- 丢包排行 (按丢包率降序) ---")
    for i, r in enumerate(sorted(results, key=lambda x: -x["loss_pct"]), 1):
        flag = " ★全丢" if r["recv"] == 0 else ""
        print(f"  {i:>2}. {r['target']:<24} 丢包 {r['loss_pct']:.0f}% ({r['sent']-r['recv']}/{r['sent']}){flag}")


def main():
    count = PING_COUNT
    list_file = LIST_FILE
    if len(sys.argv) >= 2:
        list_file = sys.argv[1]
    if len(sys.argv) >= 3:
        count = int(sys.argv[2])

    targets = load_targets(list_file)
    if not targets:
        print(f"{list_file} 中没有可用条目。")
        return    print(f"开始探测 {len(targets)} 个目标，每个 {count} 次 ...")
    results = [probe(t, count) for t in targets]
    print_report(results, count)


if __name__ == "__main__":
    main()
