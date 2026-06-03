#!/usr/bin/env bash
# test_163_cn2.sh - 判断到目标 IP 的链路走中国电信 163 骨干网 (AS4134)
# 还是 CN2 (AS4809)。
#
# 用法:
#   ./test_163_cn2.sh <目标IP> [目标IP ...]
#
# 依赖: mtr, curl
# API : ip-api.com (免费, 限速 45 req/min; 本脚本使用批量接口, 一次最多 100 个)

set -u

API="http://ip-api.com/batch"
FIELDS="query,as,org,isp,country,status,message"

# ---------- 已知中国电信骨干段（兜底用，避免对每个跳都打 API） ----------
# 163 骨干网 AS4134 关键段
AS4134_CIDRS=(
    "202.97.0.0/19"   "202.97.32.0/19"  "202.97.64.0/19"  "202.97.96.0/19"
    "202.97.128.0/17" "202.97.33.0/24"  "202.97.34.0/24"  "202.97.35.0/24"
    "202.97.36.0/24"  "202.97.37.0/24"  "202.97.38.0/24"  "202.97.39.0/24"
    "202.97.40.0/24"  "202.97.41.0/24"  "202.97.42.0/24"  "202.97.43.0/24"
    "202.97.44.0/24"  "202.97.45.0/24"  "202.97.46.0/24"  "202.97.47.0/24"
    "202.97.48.0/24"  "202.97.49.0/24"  "202.97.50.0/24"  "202.97.51.0/24"
    "202.97.52.0/24"  "202.97.53.0/24"  "202.97.54.0/24"  "202.97.55.0/24"
    "202.97.56.0/24"  "202.97.57.0/24"  "202.97.58.0/24"  "202.97.59.0/24"
    "202.97.60.0/24"  "202.97.61.0/24"  "202.97.62.0/24"  "202.97.63.0/24"
    "202.97.64.0/18"  "202.97.128.0/18"
    "61.135.0.0/16"   "61.136.0.0/16"   "61.137.0.0/16"   "61.138.0.0/16"
    "61.139.0.0/16"   "61.140.0.0/16"   "61.141.0.0/16"   "61.142.0.0/16"
    "61.143.0.0/16"   "61.144.0.0/16"   "61.145.0.0/16"   "61.146.0.0/16"
    "61.147.0.0/16"   "61.148.0.0/16"   "61.149.0.0/16"   "61.150.0.0/16"
    "61.151.0.0/16"   "61.152.0.0/16"   "61.153.0.0/16"   "61.154.0.0/16"
    "61.155.0.0/16"   "61.156.0.0/16"   "61.157.0.0/16"   "61.158.0.0/16"
    "61.159.0.0/16"   "116.6.0.0/16"    "116.7.0.0/16"    "116.8.0.0/16"
    "116.9.0.0/16"    "116.10.0.0/16"   "116.11.0.0/16"   "116.12.0.0/16"
    "116.13.0.0/16"   "116.14.0.0/16"   "116.15.0.0/16"   "116.16.0.0/16"
    "222.42.0.0/16"   "222.43.0.0/16"   "222.44.0.0/16"   "222.45.0.0/16"
    "222.46.0.0/16"   "222.47.0.0/16"   "222.48.0.0/16"   "222.49.0.0/16"
)

# CN2 AS4809 关键段
AS4809_CIDRS=(
    "59.43.0.0/16"    "59.42.0.0/16"    "59.40.0.0/16"    "59.39.0.0/16"
    "59.38.0.0/16"    "59.37.0.0/16"    "59.36.0.0/16"    "59.35.0.0/16"
    "59.34.0.0/16"    "59.33.0.0/16"    "59.32.0.0/16"
    "202.97.0.0/19"   # 部分 CN2 节点
    "218.30.0.0/16"   "218.31.0.0/16"
    "115.168.0.0/14"
    "110.96.0.0/16"
    "110.97.0.0/16"
    "124.72.0.0/16"
    "124.73.0.0/16"
    "124.74.0.0/16"
    "124.75.0.0/16"
    "124.76.0.0/16"
    "180.97.0.0/16"
)

# ---------- IP 工具函数（无 python 依赖，用 bash + awk 实现 CIDR 判断） ----------

# ip2int <ip> -> 32位整数
ip2int() {
    local IFS=.
    read -r o1 o2 o3 o4 <<<"$1"
    echo $(( (o1 << 24) + (o2 << 16) + (o3 << 8) + o4 ))
}

# 检查 ip 是否在 cidr 内。cider_check <ip> <cidr>
cidr_check() {
    local ip=$1 cidr=$2
    local net=${cidr%/*} mask=${cidr#*/}
    local ip_n net_n
    ip_n=$(ip2int "$ip")
    net_n=$(ip2int "$net")
    if [[ $mask -eq 0 ]]; then
        return 0
    fi
    local full=$((32 - mask))
    local m
    if (( full >= 32 )); then m=0
    else m=$(( 0xFFFFFFFF << full )); m=$(( m & 0xFFFFFFFF ))
    fi
    if (( (ip_n & m) == (net_n & m) )); then
        return 0
    fi
    return 1
}

