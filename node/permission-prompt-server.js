#!/usr/bin/env node
/**
 * Claude Code Permission Prompt MCP Server
 *
 * 当 Claude Code 需要执行敏感操作时，弹出系统对话框请求用户确认。
 *
 * 用法：
 *   claude --permission-prompt-tool mcp__permission__ask_permission
 *
 * 配置到 ~/.claude.json：
 *   {
 *     "mcpServers": {
 *       "permission": {
 *         "command": "node",
 *         "args": ["/path/to/permission-prompt-server.js"]
 *       }
 *     }
 *   }
 */

const readline = require("readline");

// ─── MCP stdio 传输层 ───────────────────────────────────────────────

const rl = readline.createInterface({ input: process.stdin });
let buffer = "";

rl.on("line", (line) => {
  buffer += line;
  try {
    const msg = JSON.parse(buffer);
    buffer = "";
    handleMessage(msg);
  } catch {
    // 继续累积（换行分隔的 JSON）
  }
});

function send(obj) {
  process.stdout.write(JSON.stringify(obj) + "\n");
}

// ─── MCP 协议处理 ───────────────────────────────────────────────────

function handleMessage(msg) {
  const { id, method, params } = msg;

  switch (method) {
    case "initialize":
      send({
        jsonrpc: "2.0",
        id,
        result: {
          protocolVersion: "2024-11-05",
          capabilities: { tools: {} },
          serverInfo: { name: "permission-prompt", version: "1.0.0" },
        },
      });
      break;

    case "notifications/initialized":
      // 无需响应
      break;

    case "tools/list":
      send({
        jsonrpc: "2.0",
        id,
        result: {
          tools: [
            {
              name: "ask_permission",
              description:
                "当 Claude Code 需要执行操作时，弹出系统对话框请求用户确认",
              inputSchema: {
                type: "object",
                properties: {
                  tool_name: {
                    type: "string",
                    description: "要执行的工具名称，如 Bash、Write、Edit",
                  },
                  tool_input: {
                    type: "object",
                    description: "工具的输入参数",
                  },
                },
                required: ["tool_name", "tool_input"],
              },
            },
          ],
        },
      });
      break;

    case "tools/call":
      if (params?.name === "ask_permission") {
        handlePermissionRequest(id, params.arguments);
      } else {
        send({
          jsonrpc: "2.0",
          id,
          error: { code: -32601, message: "Unknown tool" },
        });
      }
      break;

    default:
      if (id !== undefined) {
        send({
          jsonrpc: "2.0",
          id,
          error: { code: -32601, message: `Method not found: ${method}` },
        });
      }
  }
}

// ─── 权限请求处理 ────────────────────────────────────────────────────

async function handlePermissionRequest(id, args) {
  const toolName = args?.tool_name || "未知工具";
  const toolInput = args?.tool_input || {};

  // 格式化显示内容
  const inputSummary = formatInput(toolName, toolInput);
  const message = `Claude Code 请求执行操作\n\n工具：${toolName}\n${inputSummary}\n\n是否允许？`;

  try {
    const allowed = await showDialog(message, toolName);

    // 返回符合 --permission-prompt-tool 协议的响应
    send({
      jsonrpc: "2.0",
      id,
      result: {
        content: [
          {
            type: "text",
            text: JSON.stringify({
              behavior: allowed ? "allow" : "deny",
              // updatedInput 可选，允许修改工具入参
            }),
          },
        ],
      },
    });
  } catch (err) {
    // 对话框出错时默认拒绝，保守处理
    send({
      jsonrpc: "2.0",
      id,
      result: {
        content: [
          {
            type: "text",
            text: JSON.stringify({ behavior: "deny" }),
          },
        ],
      },
    });
  }
}

// ─── 格式化工具入参摘要 ──────────────────────────────────────────────

function formatInput(toolName, input) {
  switch (toolName) {
    case "Bash":
      return `命令：${input.command || ""}`;
    case "Write":
    case "Edit":
    case "MultiEdit":
      return `文件：${input.file_path || input.path || ""}`;
    case "Read":
      return `读取：${input.file_path || input.path || ""}`;
    case "WebFetch":
      return `URL：${input.url || ""}`;
    default: {
      // 通用：取前两个字段显示
      const keys = Object.keys(input).slice(0, 2);
      if (keys.length === 0) return "";
      return keys.map((k) => `${k}：${String(input[k]).slice(0, 80)}`).join("\n");
    }
  }
}

// ─── 系统弹窗（macOS / Linux）────────────────────────────────────────

function showDialog(message, toolName) {
  return new Promise((resolve, reject) => {
    const { execFile } = require("child_process");
    const platform = process.platform;

    if (platform === "darwin") {
      // macOS：使用 osascript
      const script = `
        set result to display dialog ${JSON.stringify(message)} ¬
          with title "Claude Code 权限请求" ¬
          buttons {"拒绝", "允许"} ¬
          default button "允许" ¬
          with icon caution
        return button returned of result
      `;
      execFile("osascript", ["-e", script], (err, stdout) => {
        if (err) {
          // 用户点了关闭按钮 → 拒绝
          resolve(false);
        } else {
          resolve(stdout.trim() === "允许");
        }
      });
    } else if (platform === "linux") {
      // Linux：优先用 zenity，其次 kdialog，最后终端交互
      execFile("which", ["zenity"], (err) => {
        if (!err) {
          execFile(
            "zenity",
            [
              "--question",
              "--title=Claude Code 权限请求",
              `--text=${message}`,
              "--ok-label=允许",
              "--cancel-label=拒绝",
              "--width=400",
            ],
            (err2) => resolve(!err2)
          );
        } else {
          execFile("which", ["kdialog"], (err3) => {
            if (!err3) {
              execFile(
                "kdialog",
                ["--yesno", message, "--title", "Claude Code 权限请求"],
                (err4) => resolve(!err4)
              );
            } else {
              // fallback：终端交互
              terminalPrompt(message).then(resolve).catch(reject);
            }
          });
        }
      });
    } else {
      // Windows / 其他：终端交互
      terminalPrompt(message).then(resolve).catch(reject);
    }
  });
}

// ─── 终端交互 fallback ───────────────────────────────────────────────

function terminalPrompt(message) {
  return new Promise((resolve) => {
    // 直接写到 /dev/tty，不影响 stdout（stdout 用于 MCP 协议）
    const fs = require("fs");
    try {
      const tty = fs.createWriteStream("/dev/tty");
      const ttyIn = fs.createReadStream("/dev/tty");
      const ttyRl = readline.createInterface({
        input: ttyIn,
        output: tty,
        terminal: true,
      });
      tty.write(`\n⚠️  ${message}\n允许？(y/n): `);
      ttyRl.once("line", (line) => {
        ttyRl.close();
        tty.end();
        resolve(line.trim().toLowerCase() === "y");
      });
    } catch {
      // 实在没有 tty，默认拒绝
      resolve(false);
    }
  });
}
