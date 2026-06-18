#!/usr/bin/env node
'use strict';
/**
 * CCM Permission Dialog Hook (PreToolUse)
 * =========================================
 * 在 Claude Code 执行写操作前弹出 macOS 原生确认框，替代终端里的 y/n 提示。
 * 对 Read / LS 等只读工具不弹窗（直接放行）。
 *
 * 超时 60 秒无响应自动拒绝。
 */

const { execFileSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const HOME = process.env.HOME || process.env.USERPROFILE;
const CCM_DIR = path.join(HOME, '.ccm');
const LOG_FILE = path.join(CCM_DIR, 'permission-dialog.log');

// 需要弹窗确认的工具（写操作）
const DIALOG_TOOLS = new Set(['Bash', 'Write', 'Edit', 'MultiEdit', 'NotebookEdit']);

function log(msg) {
  try {
    fs.appendFileSync(LOG_FILE, `[${new Date().toISOString().slice(11, 19)}] ${msg}\n`);
  } catch {}
}

function formatInfo(toolName, toolInput) {
  switch (toolName) {
    case 'Bash': {
      const cmd = (toolInput.command || '').trim();
      return cmd.length > 400 ? cmd.slice(0, 400) + '\n…(已截断)' : cmd;
    }
    case 'Write':
      return `写入文件:\n${toolInput.file_path || '(unknown)'}`;
    case 'Edit':
    case 'MultiEdit':
      return `编辑文件:\n${toolInput.file_path || '(unknown)'}`;
    default:
      return JSON.stringify(toolInput).slice(0, 300);
  }
}

function showDialog(toolName, info) {
  // AppleScript 内双引号和反斜杠需要转义
  const safe = info
    .replace(/\\/g, '\\\\')
    .replace(/"/g, '\\"')
    .replace(/\n/g, '\\n');

  const title = `Claude Code — ${toolName}`;
  const safeTitle = title.replace(/"/g, '\\"');

  const script = `
try
  set dlg to display dialog "${safe}" ¬
    with title "${safeTitle}" ¬
    buttons {"拒绝", "允许"} ¬
    default button "允许" ¬
    cancel button "拒绝" ¬
    with icon caution
  return button returned of dlg
on error errMsg
  return "error:" & errMsg
end try
`.trim();

  try {
    const result = execFileSync('osascript', ['-e', script], {
      timeout: 0,
      encoding: 'utf-8',
    }).trim();
    return result === '允许';
  } catch {
    return false;
  }
}

function respond(permissionDecision, reason) {
  process.stdout.write(JSON.stringify({
    hookSpecificOutput: {
      hookEventName: 'PreToolUse',
      permissionDecision,
      permissionDecisionReason: reason,
    },
  }));
}

async function main() {
  let inputData = '';
  process.stdin.setEncoding('utf-8');
  for await (const chunk of process.stdin) {
    inputData += chunk;
  }

  let input = {};
  try {
    input = JSON.parse(inputData);
  } catch {
    process.exit(0);
    return;
  }

  const toolName = input.tool_name || '';
  const toolInput = input.tool_input || {};

  if (!DIALOG_TOOLS.has(toolName)) {
    // 只读工具，直接放行
    process.exit(0);
    return;
  }

  const info = formatInfo(toolName, toolInput);
  log(`dialog: tool=${toolName} info=${info.slice(0, 80).replace(/\n/g, ' ')}`);

  const allowed = showDialog(toolName, info);
  log(`decision: ${allowed ? 'allow' : 'deny'}`);

  if (allowed) {
    respond('allow', 'dialog: user approved');
  } else {
    respond('deny', 'dialog: user denied or timeout');
  }
  process.exit(0);
}

main().catch(e => {
  log(`fatal: ${e.message}`);
  process.exit(0);
});
