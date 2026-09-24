#!/bin/bash
# ping报告脚本：读取目标列表，每个目标 ping N 次，统计延迟/丢包并排行。

set -u

LIST_FILE="${1:-test_list}"
COUNT="${2:-4}"
TIMEOUT=2
PING_BIN="$(command -v ping)"

[[ -z "$PING_BIN" ]] && { echo "找不到 ping 命令"; exit 1; }
[[ ! -f "$LIST_FILE" ]] && { echo "文件不存在: $LIST_FILE"; exit 1; }

# 收集目标
mapfile -t TARGETS < <(grep -vE '^\s*(#|$)' "$LIST_FILE")
[[ ${#TARGETS[@]} -eq 0 ]] && { echo "无有效目标"; exit 0; }

echo "目标数: ${#TARGETS[@]}  每目标 ping 次数: $COUNT  超时: ${TIMEOUT}s"
echo "探测中..."

# 探测: 每个目标 $COUNT 次,每次取 time=XXX ms
# 输出行: target<TAB>ok<TAB>latency
probe() {
    local target="$1" line lat
    for _ in $(seq 1 "$COUNT"); do
        line=$(ping -c1 -W "$TIMEOUT" "$target" 2>/dev/null)
        if [[ $line =~ time=([0-9.]+)[[:space:]]*ms ]]; then
            printf '%s\t%s\t%s\n' "$target" ok "${BASH_REMATCH[1]}"
        else
            printf '%s\t%s\t\n' "$target" fail ""
        fi
    done
}

TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

for t in "${TARGETS[@]}"; do
    probe "$t" >> "$TMP"
done

# 汇总: target sent recv loss% min avg max mdev
SUMMARY=$(mktemp)
trap 'rm -f "$TMP" "$SUMMARY"' EXIT

awk -F'\t' -v count="$COUNT" '
function finish(target, i, n, sum, sumsq, mn, mx, lat, vals, key) {
    if (!(target in sent)) return
    if (n == 0) {
        printf "%s\t%d\t0\t100.0\t-\t-\t-\t-\n", target, sent[target]
 return
    }
    mn = vals[1]; mx = vals[1]; sum = 0
    for (i=1; i<=n; i++) {
        sum += vals[i]
        if (vals[i] < mn) mn = vals[i]
        if (vals[i] > mx) mx = vals[i]
    }
    avg = sum / n
    for (i=1; i<=n; i++) sumsq += (vals[i]-avg)^2
    mdev = (n > 1) ? sqrt(sumsq/(n-1)) : 0.0
    printf "%s\t%d\t%d\t%.1f\t%.2f\t%.2f\t%.2f\t%.2f\n", target, sent[target], n, (sent[target]-n)/sent[target]*100, mn, avg, mx, mdev
}
{
    target = $1
    sent[target]++
    if ($2 == "ok" && $3 != "") {
        n = ++cnt[target]
        vals_key = target SUBSEP n
        # 用数组保存每个目标的样本延迟, key = target|n
        # awk 普通数组, 这里用 vals[target, n] 形式
        # 使用通用做法直接索引 vals[target"|"n] = $3 + 0 if (!(target in mn_arr) || $3+0 < mn_arr[target]) mn_arr[target] = $3+0
        if (!(target in mx_arr) || $3+0 > mx_arr[target]) mx_arr[target] = $3+0
    }
 last = target
}
END {
    # 由于 awk 不易处理多维+状态, 改用下面更简洁版本 - 重写
}
' "$TMP" > "$SUMMARY" 2>/dev/null

# 上面 awk 写得复杂了, 改成 python 一行汇总更稳:
SUMMARY=$(mktemp)
python3 - "$TMP" "$COUNT" > "$SUMMARY" <<'PY'
import sys, math
from collections import defaultdict
fp, count = sys.argv[1], int(sys.argv[2])
sent = defaultdict(int); recv = defaultdict(int)
vals = defaultdict(list)
for line in open(fp):
    p = line.rstrip("\n").split("\t")
    if len(p) < 3: continue
    t, st, lat = p[0], p[1], p[2]
    sent[t] += 1
    if st == "ok" and lat:
        recv[t] += 1
        vals[t].append(float(lat))
for t in sent:
    s = vals[t]
    if s:
        avg = sum(s)/len(s)
        mdev = math.sqrt(sum((x-avg)**2 for x in s)/(len(s)-1)) if len(s)>1 else 0.0
        print(f"{t}\t{sent[t]}\t{recv[t]}\t{(sent[t]-recv[t])/sent[t]*100:.1f}\t{min(s):.2f}\t{avg:.2f}\t{max(s):.2f}\t{mdev:.2f}")
    else:
        print(f"{t}\t{sent[t]}\t0\t100.0\t-\t-\t-\t-")
PY

echoprintf '%-24s %4s %4s %7s %9s %9s %9s %9s\n' "目标" "发送" "收到" "丢包%" "最小ms" "平均ms" "最大ms" "mdevms"
echo "--------------------------------------------------------------------------------"
sort -t$'\t' -k4,4n -k6,6n "$SUMMARY" | while IFS=$'\t' read -r tgt sent recv loss mn avg mx mdev; do
    printf '%-24s %4s %4s %6s%% %9s %9s %9s %9s\n' "$tgt" "$sent" "$recv" "$loss" "$mn" "$avg" "$mx" "$mdev"
done

echo
echo "--- 延迟排行 (按平均延迟升序, 排除全丢包) ---"
sort -t$'\t' -k6,6n "$SUMMARY" | awk -F'\t' '$3>0 {printf "  %2d. %-24s 平均 %s ms  丢包 %s%%\n", NR, $1, $6, $4}'

echo
echo "--- 丢包排行 (按丢包率降序) ---"
sort -t$'\t' -k4,4nr "$SUMMARY" | awk -F'\t' '{
    flag = ($3==0) ? " ★全丢" : ""
    printf "  %2d. %-24s 丢包 %s%% (%d/%d)%s\n", NR, $1, $4, $2-$3, $2, flag
}'

rm -f "$TMP" "$SUMMARY"