# 从字符串中提取合法 IPv4
extract_ipv4() {
    grep -oE '\b([0-9]{1,3}\.){3}[0-9]{1,3}\b' <<<"$1" | head -1
}

# 判断 ip 是否在 AS4134 段内
is_as4134() {
    local ip=$1 cidr
    for cidr in "${AS4134_CIDRS[@]}"; do
        cidr_check "$ip" "$cidr" && return 0
    done
    return 1
}

# 判断 ip 是否在 AS4809 段内
is_as4809() {
    local ip=$1 cidr
    for cidr in "${AS4809_CIDRS[@]}"; do
        cidr_check "$ip" "$cidr" && return 0
    done
    return 1
}

# ---------- 批量查询 IP 归属 ----------
# query_ips "ip1\nip2\n..." -> 输出三列: ip<TAB>as<TAB>org
query_ips() {
    local body=$1
    # 构造 batch 请求体: [{"query":"1.2.3.4","fields":"..."}, ...]
    local req
    req=$(awk -v fields="$FIELDS" '{
        gsub(/"/, "\\\"", $1)
        printf "{\"query\":\"%s\",\"fields\":\"%s\"}", $1, fields
        if (NR > 1) printf ","
    } END { print "" }' <<<"$body")

    # 一次最多 100 个, 超过分批
    local out
    out=$(curl -s --max-time 10 -H "Content-Type: application/json" \
              -X POST -d "[$req]" "$API" 2>/dev/null)

    # 用 awk 解析
    echo "$out" | awk -v ORS='\n' '
        function get_str(s, key,    n, m, r) {
            # 提取 "key":"..." 的值
            r = ""
            m = "\""
            n = index(s, "\"" key "\":\"")
            if (n == 0) return ""
            n += length(key) + 4
            while (n <= length(s) && substr(s, n, 1) != "\"") {
                c = substr(s, n, 1)
                if (c == "\\") { n++; c = substr(s, n, 1) }
                r = r c
                n++
            }
            return r
        }
        {
            for (i = 1; i <= length($0); i++) {
                c = substr($0, i, 1)
                if (depth == 0 && c == "{") depth = 1
                else if (depth == 1) {
                    if (c == "}") { depth = 0; print ""; break }
                    cur = cur c
                    if (c == "}") { depth = 0; print ""; cur = ""; break }
                } else if (c == "}") { depth = 0; print ""; cur = "" }
                if (c == "{" && depth == 0) depth = 1
            }
        }'
}

# 解析 mtr 输出
# 输入: mtr -r -c 1 -n 输出
# 输出: 每行一个 hop IP
parse_mtr() {
    # mtr -r 报告的格式: 第 1 行是表头 (HOST:..., Loss%..., Snt..., ...)
    # 后面每行: " 1.|-- 1.2.3.4  0.0%  1  0.5"
    awk '
        NR == 1 { next }
        /^[[:space:]]*$/ { next }
        {
            # 提取 IP
            for (i = 1; i <= NF; i++) {
                if ($i ~ /^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$/) {
                    print $i
                    break
                }
            }
        }
    '
}

