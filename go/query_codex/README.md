# query-codex

一个零依赖的 Go 命令行程序，用来快速查看当前 Codex 账号的短周期、长周期剩余额度和重置时间。

## 使用

```bash
go build -o query-codex .
./query-codex
```

查询前，程序会自动读取 `~/.codex/.env` 并将其中的变量设置到当前进程环境；文件不存在时会直接跳过。随后读取 `~/.codex/auth.json`，因此需要先通过 Codex CLI 登录：

```bash
codex login
```

其他选项：

```bash
./query-codex -json                 # 查看完整接口响应
./query-codex -timeout 3s           # 设置请求超时
./query-codex -auth-file /path/auth.json
./query-codex -endpoint https://example.test/usage
```

也可以用环境变量 `CODEX_USAGE_ENDPOINT` 覆盖额度接口，便于接口迁移或代理转发。

## 说明

- 只使用 Go 标准库，不依赖 Codex CLI 进程或其他项目目录。
- `.env` 支持空行、注释、`export KEY=VALUE`、单双引号；其中的值会覆盖当前进程已有的同名变量。
- 登录令牌只作为 HTTPS Authorization 请求头发送，不会打印到终端。
- 默认额度地址是 Codex 客户端当前使用的 ChatGPT 内部接口；它不是稳定的公开 API。如上游调整，可通过 `-endpoint` 或环境变量覆盖。