# ---------- 主流程 ----------
if [[ $# -lt 1 ]]; then
    echo "用法: $0 <目标IP> [目标IP ...]" >&2
    echo "示例: $0 8.8.8.8" >&2
    exit 1
fi

for target in "$@"; do
    echo "=========================================================="
    echo "目标: $target"
    echo "=========================================================="

    # 1. mtr 跟踪
    echo "[*] 正在使用 mtr 收集路由 ..."
    mtr_out=$(mtr -r -c 1 -n -w "$target" 2>/dev/null)
    if [[ -z "$mtr_out" ]]; then
        echo "[!] mtr 未获取到结果" >&2
        continue
    fi

    # 2. 解析每一跳 IP（去重保留首次出现）
    hops=$(echo "$mtr_out" | parse_mtr | awk '!seen[$0]++')
    if [[ -z "$hops" ]]; then
        echo "[!] 未解析到任何跳" >&2
        continue
    fi

    echo "[*] 探测到的路径节点 (按跳数):"
    echo "$hops" | nl

    # 3. 先用本地 CIDR 列表初筛
    local_4134=()
    local_4809=()
    unknown=()
    declare -A hop_as_map

    while IFS= read -r ip; do
        [[ -z "$ip" ]] && continue
        if is_as4134 "$ip"; then
            local_4134+=("$ip")
            hop_as_map[$ip]="AS4134 (163, 命中本地段)"
        elif is_as4809 "$ip"; then
            local_4809+=("$ip")
            hop_as_map[$ip]="AS4809 (CN2, 命中本地段)"
        else
            unknown+=("$ip")
        fi
    done <<<"$hops"

    # 4. 对未识别 IP 批量查询
    api_4134=()
    api_4809=()
    api_other=()
    if (( ${#unknown[@]} > 0 )); then
        echo "[*] ${#unknown[@]} 个节点未命中本地段, 调用 ip-api.com 批量查询 ..."
        # 分批 100 个
        tmpdir=$(mktemp -d)
        printf "%s\n" "${unknown[@]}" > "$tmpdir/list.txt"
        split -l 100 "$tmpdir/list.txt" "$tmpdir/batch_"

        for bf in "$tmpdir"/batch_*; do
            [[ -f "$bf" ]] || continue
            # 构造 JSON 数组
            req_body=$(awk -v fields="$FIELDS" '{
                gsub(/"/, "\\\"", $1)
                if (NR > 1) printf ","
                printf "{\"query\":\"%s\",\"fields\":\"%s\"}", $1, fields
            } END { print "" }' "$bf")

            resp=$(curl -s --max-time 15 -H "Content-Type: application/json" \
                        -X POST -d "[$req_body]" "$API" 2>/dev/null)
            [[ -z "$resp" ]] && continue

            # 简单逐对象解析: 用 grep -oE 抓每个 { ... } 块
            # 因为 ip-api.com 的 batch 返回扁平对象数组, 每个对象无嵌套, 这种方式够用
            echo "$resp" | tr '}' '\n' | while IFS= read -r obj; do
                obj="{$obj}"
                ip=$(  echo "$obj" | grep -oE '"query"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed -E 's/.*"query"[[:space:]]*:[[:space:]]*"([^"]*)".*/\1/')
                asn=$( echo "$obj" | grep -oE '"as"[[:space:]]*:[[:space:]]*"[^"]*"'   | head -1 | sed -E 's/.*"as"[[:space:]]*:[[:space:]]*"([^"]*)".*/\1/')
                [[ -z "$ip" ]] && continue
                echo "$ip|$asn" >> "$tmpdir/result.txt"
            done
        done

        if [[ -f "$tmpdir/result.txt" ]]; then
            while IFS='|' read -r ip asn; do
                [[ -z "$ip" ]] && continue
                if   [[ "$asn" == AS4134* ]]; then
                    api_4134+=("$ip")
                    hop_as_map[$ip]="AS4134 (163, $asn)"
                elif [[ "$asn" == AS4809* ]]; then
                    api_4809+=("$ip")
                    hop_as_map[$ip]="AS4809 (CN2, $asn)"
                else
                    api_other+=("$ip")
                    [[ -n "$asn" ]] && hop_as_map[$ip]="$asn" || hop_as_map[$ip]="未知"
                fi
            done < "$tmpdir/result.txt"
        fi
        rm -rf "$tmpdir"
    fi

    # 5. 汇总输出
    echo
    echo "---------- 路径分析 ----------"
    printf "%-4s %-18s %s\n" "跳" "IP" "归属"
    i=0
    while IFS= read -r ip; do
        ((i++))
        printf "%-4d %-18s %s\n" "$i" "$ip" "${hop_as_map[$ip]:-未查询}"
    done <<<"$hops"

    echo
    echo "---------- 统计 ----------"
    echo "命中 AS4134 (163 骨干网): $(( ${#local_4134[@]} + ${#api_4134[@]} )) 个"
    (( ${#local_4134[@]} )) && echo "  本地段: ${local_4134[*]}"
    (( ${#api_4134[@]}   )) && echo "  API 段: ${api_4134[*]}"
    echo "命中 AS4809 (CN2):       $(( ${#local_4809[@]} + ${#api_4809[@]} )) 个"
    (( ${#local_4809[@]} )) && echo "  本地段: ${local_4809[*]}"
    (( ${#api_4809[@]}   )) && echo "  API 段: ${api_4809[*]}"
    echo "其他 ISP / 未知:         ${#api_other[@]} 个"
    (( ${#api_other[@]} )) && echo "  ${api_other[*]}"

    # 6. 结论
    echo
    echo "---------- 结论 ----------"
    has_163=0
    has_cn2=0
    (( ${#local_4134[@]} + ${#api_4134[@]} )) && has_163=1
    (( ${#local_4809[@]} + ${#api_4809[@]} )) && has_cn2=1

    if (( has_163 && has_cn2 )); then
        echo "→ 混合: 路径上同时出现 163 与 CN2 节点（可能在边界互联处）。"
    elif (( has_163 )); then
        echo "→ 走 163 骨干网 (AS4134)"
    elif (( has_cn2 )); then
        echo "→ 走 CN2 (AS4809)"
    else
        echo "→ 路径上未发现中国电信 163 / CN2 节点（可能为其他运营商或本地局域网）。"
    fi
    echo
done
